package middleware

import (
	"net/http"

	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *service.AuthenticationService
}

func (a *AuthMiddleware) Handle(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")
	if len(auth) < 10 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Header \"Authorization\" is empty")
		return
	}

	if err := a.authService.CheckAuthentication(auth); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Next()
}

func NewAuthMiddleware(authService *service.AuthenticationService) *AuthMiddleware {
	instance := new(AuthMiddleware)
	instance.authService = authService
	return instance
}
