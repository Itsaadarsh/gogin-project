package controller

import (
	"fmt"
	"gogin/dto"
	"gogin/helper"
	"gogin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserController interface {
	Update(ctx *gin.Context)
	Profile(ctx *gin.Context)
}

type userController struct {
	userSer service.UserService
	jwtSer  service.JWTService
}

func NewUserController(userSer service.UserService, jwtSer service.JWTService) UserController {
	return &userController{
		userSer: userSer,
		jwtSer:  jwtSer,
	}
}

func (c *userController) Update(ctx *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errDTO := ctx.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtSer.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)

	if err != nil {
		panic(err.Error())
	}

	userUpdateDTO.ID = id
	updatedUser := c.userSer.Update(userUpdateDTO)
	response := helper.BuildResponse(true, "OK", updatedUser)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) Profile(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtSer.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	user := c.userSer.Profile(fmt.Sprintf("%v", claims["user_id"]))
	response := helper.BuildResponse(true, "OK", user)
	ctx.JSON(http.StatusOK, response)
}
