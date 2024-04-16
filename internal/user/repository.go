package user

import (
	"context"
	"log"

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/domain"
)

// DB estructura que contiene los datos de usuario y un contador para el ID máximo.
type DB struct {
	Users     []domain.User // Lista de usuarios
	MaxUserID uint64        // ID máximo para generar IDs automáticos
}

// Repository es una interfaz que define los métodos para interactuar con el repositorio de usuarios.
type Repository interface {
	Create(ctx context.Context, user *domain.User) error // Método para crear un nuevo usuario
	GetAll(ctx context.Context) ([]domain.User, error)   // Método para obtener todos los usuarios
}

// repo es una implementación de la interfaz Repository.
type repo struct {
	db  DB          // Base de datos de usuarios
	log *log.Logger // Logger para registrar eventos
}

// NewRepo es una función constructora que devuelve una nueva instancia del repositorio.
func NewRepo(db DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

// Create crea un nuevo usuario en la base de datos.
func (r *repo) Create(ctx context.Context, user *domain.User) error {
	r.db.MaxUserID++                       // Incrementar el ID máximo
	user.ID = r.db.MaxUserID               // Asignar el nuevo ID al usuario
	r.db.Users = append(r.db.Users, *user) // Agregar el usuario a la lista de usuarios en la base de datos
	r.log.Println("repository create")     // Registrar en el logger que se ha creado un nuevo usuario
	return nil
}

// GetAll devuelve todos los usuarios almacenados en la base de datos.
func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	r.log.Println("repository get all") // Registrar en el logger que se está obteniendo la lista de usuarios
	return r.db.Users, nil              // Devolver la lista de usuarios
}
