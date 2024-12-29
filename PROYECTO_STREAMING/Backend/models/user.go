/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 22/11/2024
Lenguaje: Golang
Descipcion: Asignacion de la clase user, con sus respectivas
funciones para el manejo de datos
(para la estructura de datos)*/

package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

// Usuario representa a los consumidores.
type Usuario struct {
	ID       int    `json:"user_id"`       // Identificación única de usuario
	Nombre   string `json:"user_name"`     // Nombre del usuario
	Email    string `json:"user_email"`    // Correo electrónico del usuario
	Password string `json:"user_password"` // Contraseña del usuario
	Token    string `json:"token"`         // Token generado para el usuario
	ExpireAt string `json:"expire_at"`     // Fecha de expiración del token
	Role     string `json:"user_role"`     // Rol del usuario (nuevo campo)
}

// NewUsuario crea una nueva instancia de Usuario con validación de campos.
func NewUsuario(id int, nombre, email, password string) (*Usuario, error) {
	// Validar que los campos requeridos no estén vacíos
	if nombre == "" || email == "" || password == "" {
		return nil, fmt.Errorf("todos los campos son requeridos: nombre, email y contraseña")
	}

	// Crear el objeto Usuario
	usuario := &Usuario{
		ID:       id,
		Nombre:   nombre,
		Email:    email,
		Password: password,
	}

	// Generar un token único para el usuario y establecer expiración
	usuario.Token = generarTokenUnico(email)
	usuario.ExpireAt = time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	return usuario, nil
}

// Lista simulada de usuarios registrados
var UsuariosRegistrados = []Usuario{
	{ID: 1, Nombre: "Juan Perez", Email: "juan.perez@example.com", Password: "password123"},
	{ID: 2, Nombre: "Ana Gomez", Email: "ana.gomez@example.com", Password: "securepass456"},
	{ID: 3, Nombre: "Carlos Lopez", Email: "carlos.lopez@example.com", Password: "qwerty789"},
}

// ValidarORegistrarUsuario verifica si un usuario está registrado por su nombre.
// Si el usuario no está registrado, lo registra en la lista.
func ValidarORegistrarUsuario(id int, nombre, email, password string) (*Usuario, error) {
	// Buscar si el usuario ya existe por nombre
	for _, usuario := range UsuariosRegistrados {
		if usuario.Nombre == nombre {
			return &usuario, nil // Usuario encontrado
		}
	}

	// Usuario no encontrado, proceder a registrarlo
	nuevoUsuario, err := NewUsuario(id, nombre, email, password)
	if err != nil {
		return nil, fmt.Errorf("error al crear el usuario: %v", err)
	}

	// Agregar a la lista de usuarios registrados
	UsuariosRegistrados = append(UsuariosRegistrados, *nuevoUsuario)

	return nuevoUsuario, nil
}

// ObtenerUsuarios devuelve la lista de usuarios registrados.
func ObtenerUsuarios() []Usuario {
	return UsuariosRegistrados
}

// generarTokenUnico genera un token único y seguro para un usuario basado en su email.
func generarTokenUnico(email string) string {
	// Generar un token aleatorio
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	randomToken := base64.URLEncoding.EncodeToString(tokenBytes)

	return fmt.Sprintf("%s-%s", randomToken, email)
}

// ValidarToken verifica si un token es válido y devuelve el usuario correspondiente.
func ValidarToken(token string) (*Usuario, error) {
	for _, usuario := range UsuariosRegistrados {
		if usuario.Token == token {
			// Verificar si el token está expirado
			expireAt, err := time.Parse(time.RFC3339, usuario.ExpireAt)
			if err != nil {
				return nil, fmt.Errorf("formato de fecha inválido para expiración del token")
			}
			if time.Now().After(expireAt) {
				return nil, fmt.Errorf("el token ha expirado")
			}
			return &usuario, nil
		}
	}
	return nil, fmt.Errorf("token inválido")
}
