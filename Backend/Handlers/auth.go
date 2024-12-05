// (para el manejo de rutas/endpoints)

// Backend/models/handlers/auth.go
/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 05/12/2024
Lenguaje: Golang
Descipcion: Asignacion de la clase auth, con sus respectivas
funciones para el manejo de rutas
*/

// handlers/auth.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type AuthHandler struct {
	db *sql.DB
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
	ExpireAt string `json:"expire_at"`
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error ingreso de datos", http.StatusBadRequest)
		return
	}

	// Verificar credenciales en la base de datos
	var user LoginResponse
	err := h.db.QueryRow(
		"SELECT id, name, email, role FROM users WHERE email = ? AND password = ?",
		req.Email, req.Password,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err == sql.ErrNoRows {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Generar token
	user.Token = "token-ejemplo"
	user.ExpireAt = time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Sesión cerrada exitosamente"})
}
