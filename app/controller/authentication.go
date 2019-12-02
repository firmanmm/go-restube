package controller

import (
	"fmt"
	"net/http"

	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	authService *service.AuthenticationService
}

func (a *AuthenticationController) HandleRegister(ctx *gin.Context) {
	username := ctx.PostForm("username")
	if len(username) < 5 {
		ctx.String(http.StatusBadRequest, "Valid \"username\" post form is required")
		return
	}

	password := ctx.PostForm("password")
	if len(password) < 5 {
		ctx.String(http.StatusBadRequest, "Valid \"password\" post form is required")
		return
	}
	if err := a.authService.NewAuthentication(username, password); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.String(http.StatusOK, "Account Created")
}

func (a *AuthenticationController) HandleAuthentication(ctx *gin.Context) {
	username := ctx.PostForm("username")
	if len(username) < 5 {
		ctx.String(http.StatusBadRequest, "Valid \"username\" post form is required")
		return
	}

	password := ctx.PostForm("password")
	if len(password) < 5 {
		ctx.String(http.StatusBadRequest, "Valid \"password\" post form is required")
		return
	}
	sessionID, err := a.authService.Authenticate(username, password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Message":   "Use this session ID in the Authorization Header",
		"SessionID": sessionID,
	})
}

func (a *AuthenticationController) HandleGetByteDownloaded(ctx *gin.Context) {
	sessionID := ctx.GetHeader("Authorization")
	size, err := a.authService.GetByteDownloaded(sessionID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Usage": fmt.Sprintf("%d MB", size/1024/1024),
	})
}

func NewAuthenticationController(authService *service.AuthenticationService) *AuthenticationController {
	instance := new(AuthenticationController)
	instance.authService = authService
	return instance
}
