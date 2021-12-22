package controller

import (
	"gogin/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
}

func NewAuthController() AuthController {
	return &authController{}
}

func (c *authController) Login(ctx *gin.Context) {
	res := helper.BuildResponse(true, "Login Route", "success")
	ctx.JSON(http.StatusOK, gin.H{
		"res": res,
	})
}
func (c *authController) Register(ctx *gin.Context) {
	res := helper.BuildResponse(true, "Register Route", "success")
	ctx.JSON(http.StatusOK, gin.H{
		"res": res,
	})
}
