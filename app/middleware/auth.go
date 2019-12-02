package middleware

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
}

func (a *AuthMiddleware) Handle(ctx *gin.Context) {
	//TODO : Do
}

func NewAuthMiddleware() *AuthMiddleware {
	instance := new(AuthMiddleware)
	return instance
}
