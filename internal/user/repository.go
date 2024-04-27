package user

import (
	"context" // Paquete `context`: Este paquete proporciona un objeto de contexto que lleva información del ámbito de la solicitud.
	"errors"
	"log" // Paquete `log`: Este paquete proporciona funciones para registrar mensajes.
	"slices"

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/domain" // Paquete `internal/domain`: Este paquete (probablemente ubicado en el directorio `internal`) proporciona la estructura `User` utilizada para representar datos de usuario.
)

// DB estructura que contiene los datos de usuario y un contador para el ID máximo.
type DB struct {
	Users     []domain.User // Lista de usuarios: Representa una lista de estructuras `domain.User` para almacenar los datos de los usuarios.
	MaxUserID uint64        // ID máximo para generar IDs automáticos: Un entero sin signo de 64 bits para mantener un registro del ID de usuario máximo para la generación automática de ID.
}

// Define las operaciones básicas que debe implementar un repositorio de usuarios, como obtener, crear, actualizar y eliminar usuarios.
type Repository interface {
	Create(ctx context.Context, user *domain.User) error // Método para crear un nuevo usuario: Define un método llamado `Create` que recibe un contexto y un puntero a una estructura `domain.User` y devuelve un error.
	GetAll(ctx context.Context) ([]domain.User, error)   // Método para obtener todos los usuarios: Define un método llamado `GetAll` que recibe un contexto y devuelve una lista de estructuras `domain.User` y un error.
	Get(ctx context.Context, id uint64) (*domain.User, error)
}

// repo es una implementación de la interfaz Repository.
type repo struct {
	db  DB          // Base de datos de usuarios
	log *log.Logger // Logger para registrar eventos
}

// NewRepo es una función constructora que devuelve una nueva instancia del repositorio.
// Crea un nuevo repositorio de usuarios utilizando la base de datos simulada y el logger.
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

func (r *repo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
		return v.ID == id
	})

	if index < 0 {
		return nil, errors.New("user not found")
	}

	r.log.Println("repository get")
	return &r.db.Users[index], nil
}

/*
Capa de repositorio (Repository):

Función: Esta capa se encarga de acceder y manipular los datos de la aplicación.
Abstrae la fuente de datos: Oculta la implementación específica de la base de datos o el sistema de almacenamiento.
Proporciona una interfaz: Define una interfaz para acceder a los datos de manera independiente de la tecnología subyacente.
Realiza operaciones CRUD: Implementa las operaciones básicas de creación, lectura, actualización y eliminación (CRUD) sobre los datos.
Interactúa con la fuente de datos: Utiliza controladores específicos para conectarse a la base de datos o el sistema de almacenamiento.
*/
