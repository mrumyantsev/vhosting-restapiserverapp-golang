package handler

import (
	sess "app/internal/session"
	"app/pkg/auth"
	"app/pkg/config"
	"app/pkg/logger"
	"app/pkg/user"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, cfg *config.Config, uc auth.AuthUseCase, uuc user.UserUseCase,
	suc sess.SessUseCase, luc logger.LogUseCase) {
	h := NewAuthHandler(cfg, uc, uuc, suc, luc)

	authRoute := router.Group("/auth")
	{
		authRoute.POST("/signin", h.SignIn)
		authRoute.POST("/change_password", h.ChangePassword)
		authRoute.GET("/signout", h.SignOut)
	}
}
