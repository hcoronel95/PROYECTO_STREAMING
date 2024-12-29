/*
Autores: Henry Aliaga / Ismael Espinoza
Fecha: 21/11/2024
Lenguaje: Golang
Descripción: Asignación de la clase Library con sus respectivas
funciones para el manejo de datos, incluyendo favoritos y búsqueda.
*/

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
	SongMap     map[string][]Song `json:"song_map"`
	Favorites   map[int][]int     `json:"favorites"` // Mapa de favoritos por usuario
	CreatedAt   time.Time         `json:"created_at"`
	LastUpdated time.Time         `json:"last_updated"`
	TotalSize   int64             `json:"total_size"` // Tamaño total en bytes
}

const (
	MaxSongs    = 60               // Límite máximo de canciones
	MaxSongSize = 10 * 1024 * 1024 // 10 MB en bytes
)

// NewLibrary crea una nueva instancia de Library
func NewLibrary(userID int) *Library {
	return &Library{
		UserID:      userID,
		Songs:       make([]Song, 0),
		SongMap:     make(map[string][]Song),
		Favorites:   make(map[int][]int),
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

// AddFavorite agrega una canción a los favoritos de un usuario
func (l *Library) AddFavorite(userID, songID int) error {
	// Verificar si la canción está en la biblioteca
	_, err := l.GetSongByID(songID)
	if err != nil {
		return fmt.Errorf("la canción no se encuentra en la biblioteca")
	}

	// Agregar la canción a la lista de favoritos del usuario
	if _, exists := l.Favorites[userID]; !exists {
		l.Favorites[userID] = []int{}
	}
	l.Favorites[userID] = append(l.Favorites[userID], songID)
	l.LastUpdated = time.Now()
	return nil
}

// RemoveFavorite elimina una canción de los favoritos de un usuario
func (l *Library) RemoveFavorite(userID, songID int) error {
	// Verificar si la canción está en la biblioteca
	_, err := l.GetSongByID(songID)
	if err != nil {
		return fmt.Errorf("la canción no se encuentra en la biblioteca")
	}

	// Eliminar la canción de la lista de favoritos
	if _, exists := l.Favorites[userID]; exists {
		for i, id := range l.Favorites[userID] {
			if id == songID {
				l.Favorites[userID] = append(l.Favorites[userID][:i], l.Favorites[userID][i+1:]...)
				break
			}
		}
	}

	l.LastUpdated = time.Now()
	return nil
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

// FormatLastUpdated retorna la fecha de última actualización en formato legible
func (l *Library) FormatLastUpdated() string {
	return l.LastUpdated.Format("02-01-2006 15:04:05")
}
