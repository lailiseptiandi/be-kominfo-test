package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lailiseptiandi/be-kominfo-test/internal/handlers"
	"github.com/lailiseptiandi/be-kominfo-test/internal/middlewares"
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
	roomRepo := repository.NewRoomRepository(dbs)
	bookRepo := repository.NewBookRepository(dbs)

	// service
	authService := services.NewUserService(userRepo)
	roomService := services.NewRoomService(roomRepo)
	bookService := services.NewBookingService(bookRepo)

	// handler
	authHandler := handlers.NewAuthHandler(authService)
	roomHandler := handlers.NewRoomHandler(roomService)
	bookHandler := handlers.NewBookingHandler(bookService)

	// auth
	routeGroup := r.Group("/api/v1")
	routeGroup.POST("/login", authHandler.Login)

	routeGroup.POST("/rooms", middlewares.AuthRequired(), roomHandler.Create)
	routeGroup.GET("/rooms/available", middlewares.AuthRequired(), roomHandler.GetAvailable)
	routeGroup.POST("/bookings", middlewares.AuthRequired(), bookHandler.Create)
	routeGroup.PATCH("/bookings/:id/cancel", middlewares.AuthRequired(), bookHandler.Cancel)
}
