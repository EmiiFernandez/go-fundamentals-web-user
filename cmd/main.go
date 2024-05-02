package main

import (
	"context"
	"fmt"      // Paquete para formateo de salida
	"log"      // Paquete para registro de errores
	"net/http" // Paquete para crear servidores HTTP
	"os"

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/user"
	"github.com/EmiiFernandez/go-fundamentals-web-users/pkg/bootstrap"
	"github.com/EmiiFernandez/go-fundamentals-web-users/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	//Importo las variables de entorno
	//por defecto toma el archivo .env. Si cambio el nombre del archivo si necesitaria agregarlo al Load
	_ = godotenv.Load()

	// Crea un nuevo multiplexor para manejar las solicitudes HTTP
	// el multiplexor decide a qué función o controlador debe enviar esa solicitud para ser procesada
	server := http.NewServeMux()

	// Conexión a la base de datos MySQL utilizando Docker
	db, err := bootstrap.NewBD()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Cerrar la conexión a la base de datos al finalizar

	// Verifica la conexión con la base de datos MySQL
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Crea un logger para registrar mensajes en la salida estándar
	logger := bootstrap.NewLogger()

	// Crea un repositorio de usuarios utilizando la base de datos y el logger
	repo := user.NewRepo(db, logger)

	// Crea un servicio de usuarios utilizando el logger y el repositorio
	service := user.NewService(logger, repo)

	// Crea un contexto de fondo para las solicitudes HTTP
	ctx := context.Background()

	// Configura el servidor HTTP para manejar las solicitudes relacionadas con usuarios
	handler.NewUserHTTPServer(ctx, server, user.MakeEndpoints(ctx, service))

	//Importo el puerto desde las variables de entorno
	port := os.Getenv("PORT")
	// Imprime un mensaje indicando el puerto donde se inicia el servidor
	fmt.Println("Server started at port ", port)

	// Inicia el servidor HTTP en el puerto 8080 y maneja cualquier error fatal
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}

/*
Funcionamiento de la arquitectura:

1. Solicitud del cliente: El usuario envía una solicitud HTTP a la aplicación, por ejemplo, a través de un navegador web o una API externa.
2. Capa de presentación: El controlador recibe la solicitud, la analiza y valida los datos recibidos.
3. Decodificación y validación: Los datos se decodifican del formato HTTP (JSON, HTML, etc.) a estructuras de datos Go y se valida su integridad.
4. Llamada al servicio: El controlador llama al servicio correspondiente para realizar la operación de negocio solicitada.
5. Lógica de negocio: El servicio aplica las reglas de negocio, procesa los datos y llama al repositorio para acceder a la fuente de datos.
6. Acceso a datos: El repositorio realiza las operaciones CRUD necesarias sobre la base de datos o el sistema de almacenamiento.
7. Preparación de la respuesta: El servicio recibe los resultados del repositorio y los procesa para construir la respuesta.
8. Codificación y envío: El controlador codifica la respuesta en formato JSON o HTML y la envía al cliente junto con el código de estado HTTP correspondiente.
9. Recepción de la respuesta: El cliente recibe la respuesta y la procesa según el tipo de contenido y el código de estado.
*/
