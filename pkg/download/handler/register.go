package handler

import (
	sess "app/internal/session"
	"app/pkg/auth"
	"app/pkg/config"
	"app/pkg/download"
	"app/pkg/logger"
	"app/pkg/user"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, cfg *config.Config, uc download.DownloadUseCase, luc logger.LogUseCase, auc auth.AuthUseCase, suc sess.SessUseCase, uuc user.UserUseCase) {
	h := NewDownloadHandler(cfg, uc, luc, auc, suc, uuc)

	downloadRoute := router.Group("/download")
	{
		downloadRoute.GET("/:file_dir/:file_name", h.DownloadFile)
	}
}
