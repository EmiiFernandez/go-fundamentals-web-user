package user

import (
	"context"
	"log"

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/domain"
)

// Define la interfaz del servicio de usuarios
// Define las operaciones, como crear, obtener, actualizar y eliminar usuarios.
type Service interface {
	// Crea un nuevo usuario
	//El contexto se utiliza para pasar información adicional a la función Create
	Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error)

	// Obtiene todos los usuarios
	GetAll(ctx context.Context) ([]domain.User, error)
}

// Implementación del servicio de usuarios
type service struct {
	log  *log.Logger // Instancia del logger para registrar mensajes
	repo Repository  // Instancia del repositorio de usuarios
}

// Función constructora del servicio
// Crea un nuevo servicio de usuarios utilizando el logger y el repositorio.
func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

// Implementación del método Create
func (s *service) Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error) {
	// Crea una nueva instancia de domain.User con los datos proporcionados
	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	// Delega la creación del usuario al repositorio
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Registra un mensaje en el logger indicando la creación del usuario
	s.log.Println("Usuario creado correctamente")

	// Retorna la instancia del usuario creado
	return user, nil
}

// Implementación del método GetAll
func (s *service) GetAll(ctx context.Context) ([]domain.User, error) {
	// Delega la obtención de todos los usuarios al repositorio
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Registra un mensaje en el logger indicando la obtención de todos los usuarios
	s.log.Println("Se han obtenido todos los usuarios")

	// Retorna la lista de usuarios obtenida del repositorio
	return users, nil
}

/*
Capa de servicio (Service):

Función: Esta capa implementa la lógica de negocio de la aplicación.
Realiza las operaciones: Ejecuta las tareas de negocio relacionadas con los datos, como crear, obtener, actualizar o eliminar información.
Aplica las reglas: Valida las entradas, aplica reglas de negocio y procesa los datos.
Interactúa con la capa de repositorio: Accede a los datos a través de la capa de repositorio.
Encapsula la lógica: Oculta los detalles de implementación de la lógica de negocio a la capa de presentación.
*/
