package bootstrap

import (
	"log"
	"os"

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/domain"
	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/user"
)

// Crea y devuelve un objeto *log.Logger que se utiliza para registrar mensajes en la consola.
func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

// Crea y devuelve una instancia de user.DB que representa la base de datos de usuarios.
func NewBD() user.DB {
	return user.DB{
		Users: []domain.User{
			{ID: 1, FirstName: "Nahuel", LastName: "Costamagna", Email: "nahuel@domain.com"},
			{ID: 2, FirstName: "Eren", LastName: "Jaeger", Email: "eren@domain.com"},
			{ID: 3, FirstName: "Paco", LastName: "Costa", Email: "paco@domain.com"},
		},
		MaxUserID: 3,
	}
}
