// Backend/models/database/mysql.go
/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 05/12/2024
Lenguaje: Golang
Descripción: Asignación de la clase MySQL, con sus respectivas
funciones para la conexión con la base de datos y recomendaciones.
*/

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

// GetRecommendationsByGenres devuelve una lista de canciones recomendadas para múltiples géneros.
func GetRecommendationsByGenres(genres []string, limit int) (map[string][]string, error) {
	recommendations := make(map[string][]string)

	for _, genre := range genres {
		query := `
			SELECT title 
			FROM songs 
			WHERE genre = ? 
			ORDER BY created_at DESC 
			LIMIT ?`
		rows, err := db.Query(query, genre, limit)
		if err != nil {
			log.Printf("Error al obtener recomendaciones para el género %s: %v", genre, err)
			return nil, err
		}
		defer rows.Close()

		var songs []string
		for rows.Next() {
			var title string
			if err := rows.Scan(&title); err != nil {
				log.Printf("Error al leer la fila de canciones: %v", err)
				return nil, err
			}
			songs = append(songs, title)
		}

		recommendations[genre] = songs
	}

	return recommendations, nil
}

// AddFavorite agrega una canción favorita para un usuario
func AddFavorite(userID, songID string) error {
	query := "INSERT INTO user_favorites (user_id, song_id) VALUES (?, ?)"
	_, err := db.Exec(query, userID, songID)
	return err
}

// GetFavorites obtiene las canciones favoritas de un usuario
func GetFavorites(userID string) ([]map[string]string, error) {
	query := `
		SELECT songs.id, songs.title, songs.artist 
		FROM songs
		JOIN user_favorites ON songs.id = user_favorites.song_id
		WHERE user_favorites.user_id = ?`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []map[string]string
	for rows.Next() {
		var id, title, artist string
		if err := rows.Scan(&id, &title, &artist); err != nil {
			return nil, err
		}
		favorites = append(favorites, map[string]string{
			"id":     id,
			"title":  title,
			"artist": artist,
		})
	}
	return favorites, nil
}
