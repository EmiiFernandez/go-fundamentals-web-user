package domain

//Define la estructura de datos para representar un usuario
type User struct {
	ID uint64 `json:"id"` // Identificador único del usuario

	FirstName string `json:"first_name"` // Nombre del usuario

	LastName string `json:"last_name"` // Apellido del usuario

	Email string `json:"email"` // Dirección de correo electrónico del usuario
}
