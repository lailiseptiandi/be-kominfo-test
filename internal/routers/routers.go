package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lailiseptiandi/be-kominfo-test/internal/handlers"
	"github.com/lailiseptiandi/be-kominfo-test/internal/repository"
	"github.com/lailiseptiandi/be-kominfo-test/internal/services"
	"gorm.io/gorm"
)

func InitRoutes(dbs *gorm.DB, r *gin.Engine) {
	r.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Success check health")
	})

	// repo
	userRepo := repository.NewUserRepository(dbs)

	// service
	authService := services.NewUserService(userRepo)

	// handler
	authHandler := handlers.NewAuthHandler(authService)

	// auth
	routeGroup := r.Group("/api/v1")
	routeGroup.POST("/login", authHandler.Login)
}
