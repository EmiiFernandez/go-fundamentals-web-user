package main

import (
	"context"  // Proporciona funcionalidades para manejar contextos en Go
	"fmt"      // Paquete para formateo de salida
	"log"      // Paquete para registro de errores
	"net/http" // Paquete para crear servidores HTTP
	"os"       // Proporciona funciones para interactuar con el sistema operativo

	"github.com/EmiiFernandez/go-fundamentals-web-users/internal/user"
	"github.com/EmiiFernandez/go-fundamentals-web-users/pkg/bootstrap"
	"github.com/EmiiFernandez/go-fundamentals-web-users/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	// Importo las variables de entorno desde el archivo .env
	_ = godotenv.Load()

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
	h := handler.NewUserHTTPServer(user.MakeEndpoints(ctx, service))

	// Importo el puerto desde las variables de entorno
	port := os.Getenv("PORT")
	// Imprime un mensaje indicando el puerto donde se inicia el servidor
	fmt.Println("Server started at port ", port)

	address := fmt.Sprintf("127.0.0.1:%s", port)

	// Configura y lanza el servidor HTTP
	srv := &http.Server{
		Handler: accessControl(h),
		Addr:    address,
	}

	log.Fatal(srv.ListenAndServe())
}

// accessControl agrega encabezados de control de acceso a todas las solicitudes HTTP.
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configura los encabezados de control de acceso permitiendo solicitudes desde cualquier origen.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Configura los métodos HTTP permitidos.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		// Configura los encabezados HTTP permitidos.
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		// Maneja las solicitudes de opción (preflight) y responde directamente sin pasarlas al manejador principal.
		if r.Method == "OPTIONS" {
			return
		}

		// Pasa la solicitud al siguiente manejador en la cadena.
		h.ServeHTTP(w, r)
	})
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

/*
1. El archivo main.go es el punto de entrada de la aplicación.
2. Se importan los paquetes necesarios.
3. La función main() inicia la aplicación:
-- Carga las variables de entorno desde el archivo .env.
-- Establece una conexión a la base de datos MySQL utilizando Docker.
-- Crea un logger para registrar mensajes.
-- Crea un repositorio y un servicio para gestionar usuarios.
-- Configura el servidor HTTP para manejar las solicitudes relacionadas con usuarios.
-- Configura y lanza el servidor en el puerto especificado en las variables de entorno.
4. La función accessControl() es un middleware que agrega encabezados de control de acceso a todas las solicitudes HTTP, permitiendo solicitudes desde cualquier origen, configurando los métodos HTTP permitidos y los encabezados permitidos. Además, maneja las solicitudes de opción (preflight) y las responde directamente sin pasarlas al manejador principal.
*/
