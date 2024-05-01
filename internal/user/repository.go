package user

import (
	"context" // Paquete `context`: Proporciona un objeto de contexto que lleva información del ámbito de la solicitud.
	"database/sql"
	"log" // Paquete `log`: Proporciona funciones para registrar mensajes.

	// Paquete `slices`: Proporciona funciones para trabajar con slices.
	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/domain" // Paquete `internal/domain`: Proporciona la estructura `User` utilizada para representar datos de usuario.
)

/*// DB estructura que contiene los datos de usuario y un contador para el ID máximo.
type DB struct {
	Users     []domain.User // Lista de usuarios: Representa una lista de estructuras `domain.User` para almacenar los datos de los usuarios.
	MaxUserID uint64        // ID máximo para generar IDs automáticos: Un entero sin signo de 64 bits para mantener un registro del ID de usuario máximo para la generación automática de ID.
}*/

// Repository define las operaciones básicas que debe implementar un repositorio de usuarios.
type Repository interface {
	// Create crea un nuevo usuario en la base de datos.
	Create(ctx context.Context, user *domain.User) error
	// GetAll devuelve todos los usuarios almacenados en la base de datos.
	GetAll(ctx context.Context) ([]domain.User, error)
	// Get devuelve un usuario específico basado en su ID.
	Get(ctx context.Context, id uint64) (*domain.User, error)
	// Update actualiza los datos de un usuario existente.
	Update(ctx context.Context, id uint64, firstName, lastName, email *string) error
}

// repo es una implementación de la interfaz Repository.
type repo struct {
	db  *sql.DB     // Base de datos de usuarios
	log *log.Logger // Logger para registrar eventos
}

// NewRepo es una función constructora que devuelve una nueva instancia del repositorio.
func NewRepo(db *sql.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

// Create crea un nuevo usuario en la base de datos.
/*
	BD en memoria
	r.db.MaxUserID++                       // Incrementar el ID máximo
	user.ID = r.db.MaxUserID               // Asignar el nuevo ID al usuario
	r.db.Users = append(r.db.Users, *user) // Agregar el usuario a la lista de usuarios en la base de datos
	r.log.Println("repository create")     // Registrar en el logger que se ha creado un nuevo usuario
*/
func (r *repo) Create(ctx context.Context, user *domain.User) error {
	sqlQ := "INSERT INTO users(first_name, last_name, email) VALUES(?,?,?)"
	//id autoincremental
	res, err := r.db.Exec(sqlQ, user.FirstName, user.LastName, user.Email)
	if err != nil {
		r.log.Println(err.Error())
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	user.ID = uint64(id)
	r.log.Println("user created with id: ", id)
	//r.log.Println("repository create")
	return nil
}

// GetAll devuelve todos los usuarios almacenados en la base de datos.
func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {

	var users []domain.User
	sqlQ := "SELECT id, first_name, last_name, email FROM users"
	rows, err := r.db.Query(sqlQ)
	if err != nil {
		r.log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {

			r.log.Println(err.Error())
			return nil, err
		}
		users = append(users, u)
	}

	//r.log.Println("repository get all")
	r.log.Println("user get all: ", len(users)) // Registrar en el logger que se está obteniendo la lista de usuarios
	return users, nil                           // Devolver la lista de usuarios
}

// Get devuelve un usuario específico basado en su ID.
func (r *repo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	/*// Buscar el usuario en la lista de usuarios por su ID
	index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
		return v.ID == id
	})

	// Si no se encuentra el usuario, devolver un error
	if index < 0 {
		return nil, ErrNotFound{id}
	}

	r.log.Println("repository get") // Registrar en el logger que se ha obtenido un usuario*/
	return nil, nil // Devolver el usuario encontrado
}

// Update actualiza los datos de un usuario existente.
func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	/*user, err := r.Get(ctx, id) // Obtener el usuario existente
	if err != nil {
		return err // Devolver el error si el usuario no existe
	}

	// Actualizar los campos del usuario con los nuevos valores si se proporcionan
	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if email != nil {
		user.Email = *email
	}

	r.log.Println("repository update") // Registrar en el logger que se ha actualizado un usuario*/
	return nil
}

/*
Capa de repositorio (Repository):

Función: La capa de repositorio se encarga de acceder y manipular los datos de la aplicación.
Abstracción de la fuente de datos: Oculta la implementación específica de la base de datos o el sistema de almacenamiento.
Interfaz definida: Define una interfaz para acceder a los datos de manera independiente de la tecnología subyacente.
Operaciones CRUD: Implementa las operaciones básicas de creación, lectura, actualización y eliminación (CRUD) sobre los datos.
Interacción con la fuente de datos: Utiliza controladores específicos para conectarse a la base de datos o el sistema de almacenamiento.
*/
