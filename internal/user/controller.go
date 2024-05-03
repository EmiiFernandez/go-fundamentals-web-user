package user

import (
	"context" // El paquete `context` proporciona un objeto de contexto para llevar información del ámbito de la solicitud.
	"errors"
	"fmt" // El paquete `fmt` proporciona funciones para el formateo de salida de datos.

	"github.com/EmiiFernandez/go-fundamentals-response/response"
)

// Definición de tipos

type (
	// Controller: Define un tipo alias para una función que recibe un objeto `http.ResponseWriter` para escribir la respuesta y un objeto `http.Request` para leer la solicitud del cliente.
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	// Endpoints: Define una estructura `Endpoints` que agrupa los controladores para los endpoints (rutas) de la API.
	Endpoints struct {
		Create Controller // Campo `Create` de tipo `Controller` que almacena el controlador para el endpoint de creación de usuarios.
		GetAll Controller // Campo `GetAll` de tipo `Controller` que almacena el controlador para el endpoint de obtención de todos los usuarios.
		Get    Controller // Campo `Get` de tipo `Controller` que almacena el controlador para el endpoint de obtención de un usuario por ID.
		Update Controller // Campo `Update` de tipo `Controller` que almacena el controlador para el endpoint de actualización de un usuario por ID.
		Delete Controller // // Campo `Delete` de tipo `Controller` que almacena el controlador para el endpoint de eliminación de un usuario por ID.
	}

	GetReq struct {
		ID uint64 // ID del usuario a obtener
	}

	// CreateReq: Define una estructura `CreateReq` para representar la solicitud de creación de un nuevo usuario.
	CreateReq struct {
		FirstName string `json:"first_name"` // Campo `FirstName` de tipo cadena para almacenar el nombre del usuario.
		LastName  string `json:"last_name"`  // Campo `LastName` de tipo cadena para almacenar el apellido del usuario.
		Email     string `json:"email"`      // Campo `Email` de tipo cadena para almacenar el correo electrónico del usuario.
		// La etiqueta `json:"first_name"` indica la clave que se usará al codificar el campo a JSON.
	}

	// UpdateReq: Define una estructura `UpdateReq` para representar la solicitud de actualización de un usuario.
	UpdateReq struct {
		ID        uint64  // ID del usuario a actualizar
		FirstName *string `json:"first_name"` // Campo `FirstName` de tipo puntero a cadena para almacenar el nombre del usuario.
		LastName  *string `json:"last_name"`  // Campo `LastName` de tipo puntero a cadena para almacenar el apellido del usuario.
		Email     *string `json:"email"`      // Campo `Email` de tipo puntero a cadena para almacenar el correo electrónico del usuario.
	}

	// DeleteReq: Define una estructura `DeleteReq` para representar la solicitud de eliminación de un usuario.
	DeleteReq struct {
		ID uint64 // ID del usuario a eliminar
	}
)

// Funciones del controlador

// MakeEndpoints crea los endpoints (rutas) de la API y asigna los controladores correspondientes.
func MakeEndpoints(ctx context.Context, s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Get:    makeGetEndopoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

// makeCreateEndpoint crea un controlador para el endpoint de creación de usuarios.
func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Esta función maneja las solicitudes POST para crear un nuevo usuario en la API de usuarios.

		// Convierte la interfaz `data` a la estructura `CreateReq` para acceder a los campos del usuario.
		req := request.(CreateReq)

		// Valida los campos obligatorios de la solicitud (nombre, apellido y correo electrónico).
		if req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}
		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		// Llama a la función `Create` del servicio `Service` para crear el nuevo usuario.
		// Esta función (que probablemente se encuentre en otro paquete) se encarga de la lógica de negocio para persistir el usuario en un repositorio.
		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)

		// Maneja el error en caso de que falle la creación del usuario.
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.Created("success", user), nil
	}
}

// makeGetAllEndpoint crea un controlador para el endpoint de obtención de todos los usuarios.
func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Esta función recupera todos los usuarios del servicio y los envía como respuesta.

		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", users), nil
	}
}

// makeGetEndopoint crea un controlador para el endpoint de obtención de un usuario por ID.
func makeGetEndopoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Esta función maneja las solicitudes GET para obtener un usuario por su ID.

		// Convierte la interfaz `data` a la estructura `GetReq` para acceder al ID del usuario.
		req := request.(GetReq)

		// Llama a la función `Get` del servicio `Service` para obtener el usuario por su ID.
		user, err := s.Get(ctx, req.ID)

		// Maneja el error en caso de que falle la obtención del usuario.
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		fmt.Println(req)
		return response.OK("success", user), nil
	}
}

// makeUpdateEndpoint crea un controlador para el endpoint de actualización de un usuario por ID.
func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Esta función maneja las solicitudes PATCH para actualizar un usuario por su ID.

		// Convierte la interfaz `data` a la estructura `UpdateReq` para acceder a los campos de actualización del usuario.
		req := request.(UpdateReq)

		// Valida los campos obligatorios de la solicitud (nombre y apellido).
		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}
		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		// Llama a la función `Update` del servicio `Service` para actualizar los datos del usuario.
		if err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email); err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		fmt.Println(req)
		return response.OK("success", nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Verifica si la solicitud tiene el formato esperado
		req, ok := request.(DeleteReq)
		if !ok {
			// Si la solicitud no tiene el formato esperado, devuelve un error
			return nil, errors.New("invalid request format")
		}

		// Llama al método Delete del servicio para eliminar el usuario
		_, err := s.Delete(ctx, req.ID)
		if err != nil {
			// Maneja el error en caso de que falle la eliminación del usuario
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		// Devuelve un mensaje de éxito indicando que el usuario fue eliminado
		return response.OK("user deleted successfully", nil), nil
	}
}

/*
Capa de presentación (Controller):

Función: La capa de presentación interactúa directamente con el usuario o cliente.
Maneja las solicitudes HTTP: Recibe las solicitudes del cliente, las analiza y las valida.
Decodifica los datos: Convierte los datos recibidos en el formato adecuado para su procesamiento.
Llama a la capa de servicio: Delega la lógica de negocio a la capa de servicio.
Codifica las respuestas: Convierte los resultados de la capa de servicio en formato JSON o HTML para enviar al cliente.
Envía las respuestas: Envía las respuestas al cliente junto con el código de estado HTTP correspondiente.
*/
