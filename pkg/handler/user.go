package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/EmiiFernandez/go-fundamentals-response/response"
	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/user"
	"github.com/EmiiFernandez/go-fundamentals-web-users/pkg/transport"
	"github.com/gin-gonic/gin"
)

// NewUserHTTPServer configura un servidor HTTP utilizando Gin para los endpoints relacionados con usuarios.
func NewUserHTTPServer(endpoints user.Endpoints) http.Handler {
	// Se crea un nuevo enrutador Gin con la configuración predeterminada.
	r := gin.Default()

	// Configuración de los endpoints para crear, obtener todos, obtener uno y actualizar usuarios.
	r.POST("/users", transport.GinServer(
		transport.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
		encodeError,
	))
	r.GET("/users", transport.GinServer(
		transport.Endpoint(endpoints.GetAll),
		decodeGetAllUser,
		encodeResponse,
		encodeError,
	))
	r.GET("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		encodeError,
	))
	r.PATCH("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		encodeError,
	))
	r.DELETE("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.Delete),
		decodeDeleteUser,
		encodeResponse,
		encodeError,
	))

	return r // Retorna el enrutador Gin como un manejador HTTP.
}

// decodeGetUser decodifica los parámetros de la solicitud para obtener el ID del usuario.
func decodeGetUser(c *gin.Context) (interface{}, error) {
	// Verifica si el token de autorización es válido.
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	// Obtiene el ID del usuario de los parámetros de la URL.
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	// Retorna un objeto GetReq que contiene el ID del usuario.
	return user.GetReq{
		ID: id,
	}, nil
}

// decodeGetAllUser decodifica los parámetros de la solicitud para obtener todos los usuarios.
func decodeGetAllUser(c *gin.Context) (interface{}, error) {
	// Verifica si el token de autorización es válido.
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	// No se necesitan parámetros adicionales para obtener todos los usuarios.
	return nil, nil
}

// decodeCreateUser decodifica los datos de la solicitud para crear un nuevo usuario.
func decodeCreateUser(c *gin.Context) (interface{}, error) {
	// Verifica si el token de autorización es válido.
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	// Decodifica el cuerpo JSON de la solicitud en la estructura CreateReq.
	var req user.CreateReq
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

// decodeUpdateUser decodifica los datos de la solicitud para modificar un atributo del usuario.
func decodeUpdateUser(c *gin.Context) (interface{}, error) {
	// Se declara una variable para contener los datos de la solicitud de actualización del usuario.
	var req user.UpdateReq

	// Decodifica los datos JSON de la solicitud en la estructura user.UpdateReq.
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid request format: '%v'", err.Error()))
	}

	// Verifica si el token de autorización es válido.
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	// Convierte el ID de usuario de tipo cadena a tipo uint64 para usarlo en la solicitud de actualización.
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error()) // Se devuelve un error si no se puede convertir el ID a uint64.
	}

	// Asigna el ID de usuario convertido a la solicitud de actualización antes de devolverla.
	req.ID = id
	return req, nil // Se devuelve la solicitud de actualización decodificada y sin errores.
}

// decodeDeleteUser decodifica los parámetros de la solicitud para obtener el ID del usuario y eliminar el usuario.
func decodeDeleteUser(c *gin.Context) (interface{}, error) {
	// Verifica si el token de autorización es válido.
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	// Obtiene el ID del usuario de los parámetros de la URL.
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	// Crea una instancia de DeleteReq con el ID del usuario.
	req := user.DeleteReq{
		ID: id,
	}

	// Devuelve la estructura DeleteReq.
	return req, nil
}

/*
// decodeGetUser decodifica los parámetros de la solicitud para obtener
	// Obtiene el ID del usuario de los parámetros de la URL.
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	// Retorna un objeto GetReq que contiene el ID del usuario.

}
*/

// tokenVerify verifica si el token proporcionado coincide con el token almacenado en las variables de entorno.
// Si el token no coincide, devuelve un error.
func tokenVerify(token string) error {
	// Compara el token proporcionado con el token almacenado en las variables de entorno.
	// Si no coinciden, devuelve un error de "token inválido".
	if os.Getenv("TOKEN") != token {
		return errors.New("invalid token")
	}

	// Si los tokens coinciden, devuelve nil (sin error).
	return nil
}

// encodeResponse codifica la respuesta en formato JSON.
func encodeResponse(c *gin.Context, resp interface{}) {
	// Obtiene la respuesta como una estructura de respuesta genérica.
	r := resp.(response.Response)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(r.StatusCode(), resp) // Codifica la respuesta como JSON y la envía al cliente.
}

// encodeError codifica los errores en formato JSON.
func encodeError(c *gin.Context, err error) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	c.JSON(resp.StatusCode(), resp) // Codifica el error como JSON y lo envía al cliente.
}
