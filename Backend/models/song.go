// Backend/models/song.go
/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 21/11/2024
Lenguaje: Golang
Descipcion: Asignacion de la clase song, con sus respectivas
funciones para el manejo de datos
(para la estructura de datos)*/
package models

import (
	"fmt"
	"strings"
	"time"
)

type Song struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	Genre      string    `json:"genre"`
	FileSize   int       `json:"file_size"`
	AddedAt    time.Time `json:"added_at"`
	PlayCount  int       `json:"play_count"`
	LastPlayed time.Time `json:"last_played"`
}

// Constructor para Song
func NewSong(id int, title, artist, genre string, fileSize int) (*Song, error) {
	if title == "" || artist == "" || genre == "" {
		return nil, fmt.Errorf("todos los campos son requeridos")
	}

	song := &Song{
		ID:       id,
		Title:    title,
		Artist:   artist,
		Genre:    genre,
		FileSize: fileSize,
		AddedAt:  time.Now(),
	}

	if err := song.ValidateSize(); err != nil {
		return nil, err
	}

	return song, nil
}

// ValidateSize  Se asegura que la canción no sea más grande de 10MB antes de agregarla
func (s *Song) ValidateSize() error {
	const MaxSize = 10 * 1024 * 1024 // 10MB en bytes
	if s.FileSize > MaxSize {
		return fmt.Errorf("el tamaño del archivo excede el límite de 10MB")
	}
	return nil
}

/*
	Getters entre ellos el GetTitle nos devuelve  el nombre de la canción

GetArtist nos devuelve el nombre de la canción, Nos dice quién es el artista de la canción
Nos indica qué tipo de música es (rock, pop, etc.)
*/
func (s *Song) GetTitle() string {
	return s.Title
}

func (s *Song) GetArtist() string {
	return s.Artist
}

func (s *Song) GetGenre() string {
	return s.Genre
}

// Metodos adicionales para el manejo de la data

// GetFormattedFileSize Muestra el tamaño del archivo
func (s *Song) GetFormattedFileSize() string {
	const unit = 1024
	if s.FileSize < unit {
		return fmt.Sprintf("%d B", s.FileSize)
	}
	div, exp := int64(unit), 0
	for n := s.FileSize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(s.FileSize)/float64(div), "KMGTPE"[exp])
}

// GetInfo Junta toda la información importante de la canción en un solo texto
func (s *Song) GetInfo() string {
	return fmt.Sprintf("ID: %d\nTítulo: %s\nArtista: %s\nGénero: %s\nTamaño: %s\nReproducciones: %d",
		s.ID, s.Title, s.Artist, s.Genre, s.GetFormattedFileSize(), s.PlayCount)
}

// IncrementPlayCount  Suma uno al contador cada vez que se reproduce la canción
func (s *Song) IncrementPlayCount() {
	s.PlayCount++
	s.LastPlayed = time.Now()
}

// UpdateMetadata  Permite cambiar la información básica de la canción
func (s *Song) UpdateMetadata(title, artist, genre string) error {
	if title != "" {
		s.Title = title
	}
	if artist != "" {
		s.Artist = artist
	}
	if genre != "" {
		s.Genre = genre
	}
	return nil
}

// MatchesSearch Revisa si la canción coincide con lo que está buscando el usuario
func (s *Song) MatchesSearch(query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(s.Title), query) ||
		strings.Contains(strings.ToLower(s.Artist), query) ||
		strings.Contains(strings.ToLower(s.Genre), query)
}
