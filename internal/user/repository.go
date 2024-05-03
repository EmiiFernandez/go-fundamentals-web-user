package user

import (
	"context" // Paquete `context`: Proporciona un objeto de contexto que lleva información del ámbito de la solicitud.
	"database/sql"
	"fmt"
	"log" // Paquete `log`: Proporciona funciones para registrar mensajes.
	"strings"

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
	// Elimina un usuario específico basado en su ID.
	Delete(ctx context.Context, id uint64) (*domain.User, error)
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

/*
	BD en memoria
	r.db.MaxUserID++                       // Incrementar el ID máximo
	user.ID = r.db.MaxUserID               // Asignar el nuevo ID al usuario
	r.db.Users = append(r.db.Users, *user) // Agregar el usuario a la lista de usuarios en la base de datos
	r.log.Println("repository create")     // Registrar en el logger que se ha creado un nuevo usuario
*/
// Create crea un nuevo usuario en la base de datos.
func (r *repo) Create(ctx context.Context, user *domain.User) error {
	// Query SQL para insertar un nuevo usuario en la base de datos.
	sqlQ := "INSERT INTO users(first_name, last_name, email) VALUES(?,?,?)"
	// Ejecutar la consulta SQL y obtener el resultado.
	res, err := r.db.Exec(sqlQ, user.FirstName, user.LastName, user.Email)
	if err != nil {
		// Si ocurre un error al ejecutar la consulta, registrar el error y devolverlo.
		r.log.Println(err.Error())
		return err
	}
	// Obtener el ID del usuario recién creado.
	id, err := res.LastInsertId()
	if err != nil {
		r.log.Println(err.Error())
		return err
	}
	// Asignar el ID al usuario y registrar el éxito en el log.
	user.ID = uint64(id)
	r.log.Println("user created with id: ", id)
	return nil
}

// GetAll devuelve todos los usuarios almacenados en la base de datos.
func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	// Consulta SQL para obtener todos los usuarios.
	sqlQ := "SELECT id, first_name, last_name, email FROM users"
	// Ejecutar la consulta SQL.
	rows, err := r.db.Query(sqlQ)
	if err != nil {
		r.log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Iterar sobre los resultados y almacenar los usuarios en un slice.
	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
			r.log.Println(err.Error())
			return nil, err
		}
		users = append(users, u)
	}

	// Registrar la cantidad de usuarios obtenidos en el log y devolver el slice de usuarios.
	r.log.Println("user get all: ", len(users))
	return users, nil
}

/* Usado en BD en memoria
// Buscar el usuario en la lista de usuarios por su ID
index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
	return v.ID == id
})

// Si no se encuentra el usuario, devolver un error
if index < 0 {
	return nil, ErrNotFound{id}
}

r.log.Println("repository get") // Registrar en el logger que se ha obtenido un usuario*/
// Get devuelve un usuario específico basado en su ID.
func (r *repo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	// Consulta SQL para obtener un usuario por su ID.
	sqlQ := "SELECT id, first_name, last_name, email FROM users WHERE id = ?"
	// Variables para almacenar el usuario y los posibles errores.
	var u domain.User
	// Ejecutar la consulta SQL y escanear el resultado en la estructura del usuario.
	if err := r.db.QueryRow(sqlQ, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
		// Si no se encuentra el usuario, devolver un error NotFound.
		r.log.Println(err.Error())
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{id}
		}
		return nil, err
	}
	// Registrar el éxito en el log y devolver el usuario.
	r.log.Println("get user with id: ", id)
	return &u, nil
}

/*BDD por memoria

user, err := r.Get(ctx, id) // Obtener el usuario existente
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

// Update actualiza los datos de un usuario existente en la base de datos.
func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	// Construir la lista de campos a actualizar y los valores correspondientes.
	var fields []string
	var values []interface{}

	// Verificar si se proporciona un nuevo nombre y agregarlo a los campos a actualizar.
	if firstName != nil {
		fields = append(fields, "first_name=?")
		values = append(values, *firstName)
	}

	// Verificar si se proporciona un nuevo apellido y agregarlo a los campos a actualizar.
	if lastName != nil {
		fields = append(fields, "last_name=?")
		values = append(values, *lastName)
	}

	// Verificar si se proporciona un nuevo correo electrónico y agregarlo a los campos a actualizar.
	if email != nil {
		fields = append(fields, "email=?")
		values = append(values, *email)
	}

	// Verificar si no se proporciona ningún campo para actualizar.
	if len(fields) == 0 {
		r.log.Println(ErrThereArentFields.Error())
		return ErrThereArentFields
	}

	// Agregar el ID del usuario a los valores para la consulta SQL.
	values = append(values, id)

	// Construir la consulta SQL final con los campos a actualizar.
	sqlQ := fmt.Sprintf("UPDATE users SET %s WHERE id=?", strings.Join(fields, ","))
	// Ejecutar la consulta SQL con los valores correspondientes.
	res, err := r.db.Exec(sqlQ, values...)
	if err != nil {
		// Si ocurre un error al ejecutar la consulta, registrar el error y devolverlo.
		r.log.Println(err.Error())
		return err
	}

	// Verificar si se actualizó correctamente algún registro.
	row, err := res.RowsAffected()
	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	// Si no se actualizó ningún registro, devolver un error NotFound.
	if row == 0 {
		err := ErrNotFound{id}
		r.log.Println(err.Error())
		return err
	}

	// Registrar el éxito en el log y devolver nil (sin error).
	r.log.Println("user updated id: ", id)
	return nil
}

// Delete elimina los datos de un usuario existente en la base de datos.
func (r *repo) Delete(ctx context.Context, id uint64) (*domain.User, error) {
	// Consulta SQL para eliminar un usuario por su ID
	sqlQ := "DELETE FROM users WHERE id = ?"

	// Ejecutar la consulta SQL para eliminar el usuario
	result, err := r.db.Exec(sqlQ, id)
	if err != nil {
		r.log.Println(err.Error())
		return nil, err
	}

	// Verificar si se eliminó algún registro
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Println(err.Error())
		return nil, err
	}
	if rowsAffected == 0 {
		// Si no se encontró ningún usuario para eliminar, devuelve un error
		return nil, ErrNotFound{id}

	}

	// Registrar el éxito en el log y devolver el usuario eliminado
	r.log.Println("Usuario eliminado con ID:", id)
	return &domain.User{ID: id}, nil
}

/*
Capa de repositorio (Repository):

Función: La capa de repositorio se encarga de acceder y manipular los datos de la aplicación.
Abstracción de la fuente de datos: Oculta la implementación específica de la base de datos o el sistema de almacenamiento.
Interfaz definida: Define una interfaz para acceder a los datos de manera independiente de la tecnología subyacente.
Operaciones CRUD: Implementa las operaciones básicas de creación, lectura, actualización y eliminación (CRUD) sobre los datos.
Interacción con la fuente de datos: Utiliza controladores específicos para conectarse a la base de datos o el sistema de almacenamiento.
*/
