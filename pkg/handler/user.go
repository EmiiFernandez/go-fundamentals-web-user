package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/EmiiFernandez/go-fundamentals-response/response"
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
		// Registra la URL de la solicitud
		url := r.URL.Path
		log.Println(r.Method, ": ", url)

		// Limpia la URL para obtener los parámetros necesarios
		path, pathSize := transport.Clean(url)

		// Crea un mapa para almacenar los parámetros de la URL
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
		case http.MethodPatch:
			switch pathSize {
			case 4:
				end = endpoints.Update
				deco = decodeUpdateUser
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

// decodeCreateUser decodifica los datos de la solicitud para crear un nuevo usuario.
func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateReq
	// Decodifica el cuerpo JSON de la solicitud en la estructura CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("Invalid request format: '%v'", err.Error())
	}
	return req, nil
}

// decodeUpdateUser decodifica los datos de la solicitud para modificar un atributo del usuario.
func decodeUpdateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("Invalid request format: '%v'", err.Error())
	}

	params := ctx.Value("params").(map[string]string)

	id, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, err
	}

	req.ID = id
	return req, nil
}

// encodeResponse codifica la respuesta en formato JSON.
func encodeResponse(tx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	r := resp.(response.Response)
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(resp)
}

// encodeError codifica los errores en formato JSON.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	// Establece el tipo de contenido de la respuesta como JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Convierte el error en una respuesta del tipo response.Response
	resp := err.(response.Response)

	// Establece el código de estado de la respuesta HTTP
	w.WriteHeader(resp.StatusCode())

	// Codifica el error en formato JSON y lo escribe en el cuerpo de la respuesta HTTP
	_ = json.NewEncoder(w).Encode(resp)
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
