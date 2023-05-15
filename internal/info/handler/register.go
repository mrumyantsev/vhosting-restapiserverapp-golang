package handler

import (
	"app/internal/info"
	sess "app/internal/session"
	"app/pkg/auth"
	"app/pkg/config"
	sconfig "app/pkg/config_stream"
	"app/pkg/logger"
	"app/pkg/user"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, cfg *config.Config, scfg *sconfig.Config, uc info.InfoUseCase, luc logger.LogUseCase,
	auc auth.AuthUseCase, suc sess.SessUseCase, uuc user.UserUseCase) {
	h := NewInfoHandler(cfg, scfg, uc, luc, auc, suc, uuc)

	infoRoute := router.Group("/info")
	{
		infoRoute.POST("", h.CreateInfo)
		infoRoute.GET(":id", h.GetInfo)
		infoRoute.GET("all", h.GetAllInfos)
		infoRoute.PATCH(":id", h.PartiallyUpdateInfo)
		infoRoute.DELETE(":id", h.DeleteInfo)
	}
}
