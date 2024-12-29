package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// AddFavoriteHandler maneja la solicitud para agregar canciones favoritas
func AddFavoriteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			SongID string `json:"songId"`
			UserID string `json:"userId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO user_favorites (user_id, song_id) VALUES (?, ?)"
		_, err := db.Exec(query, req.UserID, req.SongID)
		if err != nil {
			http.Error(w, "Error adding favorite", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// GetFavoritesHandler maneja la solicitud para obtener canciones favoritas
func GetFavoritesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		if userID == "" {
			http.Error(w, "Missing userId", http.StatusBadRequest)
			return
		}

		query := `
			SELECT songs.id, songs.title, songs.artist 
			FROM songs
			JOIN user_favorites ON songs.id = user_favorites.song_id
			WHERE user_favorites.user_id = ?`

		rows, err := db.Query(query, userID)
		if err != nil {
			http.Error(w, "Error fetching favorites", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var favorites []map[string]string
		for rows.Next() {
			var id, title, artist string
			if err := rows.Scan(&id, &title, &artist); err != nil {
				http.Error(w, "Error scanning favorites", http.StatusInternalServerError)
				return
			}
			favorites = append(favorites, map[string]string{
				"id":     id,
				"title":  title,
				"artist": artist,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(favorites)
	}
}
