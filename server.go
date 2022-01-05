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
	authController controller.AuthController = controller.NewAuthController(authSevice, jwtService)
)

func main() {
	defer config.CloseDBConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	r.Run()
}
