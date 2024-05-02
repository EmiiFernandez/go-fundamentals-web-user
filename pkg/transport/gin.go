/*Gin Gonic es un framework web para Go que permite construir aplicaciones web y API RESTful. Proporciona enrutamiento, middleware, gestión de solicitudes y respuestas HTTP, entre otras, lo que facilita la creación de servicios web robustos y eficientes en Go. El paquete contiene las funciones y estructuras necesarias para trabajar con Gin en una aplicación Go.
 */
package transport

import (
	"github.com/gin-gonic/gin"
)

// GinServer crea un manejador HTTP utilizando Gin Gonic.
// Toma un endpoint, funciones para decodificar, codificar y manejar errores y devuelve un manejador HTTP compatible con Gin.
// *gin.Context es una estructura que encapsula todas las variables y funcionalidades relacionadas con una solicitud HTTP y su respuesta asociada en el contexto de una aplicación web
func GinServer(
	endpoint Endpoint,
	decode func(c *gin.Context) (interface{}, error),
	encode func(c *gin.Context, resp interface{}),
	encodeError func(c *gin.Context, err error),
) func(c *gin.Context) {

	// La función anónima devuelta actúa como el manejador HTTP para Gin.
	return func(c *gin.Context) {
		// Decodifica la solicitud utilizando la función de decodificación proporcionada.
		data, err := decode(c)
		if err != nil {
			// Si hay un error durante la decodificación, se codifica el error y se envía como respuesta.
			encodeError(c, err)
			return
		}

		// Llama al endpoint proporcionado con el contexto de la solicitud y los datos decodificados.
		res, err := endpoint(c.Request.Context(), data)
		if err != nil {
			// Si hay un error al llamar al endpoint, se codifica el error y se envía como respuesta.
			encodeError(c, err)
			return
		}

		// Codifica la respuesta y la envía como respuesta HTTP.
		encode(c, res)
	}
}
