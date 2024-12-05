package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// GetDefaultConfig retorna la configuración por defecto
func GetDefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "RLPZFVKdjq5BhlGe09ol",
		DBName:   "streaming_music",
	}
}

// InitDB inicializa la conexión a la base de datos
func InitDB(config Config) error {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
		)

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Error conectando a MySQL: %v", err)
			return
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Error haciendo ping a MySQL: %v", err)
			return
		}

		// Configurar el pool de conexiones
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
	})

	return err
}

// GetDB retorna la instancia de la base de datos
func GetDB() *sql.DB {
	return db
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
