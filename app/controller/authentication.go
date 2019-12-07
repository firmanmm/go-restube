package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	authService *service.AuthenticationService
}

func (a *AuthenticationController) HandleRegister(ctx *gin.Context) {
	username := ctx.PostForm("username")
	if len(username) < 5 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Valid \"username\" post form is required",
		})
		return
	}

	password := ctx.PostForm("password")
	if len(password) < 5 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Valid \"password\" post form is required",
		})
		return
	}
	if err := a.authService.NewAuthentication(username, password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account Created",
	})
}

func (a *AuthenticationController) HandleAuthentication(ctx *gin.Context) {
	username := ctx.PostForm("username")
	if len(username) < 5 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Valid \"username\" post form is required",
		})
		return
	}

	password := ctx.PostForm("password")
	if len(password) < 5 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Valid \"password\" post form is required",
		})
		return
	}
	sessionID, err := a.authService.Authenticate(username, password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Usage": fmt.Sprintf("%d MB", size/1024/1024),
	})
}

func (a *AuthenticationController) HandleListUsage(ctx *gin.Context) {
	sessionID := ctx.GetHeader("Authorization")
	auth, err := a.authService.FindBySessionID(sessionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if auth.Username != "myadmin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You are not allowed to access this resource",
		})
		return
	}
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 0)
	offset, _ := strconv.ParseInt(ctx.Query("offset"), 10, 0)
	usages, err := a.authService.ListAllUsage(uint(limit), uint(offset))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, usages)
}

func NewAuthenticationController(authService *service.AuthenticationService) *AuthenticationController {
	instance := new(AuthenticationController)
	instance.authService = authService
	return instance
}
