// Backend/Handlers/admin.go
/* Autores: Henry Aliaga / Ismael Espinoza
Fecha: 06/12/2024
Lenguaje: Golang
Descripción: Funciones específicas para administradores.
*/

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// Listar todos los usuarios
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.db.Query("SELECT id, name, email, role FROM users")
	if err != nil {
		http.Error(w, "Error al obtener la lista de usuarios", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			http.Error(w, "Error al procesar los usuarios", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Actualizar el rol de un usuario
func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		UserID int    `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al leer el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("UPDATE users SET role = ? WHERE id = ?", input.Role, input.UserID)
	if err != nil {
		http.Error(w, "Error al actualizar el rol del usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rol actualizado correctamente"))
}

// Eliminar un usuario
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usuario eliminado correctamente"))
}
