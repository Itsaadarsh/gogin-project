package main

import (
	"gogin/config"
	"gogin/controller"
	"gogin/repository"
	"gogin/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDBConnection()
	userRep        repository.UserRepository = repository.NewUserRepository(db)
	authSevice     service.AuthService       = service.NewAuthService(userRep)
	jwtService     service.JWTService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRep)
	authController controller.AuthController = controller.NewAuthController(authSevice, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
)

func main() {
	defer config.CloseDBConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoute := r.Group("api/user")
	{
		userRoute.POST("/update", userController.Update)
		userRoute.GET("/profile", userController.Profile)
	}

	r.Run()
}
