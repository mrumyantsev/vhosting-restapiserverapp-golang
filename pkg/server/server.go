package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/group"
	grouphandler "app/internal/group/handler"
	grouprepo "app/internal/group/repository"
	groupusecase "app/internal/group/usecase"
	"app/internal/info"
	infohandler "app/internal/info/handler"
	inforepo "app/internal/info/repository"
	infousecase "app/internal/info/usecase"
	msg "app/internal/messages"
	perm "app/internal/permission"
	permhandler "app/internal/permission/handler"
	permrepo "app/internal/permission/repository"
	permusecase "app/internal/permission/usecase"
	sess "app/internal/session"
	sessrepo "app/internal/session/repository"
	sessusecase "app/internal/session/usecase"
	"app/internal/video"
	videohandler "app/internal/video/handler"
	videorepo "app/internal/video/repository"
	videousecase "app/internal/video/usecase"
	"app/pkg/auth"
	authhandler "app/pkg/auth/handler"
	authrepo "app/pkg/auth/repository"
	authusecase "app/pkg/auth/usecase"
	"app/pkg/config"
	"app/pkg/config_stream"
	sconfig "app/pkg/config_stream"
	"app/pkg/download"
	downloadhandler "app/pkg/download/handler"
	downloadusecase "app/pkg/download/usecase"
	"app/pkg/logger"
	logrepo "app/pkg/logger/repository"
	logusecase "app/pkg/logger/usecase"
	"app/pkg/stream"
	streamhandler "app/pkg/stream/handler"
	streamrepo "app/pkg/stream/repository"
	streamusecase "app/pkg/stream/usecase"
	"app/pkg/user"
	userhandler "app/pkg/user/handler"
	userrepo "app/pkg/user/repository"
	userusecase "app/pkg/user/usecase"

	"github.com/gin-gonic/gin"
)

type App struct {
	httpServer      *http.Server
	cfg             *config.Config
	scfg            *sconfig.Config
	userUseCase     user.UserUseCase
	authUseCase     auth.AuthUseCase
	sessUseCase     sess.SessUseCase
	logUseCase      logger.LogUseCase
	groupUseCase    group.GroupUseCase
	permUseCase     perm.PermUseCase
	infoUseCase     info.InfoUseCase
	videoUseCase    video.VideoUseCase
	StreamUC        stream.StreamUseCase
	downloadUseCase download.DownloadUseCase
}

func NewApp(cfg *config.Config) *App {
	userRepo := userrepo.NewUserRepository(cfg)
	authRepo := authrepo.NewAuthRepository(cfg)
	sessRepo := sessrepo.NewSessRepository(cfg)
	logRepo := logrepo.NewLogRepository(cfg)
	groupRepo := grouprepo.NewGroupRepository(cfg)
	permRepo := permrepo.NewPermRepository(cfg)
	infoRepo := inforepo.NewInfoRepository(cfg)
	videoRepo := videorepo.NewVideoRepository(cfg)
	streamRepo := streamrepo.NewStreamRepository(cfg)

	scfg := &config_stream.Config{}

	return &App{
		cfg:             cfg,
		scfg:            scfg,
		userUseCase:     userusecase.NewUserUseCase(cfg, userRepo),
		authUseCase:     authusecase.NewAuthUseCase(cfg, authRepo),
		sessUseCase:     sessusecase.NewSessUseCase(sessRepo, authRepo),
		logUseCase:      logusecase.NewLogUseCase(logRepo),
		groupUseCase:    groupusecase.NewGroupUseCase(groupRepo),
		permUseCase:     permusecase.NewPermUseCase(permRepo),
		infoUseCase:     infousecase.NewInfoUseCase(infoRepo),
		videoUseCase:    videousecase.NewVideoUseCase(videoRepo),
		StreamUC:        streamusecase.NewStreamUseCase(cfg, scfg, streamRepo),
		downloadUseCase: downloadusecase.NewDownloadUseCase(cfg),
	}
}

func (a *App) Run() error {
	// Debug mode.
	if a.cfg.ServerDebugEnable {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init engine.
	router := gin.New()

	// Init middleware.
	router.Use(CORSMiddleware())

	// Check for web directory exists and register routes.
	if _, err := os.Stat("./web"); !os.IsNotExist(err) {
		router.LoadHTMLGlob("./web/templates/*")
		streamhandler.RegisterTemplateHTTPEndpoints(router, a.cfg, a.scfg, a.StreamUC,
			a.userUseCase, a.logUseCase, a.authUseCase, a.sessUseCase)
	}

	router.StaticFS("/static", http.Dir("./web/static"))

	// Register routes.
	authhandler.RegisterHTTPEndpoints(router, a.cfg, a.authUseCase, a.userUseCase,
		a.sessUseCase, a.logUseCase)
	userhandler.RegisterHTTPEndpoints(router, a.cfg, a.userUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase)
	grouphandler.RegisterHTTPEndpoints(router, a.cfg, a.groupUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase)
	permhandler.RegisterHTTPEndpoints(router, a.cfg, a.permUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase, a.groupUseCase)
	infohandler.RegisterHTTPEndpoints(router, a.cfg, a.scfg, a.infoUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase)
	videohandler.RegisterHTTPEndpoints(router, a.cfg, a.videoUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase)
	streamhandler.RegisterStreamingHTTPEndpoints(router, a.cfg, a.scfg, a.StreamUC,
		a.userUseCase, a.logUseCase, a.authUseCase, a.sessUseCase)
	downloadhandler.RegisterHTTPEndpoints(router, a.cfg, a.downloadUseCase, a.logUseCase,
		a.authUseCase, a.sessUseCase, a.userUseCase)

	// Set HTTP server params.
	a.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", a.cfg.ServerHost, a.cfg.ServerPort),
		Handler:        router,
		ReadTimeout:    time.Duration(a.cfg.ServerReadTimeoutSeconds) * time.Second,
		WriteTimeout:   time.Duration(a.cfg.ServerWriteTimeoutSeconds) * time.Second,
		MaxHeaderBytes: a.cfg.ServerMaxHeaderBytes,
	}

	// Start HTTP server.
	var err error

	go func() {
		err = a.httpServer.ListenAndServe()
	}()

	time.Sleep(50 * time.Millisecond)

	if err != nil {
		return errors.New(fmt.Sprintf("Cannot start server. Error: %s.", err.Error()))
	}

	a.cfg.ServerIP = getOutboundIP()
	logger.Print(msg.InfoServerStartedSuccessfullyAtLocalAddress(a.cfg.ServerIP, a.cfg.ServerPort))

	// Start videostreams worker.
	go a.StreamUC.ServeStreams()

	// Listen for interrupt signal from keyboard.
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logger.Print(msg.InfoRecivedSignal(sig))
		done <- true
	}()

	<-done

	// Shut down HTTP server.
	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return errors.New(fmt.Sprintf("Cannot shut down the server correctly. Error: %s.", err.Error()))
	}

	logger.Print(msg.InfoServerShutedDownCorrectly())

	return nil
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		logger.Print(msg.WarningCannotGetLocalIP(err))
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
