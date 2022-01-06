package controller

import (
	"fmt"
	"gogin/dto"
	"gogin/entity"
	"gogin/helper"
	"gogin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type BookController interface {
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	All(ctx *gin.Context)
	FindByID(ctx *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JWTService
}

func NewBookController(bookSer service.BookService, jwtSer service.JWTService) BookController {
	return &bookController{
		bookService: bookSer,
		jwtService:  jwtSer,
	}
}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, errToken := c.jwtService.ValidateToken(token)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return fmt.Sprint("%v", claims["user_id"])
}

func (c *bookController) Insert(ctx *gin.Context) {
	var bookInsertDTO dto.BookCreatedDTO
	errDTO := ctx.ShouldBind(&bookInsertDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	userID := c.getUserIDByToken(authHeader)
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err == nil {
		bookInsertDTO.UserID = userIDUint
	}

	bookInserted := c.bookService.Insert(bookInsertDTO)
	response := helper.BuildResponse(true, "OK", bookInserted)
	ctx.JSON(http.StatusOK, response)

}

func (c *bookController) Update(ctx *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := ctx.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userId, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userId, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserID = id
		}
		bookUpdated := c.bookService.Update(bookUpdateDTO)
		response := helper.BuildResponse(true, "OK", bookUpdated)
		ctx.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You don't have permission", "Incorrect user ID", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
	}

}
func (c *bookController) Delete(ctx *gin.Context) {
	var book entity.Book

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("No param ID was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	book.ID = id
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userId, book.ID) {
		c.bookService.Delete(book)
		response := helper.BuildResponse(true, "OK", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You don't have permission", "Incorrect user ID", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
	}
}

func (c *bookController) All(ctx *gin.Context) {
	var books []entity.Book = c.bookService.All()
	response := helper.BuildResponse(true, "OK", books)
	ctx.JSON(http.StatusOK, response)
}

func (c *bookController) FindByID(ctx *gin.Context) {

	bookUint, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("No param ID was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	bookRes := c.bookService.FindByID(bookUint)
	if (bookRes == entity.Book{}) {
		response := helper.BuildErrorResponse("Data Not Found", "No data found with the given ID", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "OK", bookRes)
		ctx.JSON(http.StatusOK, response)
	}
}
