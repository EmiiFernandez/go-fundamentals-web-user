package main

import (
	"context"
	"fmt"      // Paquete para formateo de salida
	"log"      // Paquete para registro de errores
	"net/http" // Paquete para crear servidores HTTP
	"os"

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/domain"
	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/user"
)

func main() {
	// Crea un nuevo multiplexor para manejar las solicitudes HTTP
	server := http.NewServeMux()

	// Simulación de una base de datos en memoria
	db := user.DB{
		Users: []domain.User{
			{ID: 1, FirstName: "Nahuel", LastName: "Costamagna", Email: "nahuel@domain.com"},
			{ID: 2, FirstName: "Eren", LastName: "Jaeger", Email: "eren@domain.com"},
			{ID: 3, FirstName: "Paco", LastName: "Costa", Email: "paco@domain.com"},
		},
		MaxUserID: 3,
	}

	// Crea un logger para registrar mensajes en la salida estándar
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// Crea un repositorio de usuarios utilizando la base de datos y el logger
	repo := user.NewRepo(db, logger)

	// Crea un servicio de usuarios utilizando el logger y el repositorio
	service := user.NewService(logger, repo)

	// Crea un contexto de fondo para las solicitudes HTTP
	ctx := context.Background()

	// Registra un manejador para la ruta "/users" que delega el control a la función MakeEndpoints del paquete user
	server.HandleFunc("/users", user.MakeEndpoints(ctx, service))

	// Imprime un mensaje indicando el puerto donde se inicia el servidor
	fmt.Println("Server started at port 8080")

	// Inicia el servidor HTTP en el puerto 8080 y registra cualquier error fatal
	log.Fatal(http.ListenAndServe(":8080", server))
}

/*
Funcionamiento de la arquitectura:

Solicitud del cliente: El usuario envía una solicitud HTTP a la aplicación, por ejemplo, a través de un navegador web o una API externa.
Capa de presentación: El controlador recibe la solicitud, la analiza y valida los datos recibidos.
Decodificación y validación: Los datos se decodifican del formato HTTP (JSON, HTML, etc.) a estructuras de datos Go y se valida su integridad.
Llamada al servicio: El controlador llama al servicio correspondiente para realizar la operación de negocio solicitada.
Lógica de negocio: El servicio aplica las reglas de negocio, procesa los datos y llama al repositorio para acceder a la fuente de datos.
Acceso a datos: El repositorio realiza las operaciones CRUD necesarias sobre la base de datos o el sistema de almacenamiento.
Preparación de la respuesta: El servicio recibe los resultados del repositorio y los procesa para construir la respuesta.
Codificación y envío: El controlador codifica la respuesta en formato JSON o HTML y la envía al cliente junto con el código de estado HTTP correspondiente.
Recepción de la respuesta: El cliente recibe la respuesta y la procesa según el tipo de contenido y el código de estado.
*/
