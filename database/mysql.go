package database

import (
	"fmt"

	"github.com/pablojnd/api_go_saldo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL inicializa la conexión a MySQL
func InitMySQL() error {
	// Usar los nombres de variables correctos que coinciden con el archivo .env
	host := config.GetEnv("MYSQL_HOST", "localhost")
	port := config.GetEnv("MYSQL_PORT", "3306")
	user := config.GetEnv("MYSQL_USER", "root")
	password := config.GetEnv("MYSQL_PASSWORD", "")
	dbname := config.GetEnv("MYSQL_DATABASE", "")

	// Cambiando utf8mb4 a utf8 para compatibilidad con versiones antiguas de MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10s",
		user, password, host, port, dbname)

	fmt.Printf("Conectando a la base de datos en %s:%s\n", host, port)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	fmt.Println("Conexión a la base de datos establecida correctamente")
	return nil
}
