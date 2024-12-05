// Backend/models/library.go
/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 21/11/2024
Lenguaje: Golang
Descipcion: Asignacion de la clase library , con sus respectivas
funciones para el manejo de datos
(para la estructura de datos)*/
package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Library representa la gestión de canciones de un usuario
type Library struct {
	ID          int               `json:"id"`
	UserID      int               `json:"user_id"`
	Songs       []Song            `json:"songs"`
	SongMap     map[string][]Song // Organización por género
	CreatedAt   time.Time         `json:"created_at"`
	LastUpdated time.Time         `json:"last_updated"`
	TotalSize   int64             `json:"total_size"` // Tamaño total en bytes
}

const (
	MaxSongs    = 60               // Límite máximo de canciones
	MaxSongSize = 10 * 1024 * 1024 // 10MB en bytes
)

// NewLibrary crea una nueva instancia de Library
func NewLibrary(userID int) *Library {
	return &Library{
		UserID:      userID,
		Songs:       make([]Song, 0),
		SongMap:     make(map[string][]Song),
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
		TotalSize:   0,
	}
}

// AddSong añade una nueva canción a la biblioteca
func (l *Library) AddSong(song Song) error {
	// Verificar límite de canciones
	if len(l.Songs) >= MaxSongs {
		return fmt.Errorf("límite de canciones alcanzado (máximo %d)", MaxSongs)
	}

	// Verificar si la canción ya existe
	for _, s := range l.Songs {
		if s.ID == song.ID {
			return errors.New("la canción ya existe en la biblioteca")
		}
	}

	// Verificar si el nuevo tamaño total excedería el límite
	newTotalSize := l.TotalSize + int64(song.FileSize)
	if newTotalSize > int64(MaxSongs*MaxSongSize) {
		return errors.New("el tamaño total de la biblioteca excedería el límite")
	}

	// Añadir la canción
	l.Songs = append(l.Songs, song)
	l.SongMap[song.Genre] = append(l.SongMap[song.Genre], song)
	l.TotalSize += int64(song.FileSize)
	l.LastUpdated = time.Now()
	return nil
}

// GetSongsByGenre retorna todas las canciones de un género específico
func (l *Library) GetSongsByGenre(genre string) []Song {
	return l.SongMap[genre]
}

// RemoveSong elimina una canción de la biblioteca
func (l *Library) RemoveSong(songID int) error {
	for i, song := range l.Songs {
		if song.ID == songID {
			// Eliminar de la lista principal
			l.Songs = append(l.Songs[:i], l.Songs[i+1:]...)

			// Eliminar del mapa de géneros
			genre := song.Genre
			songs := l.SongMap[genre]
			for j, s := range songs {
				if s.ID == songID {
					l.SongMap[genre] = append(songs[:j], songs[j+1:]...)
					break
				}
			}

			// Actualizar tamaño total
			l.TotalSize -= int64(song.FileSize)
			l.LastUpdated = time.Now()
			return nil
		}
	}
	return errors.New("canción no encontrada")
}

// SearchSongs busca canciones por título, artista o género
func (l *Library) SearchSongs(query string) []Song {
	var results []Song
	query = strings.ToLower(query)

	for _, song := range l.Songs {
		if strings.Contains(strings.ToLower(song.Title), query) ||
			strings.Contains(strings.ToLower(song.Artist), query) ||
			strings.Contains(strings.ToLower(song.Genre), query) {
			results = append(results, song)
		}
	}
	return results
}

// GetLibraryStats retorna estadísticas de la biblioteca
func (l *Library) GetLibraryStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// Estadísticas básicas
	stats["total_songs"] = len(l.Songs)
	stats["total_size_mb"] = float64(l.TotalSize) / (1024 * 1024)
	stats["available_slots"] = MaxSongs - len(l.Songs)

	// Conteo por género
	genreCounts := make(map[string]int)
	for genre := range l.SongMap {
		genreCounts[genre] = len(l.SongMap[genre])
	}
	stats["songs_per_genre"] = genreCounts

	return stats
}

// GetAvailableSpace retorna el espacio disponible en MB
func (l *Library) GetAvailableSpace() float64 {
	totalAllowed := int64(MaxSongs * MaxSongSize)
	available := totalAllowed - l.TotalSize
	return float64(available) / (1024 * 1024)
}

// GetGenres retorna la lista de géneros únicos en la biblioteca
func (l *Library) GetGenres() []string {
	genres := make([]string, 0, len(l.SongMap))
	for genre := range l.SongMap {
		genres = append(genres, genre)
	}
	return genres
}

// GetSongByID busca una canción por su ID
func (l *Library) GetSongByID(id int) (*Song, error) {
	for _, song := range l.Songs {
		if song.ID == id {
			return &song, nil
		}
	}
	return nil, errors.New("canción no encontrada")
}

// ValidateLibraryLimits verifica si la biblioteca está dentro de los límites permitidos
func (l *Library) ValidateLibraryLimits() error {
	if len(l.Songs) > MaxSongs {
		return fmt.Errorf("la biblioteca excede el límite de %d canciones", MaxSongs)
	}

	maxTotalSize := int64(MaxSongs * MaxSongSize)
	if l.TotalSize > maxTotalSize {
		return fmt.Errorf("la biblioteca excede el límite de tamaño de %d MB", maxTotalSize/(1024*1024))
	}

	return nil
}

// FormatLastUpdated retorna la fecha de última actualización en formato legible
func (l *Library) FormatLastUpdated() string {
	return l.LastUpdated.Format("02-01-2006 15:04:05")
}
