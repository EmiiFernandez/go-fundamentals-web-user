package transport

//Separar la lógica de negocio de la capa de transporte (servidor HTTP)

import (
	"context"
	"net/http"
)

// Define el contrato para un objeto que maneja la ejecución de un endpoint.
type Transport interface {
	/*
		El método Server es responsable de:
			.Decodificar la solicitud HTTP (usando decode).
			.Ejecutar el endpoint (pasando el contexto y los datos decodificados).
			.Codificar la respuesta del endpoint (usando encode).
			.Manejar errores de decodificación, ejecución del endpoint o codificación de la respuesta (usando encodeError).
	*/
	Server(
		//función que encapsula la lógica del negocio del endpoin
		endpoint Endpoint,
		//Función para decodificar la solicitud HTTP y obtener los datos que se pasarán al endpoint
		decode func(ctx context.Context, r *http.Request) (interface{}, error),
		//Función para codificar la respuesta del endpoint y enviarla al cliente.
		encode func(ctx context.Context, w http.ResponseWriter, resp interface{}) error,
		//Función para codificar y enviar un error al cliente.
		encodeError func(ctx context.Context, err error, w http.ResponseWriter),
	)
}

// Encapsula la lógica específica de procesamiento de una solicitud HTTP. Recibe el contexto y los datos decodificados, y devuelve la respuesta o un error
type Endpoint func(ctx context.Context, request interface{}) (interface{}, error)

// Implementa la interfaz Transport. Almacena la respuesta HTTP (w), la solicitud (r) y el contexto (ctx) para su uso dentro del método Server.
type transport struct {
	w   http.ResponseWriter
	r   *http.Request
	ctx context.Context
}

// Crea una nueva instancia de transport con la respuesta HTTP, la solicitud y el contexto proporcionados
func New(w http.ResponseWriter, r *http.Request, ctx context.Context) Transport {
	return &transport{
		w:   w,
		r:   r,
		ctx: ctx,
	}
}

// comunicación HTTP en su aplicación
func (t *transport) Server(
	endpoint Endpoint,
	decode func(ctx context.Context, r *http.Request) (interface{}, error),
	encode func(ctx context.Context, w http.ResponseWriter, resp interface{}) error,
	encodeError func(ctx context.Context, err error, w http.ResponseWriter),
) {
	//Decodificación de la solicitud
	//decodificar el cuerpo de la solicitud y cualquier parámetro de ruta o consulta en un formato utilizable por el endpoint.
	data, err := decode(t.ctx, t.r)
	if err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

	//Ejecución del endpoint
	//debe procesar la solicitud y devolver la respuesta calculada en la variable res
	res, err := endpoint(t.ctx, data)
	if err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

	// Codificación de la respuesta
	//es responsable de convertir la respuesta del endpoint a un formato de respuesta HTTP adecuado
	if err := encode(t.ctx, t.w, res); err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

}
