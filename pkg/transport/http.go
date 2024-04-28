package transport

import (
	"context"
	"net/http"
	"strings"
)

// Transport define el contrato para un objeto que maneja la ejecución de un endpoint.
type Transport interface {
	// Server es responsable de:
	// - Decodificar la solicitud HTTP (usando decode).
	// - Ejecutar el endpoint (pasando el contexto y los datos decodificados).
	// - Codificar la respuesta del endpoint (usando encode).
	// - Manejar errores de decodificación, ejecución del endpoint o codificación de la respuesta (usando encodeError).
	Server(
		endpoint Endpoint,
		decode func(ctx context.Context, r *http.Request) (interface{}, error),
		encode func(ctx context.Context, w http.ResponseWriter, resp interface{}) error,
		encodeError func(ctx context.Context, err error, w http.ResponseWriter),
	)
}

// Endpoint encapsula la lógica específica de procesamiento de una solicitud HTTP.
// Recibe el contexto y los datos decodificados, y devuelve la respuesta o un error.
type Endpoint func(ctx context.Context, request interface{}) (interface{}, error)

// transport implementa la interfaz Transport.
// Almacena la respuesta HTTP (w), la solicitud (r) y el contexto (ctx) para su uso dentro del método Server.
type transport struct {
	w   http.ResponseWriter
	r   *http.Request
	ctx context.Context
}

// New crea una nueva instancia de transport con la respuesta HTTP, la solicitud y el contexto proporcionados.
func New(w http.ResponseWriter, r *http.Request, ctx context.Context) Transport {
	return &transport{
		w:   w,
		r:   r,
		ctx: ctx,
	}
}

// Server comunica con HTTP en su aplicación.
func (t *transport) Server(
	endpoint Endpoint,
	decode func(ctx context.Context, r *http.Request) (interface{}, error),
	encode func(ctx context.Context, w http.ResponseWriter, resp interface{}) error,
	encodeError func(ctx context.Context, err error, w http.ResponseWriter),
) {
	// Decodificación de la solicitud.
	// Decodifica el cuerpo de la solicitud y cualquier parámetro de ruta o consulta en un formato utilizable por el endpoint.
	data, err := decode(t.ctx, t.r)
	if err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

	// Ejecución del endpoint.
	// Debe procesar la solicitud y devolver la respuesta calculada en la variable res.
	res, err := endpoint(t.ctx, data)
	if err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

	// Codificación de la respuesta.
	// Es responsable de convertir la respuesta del endpoint a un formato de respuesta HTTP adecuado.
	if err := encode(t.ctx, t.w, res); err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}
}

// Clean elimina las barras diagonales adicionales al principio y al final de la URL y devuelve las partes de la URL divididas por las barras diagonales.
func Clean(url string) ([]string, int) {
	if url[0] != '/' {
		url = "/" + url
	}

	if url[len(url)-1] != '/' {
		url = url + "/"
	}

	parts := strings.Split(url, "/")

	return parts, len(parts)
}
