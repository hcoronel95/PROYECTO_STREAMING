// Backend/models/user.go
/*Autores: Henry Aliaga / Ismael Espinoza
Fecha: 22/11/2024
Lenguaje: Golang
Descipcion: Asignacion de la clase admin, con sus respectivas
funciones para el manejo de datos
(para la estructura de datos)*/

// (para la estructura de datos)

package models

import (
	"errors"
	"fmt"
)

// El objeto Admin representa a los administradores.
type Admin struct {
	ID       int    `json:"admin_id"`       // Identificación única del administrador
	Nombre   string `json:"user_name"`      // Nombre del administrador
	Email    string `json:"admin_email"`    // Correo electrónico del administrador
	Password string `json:"admin_password"` // Contraseña del administrador
}

// NewAdmin crea una nueva instancia de Admin con validación de campos.
func NewAdmin(id int, nombre, email, password string) (*Admin, error) {
	// Validar que los campos requeridos no estén vacíos
	if nombre == "" || email == "" || password == "" {
		return nil, fmt.Errorf("todos los campos son requeridos: nombre, email y contraseña")
	}

	// Crear el objeto Admin
	admin := &Admin{
		ID:       id,
		Nombre:   nombre,
		Email:    email,
		Password: password,
	}

	return admin, nil
}

// Lista simulada de administradores registrados
var AdminsRegistrados = []Admin{
	{ID: 1, Nombre: "HENRY ALIAGA", Email: "henry@example.com", Password: "admin123"},
	{ID: 2, Nombre: "ISMAEL ESPINOZA", Email: "ismael@example.com", Password: "admin123"},
}

// ValidarAdminPorNombre verifica si un administrador está registrado por su nombre.
// Si el administrador no está registrado, devuelve un error.
func ValidarAdminPorNombre(nombre string) (*Admin, error) {
	for _, admin := range AdminsRegistrados {
		if admin.Nombre == nombre {
			return &admin, nil // Administrador encontrado
		}
	}
	return nil, errors.New("el administrador no está registrado")
}

// ObtenerAdmins devuelve la lista de administradores registrados.
func ObtenerAdmins() []Admin {
	return AdminsRegistrados
}
