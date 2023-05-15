package handler

import (
	sess "app/internal/session"
	"app/internal/video"
	"app/pkg/auth"
	"app/pkg/config"
	"app/pkg/logger"
	"app/pkg/user"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, cfg *config.Config, uc video.VideoUseCase, luc logger.LogUseCase,
	auc auth.AuthUseCase, suc sess.SessUseCase, uuc user.UserUseCase) {
	h := NewVideoHandler(cfg, uc, luc, auc, suc, uuc)

	videoRoute := router.Group("/video")
	{
		videoRoute.POST("", h.CreateVideo)
		videoRoute.GET(":id", h.GetVideo)
		videoRoute.GET("all", h.GetAllVideos)
		videoRoute.PATCH(":id", h.PartiallyUpdateVideo)
		videoRoute.DELETE(":id", h.DeleteVideo)
	}
}
