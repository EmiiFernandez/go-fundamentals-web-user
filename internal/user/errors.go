package user

import (
	"errors"
	"fmt"
)

// ErrFirstNameRequired es un error que se produce cuando se intenta crear un usuario sin proporcionar un nombre.
var ErrFirstNameRequired = errors.New("first name ir required")

// ErrLastNameRequired es un error que se produce cuando se intenta crear un usuario sin proporcionar un apellido.
var ErrLastNameRequired = errors.New("last name ir required")

// ErrNotFound es una estructura de error personalizada que se utiliza cuando no se encuentra un usuario en la base de datos.
type ErrNotFound struct {
	ID uint64 // ID del usuario que no se encontró.
}

// Error implementa el método Error de la interfaz error para la estructura ErrNotFound.
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user id '%d' doesn`t exist", e.ID) // Retorna un mensaje de error formateado con el ID del usuario.
}
