package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pablojnd/api_go_saldo/config"
	"github.com/pablojnd/api_go_saldo/database"
	"github.com/pablojnd/api_go_saldo/routes"
)

func main() {
	// Cargar variables de entorno
	config.LoadEnv()

	// Configurar modo de Gin según entorno
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode // Usar ReleaseMode por defecto en producción
	}
	gin.SetMode(ginMode)

	// Inicializar conexión a la base de datos
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Crear router Gin
	r := gin.Default()

	// Configurar confianza en proxies - solo confiar en localhost
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Configurar rutas
	routes.Setup(r)

	// Iniciar servidor usando la variable de entorno correcta
	port := config.GetEnv("SERVER_PORT", "8080")
	log.Printf("Servidor iniciado en el puerto %s en modo %s", port, ginMode)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
