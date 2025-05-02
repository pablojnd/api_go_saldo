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

	// Agrupar rutas API
	api := router.Group("/api")
	{
		// Rutas para saldos
		api.GET("/saldos", controllers.GetSaldos)
		api.GET("/saldos/:codArt", controllers.GetSaldoById)
	}
}
