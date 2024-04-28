package user

import (
	"context" // Paquete `context`: Proporciona un objeto de contexto que lleva información del ámbito de la solicitud.
	"errors"
	"log" // Paquete `log`: Proporciona funciones para registrar mensajes.
	"slices"

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/domain" // Paquete `internal/domain`: Proporciona la estructura `User` utilizada para representar datos de usuario.
)

// DB estructura que contiene los datos de usuario y un contador para el ID máximo.
type DB struct {
	Users     []domain.User // Lista de usuarios: Representa una lista de estructuras `domain.User` para almacenar los datos de los usuarios.
	MaxUserID uint64        // ID máximo para generar IDs automáticos: Un entero sin signo de 64 bits para mantener un registro del ID de usuario máximo para la generación automática de ID.
}

// Repository define las operaciones básicas que debe implementar un repositorio de usuarios.
type Repository interface {
	// Create crea un nuevo usuario en la base de datos.
	Create(ctx context.Context, user *domain.User) error
	// GetAll devuelve todos los usuarios almacenados en la base de datos.
	GetAll(ctx context.Context) ([]domain.User, error)
	// Get devuelve un usuario específico basado en su ID.
	Get(ctx context.Context, id uint64) (*domain.User, error)
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

// Get devuelve un usuario específico basado en su ID.
func (r *repo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	// Buscar el usuario en la lista de usuarios por su ID
	index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
		return v.ID == id
	})

	// Si no se encuentra el usuario, devolver un error
	if index < 0 {
		return nil, errors.New("user not found")
	}

	r.log.Println("repository get") // Registrar en el logger que se ha obtenido un usuario
	return &r.db.Users[index], nil  // Devolver el usuario encontrado
}

/*
Capa de repositorio (Repository):

Función: La capa de repositorio se encarga de acceder y manipular los datos de la aplicación.
Abstracción de la fuente de datos: Oculta la implementación específica de la base de datos o el sistema de almacenamiento.
Interfaz definida: Define una interfaz para acceder a los datos de manera independiente de la tecnología subyacente.
Operaciones CRUD: Implementa las operaciones básicas de creación, lectura, actualización y eliminación (CRUD) sobre los datos.
Interacción con la fuente de datos: Utiliza controladores específicos para conectarse a la base de datos o el sistema de almacenamiento.
*/
