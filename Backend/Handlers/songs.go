//(para el manejo de rutas/endpoints)

// Backend/models/handlers/songs.go
/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 05/12/2024
Lenguaje: Golang
Descipcion: Asignacion de la clase songs, con sus respectivas
funciones para el manejo de rutas
*/
// handlers/song.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type SongHandler struct {
	db *sql.DB
}

type Song struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Genre    string `json:"genre"`
	FileSize int    `json:"file_size"`
	FilePath string `json:"file_path"`
}

func NewSongHandler(db *sql.DB) *SongHandler {
	return &SongHandler{db: db}
}

func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.db.Query("SELECT id, title, artist, genre, file_size, file_path FROM songs")
	if err != nil {
		http.Error(w, "Error al obtener canciones", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.Genre, &song.FileSize, &song.FilePath)
		if err != nil {
			http.Error(w, "Error al leer canción", http.StatusInternalServerError)
			return
		}
		songs = append(songs, song)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var song Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Error al leer el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	// Validar tamaño máximo (10MB)
	if song.FileSize > 10*1024*1024 {
		http.Error(w, "El tamaño del archivo excede el límite permitido", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec(
		"INSERT INTO songs (title, artist, genre, file_size, file_path) VALUES (?, ?, ?, ?, ?)",
		song.Title, song.Artist, song.Genre, song.FileSize, song.FilePath,
	)
	if err != nil {
		http.Error(w, "Error al guardar la canción", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	song.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

func (h *SongHandler) GetSongByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	songID := r.URL.Query().Get("id")
	if songID == "" {
		http.Error(w, "ID de canción no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(songID)
	if err != nil {
		http.Error(w, "ID de canción inválido", http.StatusBadRequest)
		return
	}

	var song Song
	err = h.db.QueryRow(
		"SELECT id, title, artist, genre, file_size, file_path FROM songs WHERE id = ?",
		id,
	).Scan(&song.ID, &song.Title, &song.Artist, &song.Genre, &song.FileSize, &song.FilePath)

	if err == sql.ErrNoRows {
		http.Error(w, "Canción no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error al obtener la canción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}
