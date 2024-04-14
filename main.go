package main

import (
	"encoding/json" // Paquete para codificar y decodificar datos JSON
	"fmt"           // Paquete para formateo de salida
	"log"           // Paquete para registro de errores
	"net/http"      // Paquete para crear servidores HTTP
)

// Definición de la estructura del usuario
type User struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Array de usuarios
var users []User

// Función para inicializar valores. Se ejecuta primero.
func init() {
	// Inicialización de usuarios
	users = []User{
		{
			ID:        1,
			FirstName: "Nahuel",
			LastName:  "Costamagna",
			Email:     "nahuel@domain.com",
		},
		{
			ID:        2,
			FirstName: "Eren",
			LastName:  "Jaeger",
			Email:     "eren@domain.com",
		},
		{
			ID:        3,
			FirstName: "Paco",
			LastName:  "Costamagna",
			Email:     "paco@domain.com",
		},
	}
}

func main() {
	// Inicializamos un servidor HTTP en el puerto 8080
	// Asignamos la función UserServer para manejar las solicitudes en la ruta "/users".
	http.HandleFunc("/users", UserServer)

	// Escuchamos por solicitudes entrantes en el puerto 8080.
	// La función ListenAndServe siempre devuelve un error (que será no nulo) a menos que haya un problema,
	// en cuyo caso log.Fatal manejará la salida del programa y mostrará el mensaje de error.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// UserServer maneja las solicitudes HTTP dirigidas a la ruta "/users".
func UserServer(w http.ResponseWriter, r *http.Request) {
	var status int
	switch r.Method {
	case http.MethodGet:
		// Si la solicitud es un GET, llamamos a la función GetAllUser
		GetAllUser(w)
	case http.MethodPost:
		// Si la solicitud es un POST, devolvemos un mensaje de éxito
		status = 200
		w.WriteHeader(status)
		fmt.Fprintf(w, `{ "status": %d, "message": %s}`, status, "success in post")
	default:
		// Si el método no es GET o POST, devolvemos un mensaje de ruta no encontrada
		status = 404
		w.WriteHeader(status)
		fmt.Fprintf(w, `{ "status": %d, "message": %s}`, status, "not found")
	}
}

// GetAllUser devuelve todos los usuarios en formato JSON
func GetAllUser(w http.ResponseWriter) {
	// Llamamos a DataResponse para escribir la respuesta HTTP con los usuarios en formato JSON
	DataResponse(w, http.StatusOK, users)
}

// DataResponse escribe la respuesta HTTP con los datos en formato JSON y el código de estado especificado
func DataResponse(w http.ResponseWriter, status int, data interface{}) {
	// Convertimos los datos a formato JSON
	value, _ := json.Marshal(data)
	// Escribimos la respuesta con el código de estado y los datos en formato JSON
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}
