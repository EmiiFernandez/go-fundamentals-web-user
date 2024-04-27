package user

import (
	"context" // Paquete `context`: Proporciona un objeto de contexto para llevar información del ámbito de la solicitud.
	// Paquete `encoding/json`: Proporciona funciones para codificar y decodificar datos JSON.
	"errors"
	// Paquete `fmt`: Proporciona funciones para el formateo de salida de datos.
	// Paquete `net/http`: Proporciona funciones para crear servidores HTTP.
)

// Definición de tipos

type (
	// Controller: Define un tipo alias para una función que recibe un objeto `http.ResponseWriter` para escribir la respuesta y un objeto `http.Request` para leer la solicitud del cliente.
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	// Endpoints: Define una estructura `Endpoints` que agrupa los controladores para los endpoints (rutas) de la API.
	Endpoints struct {
		Create Controller // Campo `Create` de tipo `Controller` que almacena el controlador para el endpoint de creación de usuarios.
		GetAll Controller // Campo `GetAll` de tipo `Controller` que almacena el controlador para el endpoint de obtención de todos los usuarios.
	}

	// CreateReq: Define una estructura `CreateReq` para representar la solicitud de creación de un nuevo usuario.
	CreateReq struct {
		FirstName string `json:"first_name"` // Campo `FirstName` de tipo cadena para almacenar el nombre del usuario.
		LastName  string `json:"last_name"`  // Campo `LastName` de tipo cadena para almacenar el apellido del usuario.
		Email     string `json:"email"`      // Campo `Email` de tipo cadena para almacenar el correo electrónico del usuario.
		// La etiqueta `json:"first_name"` indica la clave que se usará al codificar el campo a JSON.
	}
)

// Funciones del controlador

func MakeEndpoints(ctx context.Context, s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		// Esta función recupera todos los usuarios del servicio y los envía como respuesta.
		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		return users, nil
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Esta función maneja las solicitudes POST para crear un nuevo usuario en la API de usuarios.

		// Convierte la interfaz `data` a la estructura `CreateReq` para acceder a los campos del usuario.
		req := request.(CreateReq)

		// Valida los campos obligatorios de la solicitud (nombre, apellido y correo electrónico).
		if req.FirstName == "" {
			return nil, errors.New("first name is required")
		}
		if req.LastName == "" {
			return nil, errors.New("last name is required")
		}
		if req.Email == "" {
			return nil, errors.New("email is required")
		}

		// Llama a la función `Create` del servicio `Service` para crear el nuevo usuario.
		// Esta función (que probablemente se encuentre en otro paquete) se encarga de la lógica de negocio para persistir el usuario en un repositorio.
		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)

		// Maneja el error en caso de que falle la creación del usuario.
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

/*
Capa de presentación (Controller):

Función: Esta capa se encarga de interactuar directamente con el usuario o cliente.
Maneja las solicitudes HTTP: Recibe las solicitudes del cliente, las analiza y las valida.
Decodifica los datos: Convierte los datos recibidos en el formato adecuado para su procesamiento.
Llama a la capa de servicio: Delega la lógica de negocio a la capa de servicio.
Codifica las respuestas: Convierte los resultados de la capa de servicio en formato JSON o HTML para enviar al cliente.
Envíe las respuestas: Envía las respuestas al cliente junto con el código de estado HTTP correspondiente.
*/
