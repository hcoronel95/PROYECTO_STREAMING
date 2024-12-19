// Backend/models/handlers/user.go
/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 05/12/2024
Lenguaje: Golang
Descipcion: Asignacion de la clase users, con sus respectivas
funciones para el manejo de rutas
*/

// handlers/user.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	db *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error al leer el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Todos los campos son requeridos", http.StatusBadRequest)
		return
	}

	// Verificar si el email ya existe
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		http.Error(w, "Error al verificar el email", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "El email ya está registrado", http.StatusConflict)
		return
	}

	// Insertar nuevo usuario
	result, err := h.db.Exec(
		"INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, user.Password, "user",
	)
	if err != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)
	user.Password = "" // No devolver la contraseña

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	//  ID del usuario del token
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	var user User
	err := h.db.QueryRow(
		"SELECT id, name, email, role FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err == sql.ErrNoRows {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error al obtener el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Nuevo endpoint: Obtener recomendaciones personalizadas para el usuario
func (h *UserHandler) GetUserRecommendations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el ID del usuario de los parámetros de consulta
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	// Consulta para obtener las recomendaciones basadas en las preferencias del usuario
	rows, err := h.db.Query(`
		SELECT s.id, s.title, s.artist, s.album
		FROM user_preferences up
		JOIN songs s ON up.song_id = s.id
		WHERE up.user_id = ?`, userID)
	if err != nil {
		http.Error(w, "Error al obtener las recomendaciones", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Construir las recomendaciones
	var recommendations []struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Album  string `json:"album"`
	}

	for rows.Next() {
		var song struct {
			ID     int    `json:"id"`
			Title  string `json:"title"`
			Artist string `json:"artist"`
			Album  string `json:"album"`
		}
		if err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.Album); err != nil {
			http.Error(w, "Error al procesar las recomendaciones", http.StatusInternalServerError)
			return
		}
		recommendations = append(recommendations, song)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}
