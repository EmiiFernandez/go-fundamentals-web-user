package user

import (
	"context"       // Paquete `context`: Proporciona un objeto de contexto para llevar información del ámbito de la solicitud.
	"encoding/json" // Paquete `encoding/json`: Proporciona funciones para codificar y decodificar datos JSON.
	"fmt"           // Paquete `fmt`: Proporciona funciones para el formateo de salida de datos.
	"net/http"      // Paquete `net/http`: Proporciona funciones para crear servidores HTTP.
)

// Definición de tipos

type (
	// Controller: Define un tipo alias para una función que recibe un objeto `http.ResponseWriter` para escribir la respuesta y un objeto `http.Request` para leer la solicitud del cliente.
	Controller func(w http.ResponseWriter, r *http.Request)

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

// Crea y configura los controladores para las diferentes rutas relacionadas con la gestión de usuarios.
func MakeEndpoints(ctx context.Context, s Service) Controller {
	// Esta función crea y devuelve un controlador que maneja las diferentes solicitudes HTTP para la API de usuarios.
	// Recibe el contexto y una instancia del servicio `Service` (que probablemente contiene la lógica de negocio).

	return func(w http.ResponseWriter, r *http.Request) {
		// El controlador comprueba el método HTTP de la solicitud (`GET`, `POST`, etc.) y delega la ejecución a la función correspondiente.
		switch r.Method {
		case http.MethodGet:
			GetAllUser(ctx, s, w) // Delega la solicitud GET a la función `GetAllUser`
		case http.MethodPost:
			// Procesa una solicitud POST para crear un nuevo usuario
			decode := json.NewDecoder(r.Body) // Crea un decodificador JSON para leer el cuerpo de la solicitud
			var req CreateReq                 // Variable para almacenar la estructura de la solicitud
			if err := decode.Decode(&req); err != nil {
				MsgResponse(w, http.StatusBadRequest, err.Error()) // Maneja el error de decodificación
				return
			}
			PostUser(ctx, s, w, req) // Delega la solicitud POST a la función `PostUser` enviando la estructura decodificada
		default:
			InvalidMethod(w) // Maneja métodos HTTP no permitidos
		}
	}
}

func GetAllUser(ctx context.Context, s Service, w http.ResponseWriter) {
	// Esta función recupera todos los usuarios del servicio y los envía como respuesta.
	users, err := s.GetAll(ctx)
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error()) // Maneja el error del servicio
		return
	}
	DataResponse(w, http.StatusOK, users) // Envía la lista de usuarios con código de éxito (200)
}

func PostUser(ctx context.Context, s Service, w http.ResponseWriter, data interface{}) {
	// Esta función maneja las solicitudes POST para crear un nuevo usuario en la API de usuarios.

	// Convierte la interfaz `data` a la estructura `CreateReq` para acceder a los campos del usuario.
	req := data.(CreateReq)

	// Valida los campos obligatorios de la solicitud (nombre, apellido y correo electrónico).
	if req.FirstName == "" {
		MsgResponse(w, http.StatusBadRequest, "first name is required")
		return
	}
	if req.LastName == "" {
		MsgResponse(w, http.StatusBadRequest, "last name is required")
		return
	}
	if req.Email == "" {
		MsgResponse(w, http.StatusBadRequest, "email is required")
		return
	}

	// Llama a la función `Create` del servicio `Service` para crear el nuevo usuario.
	// Esta función (que probablemente se encuentre en otro paquete) se encarga de la lógica de negocio para persistir el usuario en un repositorio.
	user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)

	// Maneja el error en caso de que falle la creación del usuario.
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Si la creación es exitosa, envía una respuesta con el código de estado "Created" (201) y el usuario creado en formato JSON.
	DataResponse(w, http.StatusCreated, user)
}

// Funciones utilitarias del controlador para respuestas HTTP

// InvalidMethod maneja las solicitudes con métodos HTTP no permitidos
func InvalidMethod(w http.ResponseWriter) {
	// Establece el código de estado de la respuesta a "Not Found" (404)
	status := http.StatusNotFound

	// Escribe el código de estado en la respuesta
	w.WriteHeader(status)

	// Construye un mensaje de error indicando que el método no existe
	// en formato JSON {"status": 404, "message": "method doesn't exist"}
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exist"}`, status)
}

// MsgResponse construye y envía respuestas de error en formato JSON
func MsgResponse(w http.ResponseWriter, status int, message string) {
	// Escribe el código de estado en la respuesta
	w.WriteHeader(status)

	// Construye una cadena en formato JSON con el código de estado y el mensaje de error
	// Ejemplo: {"status": 400, "message": "Error al decodificar la solicitud"}
	fmt.Fprintf(w, `{"status": %d, "message": %s}`, status, message)
}

// DataResponse construye y envía respuestas con datos en formato JSON
func DataResponse(w http.ResponseWriter, status int, data interface{}) {
	// Intenta codificar los datos a formato JSON
	value, err := json.Marshal(data)
	if err != nil {
		// Si la codificación falla, envía un error con código "Bad Request" (400)
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Escribe el código de estado en la respuesta
	w.WriteHeader(status)

	// Construye una cadena en formato JSON con el código de estado y los datos codificados
	// Ejemplo: {"status": 201, "data": {"id": 1, "name": "usuario1"}}
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
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
