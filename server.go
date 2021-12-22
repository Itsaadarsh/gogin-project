package main

import (
	"gogin/config"
	"gogin/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDBConnection()
	authController controller.AuthController = controller.NewAuthController()
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
