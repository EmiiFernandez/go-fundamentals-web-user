package user

import (
	"context" // Paquete `context`: Proporciona un objeto de contexto que lleva información del ámbito de la solicitud.
	"log"     // Paquete `log`: Proporciona funciones para registrar mensajes.

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/domain" // Paquete `internal/domain`: Proporciona la estructura `User` utilizada para representar datos de usuario.
)

// Service define la interfaz del servicio de usuarios.
type Service interface {
	// Create crea un nuevo usuario con los datos proporcionados.
	// El contexto se utiliza para pasar información adicional a la función Create.
	Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error)

	// GetAll devuelve todos los usuarios almacenados en la base de datos.
	GetAll(ctx context.Context) ([]domain.User, error)

	// Get devuelve un usuario específico basado en su ID.
	Get(ctx context.Context, id uint64) (*domain.User, error)
}

// service es una implementación del servicio de usuarios.
type service struct {
	log  *log.Logger // Instancia del logger para registrar mensajes.
	repo Repository  // Instancia del repositorio de usuarios.
}

// NewService es una función constructora que devuelve una nueva instancia del servicio de usuarios.
// Crea un nuevo servicio de usuarios utilizando el logger y el repositorio.
func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

// Create crea un nuevo usuario con los datos proporcionados.
func (s *service) Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error) {
	// Crea una nueva instancia de domain.User con los datos proporcionados.
	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	// Delega la creación del usuario al repositorio.
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Registra un mensaje en el logger indicando la creación del usuario.
	s.log.Println("Usuario creado correctamente")

	// Retorna la instancia del usuario creado.
	return user, nil
}

// GetAll devuelve todos los usuarios almacenados en la base de datos.
func (s *service) GetAll(ctx context.Context) ([]domain.User, error) {
	// Delega la obtención de todos los usuarios al repositorio.
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Registra un mensaje en el logger indicando la obtención de todos los usuarios.
	s.log.Println("Se han obtenido todos los usuarios")

	// Retorna la lista de usuarios obtenida del repositorio.
	return users, nil
}

// Get devuelve un usuario específico basado en su ID.
func (s *service) Get(ctx context.Context, id uint64) (*domain.User, error) {
	// Delega la obtención del usuario al repositorio.
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Registra un mensaje en el logger indicando la obtención del usuario.
	s.log.Println("Se ha obtenido el usuario seleccionado")

	// Retorna el usuario obtenido del repositorio.
	return user, nil
}

/*
Capa de servicio (Service):

Función: La capa de servicio implementa la lógica de negocio de la aplicación.
Realiza las operaciones: Ejecuta las tareas de negocio relacionadas con los datos, como crear, obtener, actualizar o eliminar información.
Aplica las reglas: Valida las entradas, aplica reglas de negocio y procesa los datos.
Interactúa con la capa de repositorio: Accede a los datos a través de la capa de repositorio.
Encapsula la lógica: Oculta los detalles de implementación de la lógica de negocio a la capa de presentación.
*/
