package main

import (
	"gogin/config"
	"gogin/controller"
	"gogin/middleware"
	"gogin/repository"
	"gogin/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDBConnection()

	userRep    repository.UserRepository = repository.NewUserRepository(db)
	bookRep    repository.BookRepository = repository.NewBookRepository(db)
	authSevice service.AuthService       = service.NewAuthService(userRep)

	userService service.UserService = service.NewUserService(userRep)
	bookSer     service.BookService = service.NewBookService(bookRep)
	jwtService  service.JWTService  = service.NewJWTService()

	authController controller.AuthController = controller.NewAuthController(authSevice, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookController controller.BookController = controller.NewBookController(bookSer, jwtService)
)

func main() {
	defer config.CloseDBConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoute := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoute.POST("/update", userController.Update)
		userRoute.GET("/profile", userController.Profile)
	}

	bookRoute := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoute.GET("/", bookController.All)
		bookRoute.POST("/", bookController.Insert)
		bookRoute.PUT("/:id", bookController.Update)
		bookRoute.DELETE("/:id", bookController.Delete)
		bookRoute.GET("/:id", bookController.FindByID)
	}

	r.Run()
}
