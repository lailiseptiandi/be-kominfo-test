package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(dbs *gorm.DB, r *gin.Engine) {
	r.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Success check health")
	})
}
