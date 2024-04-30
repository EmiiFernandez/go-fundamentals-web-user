package bootstrap

/*
Package bootstrap proporciona funciones para inicializar componentes clave antes de que la aplicación principal comience a ejecutarse.
*/

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // _ lo importo pero no lo uso
)

// NewLogger crea y devuelve un objeto *log.Logger que se utiliza para registrar mensajes en la consola.
func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

// NewBD inicializa y devuelve una conexión a la base de datos MySQL.
func NewBD() (*sql.DB, error) {
	// Abre una conexión a la base de datos MySQL utilizando las credenciales y la cadena de conexión proporcionadas.
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3336)/go_course_users")
	if err != nil {
		return nil, err // Devuelve nil y el error si la conexión falla
	}

	return db, nil // Devuelve la conexión a la base de datos y ningún error si es exitosa
}

/*
return user.DB{
		Users: []domain.User{
			{ID: 1, FirstName: "Nahuel", LastName: "Costamagna", Email: "nahuel@domain.com"},
			{ID: 2, FirstName: "Eren", LastName: "Jaeger", Email: "eren@domain.com"},
			{ID: 3, FirstName: "Paco", LastName: "Costa", Email: "paco@domain.com"},
		},
		MaxUserID: 3,
	}
*/
