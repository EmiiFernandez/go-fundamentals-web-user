package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/user"
	"github.com/EmiiFernandez/go-fundamentals-web-users/pkg/transport"
)

// NewUserHTTPServer configura las rutas del servidor HTTP para los endpoints de usuarios.
func NewUserHTTPServer(ctx context.Context, router *http.ServeMux, endpoints user.Endpoints) {
	// Configura las rutas de los endpoints para el servidor HTTP
	router.HandleFunc("/users", UserServer(ctx, endpoints))
	router.HandleFunc("/users/", UserServer(ctx, endpoints))
}

// UserServer maneja las solicitudes HTTP relacionadas con los usuarios.
func UserServer(ctx context.Context, endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		log.Println(r.Method, ": ", url)

		// Limpia la URL para obtener los parámetros necesarios
		path, pathSize := transport.Clean(url)

		params := make(map[string]string)
		// Si hay un parámetro de ID de usuario en la URL, guárdalo
		if pathSize == 4 && path[2] != "" {
			params["userID"] = path[2]
		}

		// Crea un nuevo objeto transport para manejar la solicitud y la respuesta
		tran := transport.New(w, r, context.WithValue(ctx, "params", params))

		var end user.Controller
		var deco func(ctx context.Context, r *http.Request) (interface{}, error)

		// Determina el controlador y el decodificador adecuados según el método y la ruta de la solicitud
		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3:
				end = endpoints.GetAll
				deco = decodeGetAllUser
			case 4:
				end = endpoints.Get
				deco = decodeGetUser
			}
		case http.MethodPost:
			switch pathSize {
			case 3:
				end = endpoints.Create
				deco = decodeCreateUser
			}
		}

		// Si se encontró un controlador y un decodificador válidos, procesa la solicitud
		if end != nil && deco != nil {
			tran.Server(
				transport.Endpoint(end),
				deco,
				encodeResponse,
				encodeError,
			)
		} else {
			// Si no se encontró un método válido, devuelve un error
			InvalidMethod(w)
		}
	}
}

// decodeGetUser decodifica los parámetros de la solicitud para obtener el ID del usuario.
func decodeGetUser(ctx context.Context, r *http.Request) (interface{}, error) {
	params := ctx.Value("params").(map[string]string)

	id, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, err
	}

	// Retorna un objeto GetReq que contiene el ID del usuario
	return user.GetReq{
		ID: id,
	}, nil
}

// decodeGetAllUser decodifica los parámetros de la solicitud para obtener todos los usuarios.
func decodeGetAllUser(ctx context.Context, r *http.Request) (interface{}, error) {
	// No se necesitan parámetros adicionales para obtener todos los usuarios
	return nil, nil
}

// encodeResponse codifica la respuesta en formato JSON.
func encodeResponse(tx context.Context, w http.ResponseWriter, resp interface{}) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	status := http.StatusOK
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Escribe la respuesta JSON en el cuerpo de la respuesta HTTP
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, data)
	return nil
}

// encodeError codifica los errores en formato JSON.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Escribe el mensaje de error JSON en el cuerpo de la respuesta HTTP
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, err.Error())
}

// decodeCreateUser decodifica los datos de la solicitud para crear un nuevo usuario.
func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateReq
	// Decodifica el cuerpo JSON de la solicitud en la estructura CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("Invalid request format: '%v'", err.Error())
	}
	return req, nil
}

// InvalidMethod envía una respuesta de error cuando el método HTTP no es compatible con el endpoint.
func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound

	// Escribe el mensaje de error JSON en el cuerpo de la respuesta HTTP
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exist"}`, status)
}

// MsgResponse envía una respuesta personalizada con el estado y el mensaje proporcionados.
func MsgResponse(w http.ResponseWriter, status int, message string) {
	// Escribe el mensaje personalizado JSON en el cuerpo de la respuesta HTTP
	fmt.Fprintf(w, `{"status": %d, "message": %s}`, status, message)
}
