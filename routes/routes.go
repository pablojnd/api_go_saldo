package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pablojnd/api_go_saldo/controllers"
)

// Setup configura todas las rutas de la API
func Setup(router *gin.Engine) {
	// Ruta de estado
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// Única ruta de la API
	api := router.Group("/api")
	{
		api.GET("/allsaldos", controllers.GetAllSaldos)
	}
}
