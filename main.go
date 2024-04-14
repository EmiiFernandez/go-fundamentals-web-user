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

// autoincremental
var maxID uint64

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

	maxID = getNextID()
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
		// Si la solicitud es un POST, decodificamos el cuerpo JSON de la solicitud para obtener los datos del nuevo usuario
		decode := json.NewDecoder(r.Body)
		var u User
		if err := decode.Decode(&u); err != nil {
			// Si hay un error al decodificar los datos, respondemos con un mensaje de error y el código de estado 400 Bad Request
			MsgResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// Llamamos a la función PostUser para agregar el nuevo usuario y responder con los datos del usuario creado
		PostUser(w, u)
	default:
		// Si el método no es GET o POST, devolvemos un mensaje de ruta no encontrada y el código de estado 404 Not Found
		status = 404
		w.WriteHeader(status)
		fmt.Fprintf(w, `{ "status": %d, "message": %s}`, status, "not found")
	}
}

// GetAllUser devuelve todos los usuarios en formato JSON
func GetAllUser(w http.ResponseWriter) {
	// Llamamos a DataResponse para escribir la respuesta HTTP con los usuarios en formato JSON y el código de estado 200 OK
	DataResponse(w, http.StatusOK, users)
}

// PostUser agrega un nuevo usuario a la lista de usuarios y responde con los datos del usuario creado
func PostUser(w http.ResponseWriter, data interface{}) {
	// Convertir los datos a un objeto User
	user := data.(User)
	// Asignar el siguiente ID disponible al nuevo usuario
	user.ID = getNextID()
	// Agregar el nuevo usuario a la lista de usuarios
	users = append(users, user)
	// Responder con los datos del usuario creado y el código de estado 201 Created
	DataResponse(w, http.StatusCreated, user)
}

// MsgResponse responde con un mensaje y un código de estado HTTP específico
func MsgResponse(w http.ResponseWriter, status int, message string) {
	// Escribimos el mensaje y el código de estado en formato JSON
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, message)
}

// DataResponse escribe la respuesta HTTP con los datos en formato JSON y el código de estado especificado
func DataResponse(w http.ResponseWriter, status int, data interface{}) {
	// Convertimos los datos a formato JSON
	value, err := json.Marshal(data)
	if err != nil {
		// Si hay un error al convertir los datos a JSON, respondemos con un mensaje de error y el código de estado 400 Bad Request
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// Escribimos la respuesta con el código de estado y los datos en formato JSON
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}

// getNextID calcula el siguiente ID disponible basado en la lista de usuarios actual
func getNextID() uint64 {
	var nextID uint64
	// Iterar sobre la lista de usuarios para encontrar el máximo ID
	for _, user := range users {
		if user.ID > nextID {
			nextID = user.ID
		}
	}
	// Incrementar el máximo ID encontrado
	return nextID + 1
}
