// Backend/models/playback.go
/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 21/11/2024
Lenguaje: Golang
Descipcion: Asignacion de la clase playback, con sus respectivas
funciones para el manejo de datos
(para la estructura de datos)*/
package models

import (
	"errors"
	"fmt"
	"time"
)

type PlaybackStatus string

const (
	StatusPlaying   PlaybackStatus = "playing"
	StatusPaused    PlaybackStatus = "paused"
	StatusCompleted PlaybackStatus = "completed"
)

type Playback struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	SongID    int            `json:"song_id"`
	PlayedAt  time.Time      `json:"played_at"`
	Duration  int            `json:"duration"`  // Duración en segundos
	Completed bool           `json:"completed"` // Si se completó la reproducción
	Status    PlaybackStatus `json:"status"`
	PausedAt  time.Time      `json:"paused_at"`
}

// NewPlayback crea una nueva instancia de reproducción
func NewPlayback(userID, songID int) *Playback {
	return &Playback{
		UserID:    userID,
		SongID:    songID,
		PlayedAt:  time.Now(),
		Completed: false,
		Status:    StatusPlaying,
	}
}

// Start inicia o reinicia la reproducción
func (p *Playback) Start() error {
	if p.Completed {
		return errors.New("la reproducción ya fue completada")
	}
	if p.Status == StatusPaused {
		p.Duration += int(time.Since(p.PausedAt).Seconds())
	}
	p.Status = StatusPlaying
	return nil
}

// Pause pausa la reproducción
func (p *Playback) Pause() error {
	if p.Completed {
		return errors.New("no se puede pausar una reproducción completada")
	}
	if p.Status != StatusPlaying {
		return errors.New("la reproducción no está en curso")
	}
	p.Status = StatusPaused
	p.PausedAt = time.Now()
	return nil
}

// CompletePlayback marca la reproducción como completada
func (p *Playback) CompletePlayback() {
	p.Completed = true
	p.Status = StatusCompleted
	p.Duration += int(time.Since(p.PlayedAt).Seconds())
}

// GetPlaybackDuration retorna la duración de la reproducción
func (p *Playback) GetPlaybackDuration() time.Duration {
	if !p.Completed {
		if p.Status == StatusPaused {
			return time.Duration(p.Duration) * time.Second
		}
		return time.Since(p.PlayedAt)
	}
	return time.Duration(p.Duration) * time.Second
}

// FormatPlayedAt retorna la fecha actual de reproduccion
func (p *Playback) FormatPlayedAt() string {
	return p.PlayedAt.Format("02-01-2006 15:04:05")
}

// IsRecentPlayback verifica si la reproducción es reciente (últimas 24 horas)
func (p *Playback) IsRecentPlayback() bool {
	return time.Since(p.PlayedAt) < 24*time.Hour
}

// GetStatus Informa si la canción está sonando, pausada o terminó
func (p *Playback) GetStatus() string {
	switch p.Status {
	case StatusPlaying:
		return "Reproduciendo"
	case StatusPaused:
		return "Pausado"
	case StatusCompleted:
		return "Completado"
	default:
		return "Desconocido"
	}
}

// GetFormattedDuration Muestra el tiempo en minutos y segundos
func (p *Playback) GetFormattedDuration() string {
	duration := p.GetPlaybackDuration()
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
