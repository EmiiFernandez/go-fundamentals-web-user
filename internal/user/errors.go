package user

import (
	"errors"
	"fmt"
)

// ErrFirstNameRequired se produce cuando se intenta crear un usuario sin proporcionar un nombre.
var ErrFirstNameRequired = errors.New("first name ir required")

// ErrLastNameRequired se produce cuando se intenta crear un usuario sin proporcionar un apellido.
var ErrLastNameRequired = errors.New("last name ir required")

// ErrThereArentFields se utiliza cuando no se proporcionan campos para actualizar en la función Update del repositorio de usuarios.
var ErrThereArentFields = errors.New("there aren't fields")

// ErrNotFound es una estructura de error personalizada que se utiliza cuando no se encuentra un usuario en la base de datos.
type ErrNotFound struct {
	ID uint64 // ID del usuario que no se encontró.
}

// Error implementa el método Error de la interfaz error para la estructura ErrNotFound.
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user id '%d' doesn`t exist", e.ID) // Retorna un mensaje de error formateado con el ID del usuario.
}
