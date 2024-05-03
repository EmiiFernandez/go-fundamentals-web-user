## Proyecto GO Web API - Gestión de Usuarios

### Descripción
Este proyecto implementa un servicio web a nivel local utilizando microservicios y aplicaciones con comunicación HTTP 1.1 y el formato de intercambio JSON. La API proporciona operaciones básicas para administrar usuarios, como crear, leer, actualizar y eliminar (CRUD). Utiliza Gin Gonic como framework web y se conecta a una base de datos SQL a través de Docker para persistencia de datos.

## Estructura del Proyecto

El proyecto sigue una estructura modularizada que incluye los siguientes componentes principales:

- **Dominio**: Define la estructura de datos para representar un usuario.
- **Repositorio**: Define las operaciones básicas que se pueden realizar sobre los usuarios en la base de datos.
- **Servicio**: Define las operaciones de alto nivel que se pueden realizar con los usuarios, incluyendo la lógica de negocio.
- **Controladores**: Definen las funciones que manejan las solicitudes HTTP y las traducen a llamadas al servicio.
- **Manejo de Rutas (Handlers)**: Configura las rutas de la aplicación y asigna las funciones de controlador correspondientes.
- **Transporte**: Define funciones para decodificar los datos de la solicitud, llamar a los controladores y codificar las respuestas.
- **Bootstrap**: Configura la inicialización de la aplicación, como la conexión a la base de datos y la configuración del logger.

## Ejecución

Para ejecutar la aplicación, sigue estos pasos:

1. Instala Go en tu sistema si aún no lo tienes: https://golang.org/doc/install
2. Clona este repositorio: `git clone https://github.com/tu_usuario/tu_proyecto.git](https://github.com/EmiiFernandez/go-fundamentals-web-user`
3. Navega al directorio del proyecto: `cd go-fundamentals-web-user`
4. Instala las dependencias del proyecto: `go mod tidy`
5. Ejecuta la aplicación: `go run cmd/main.go`

La aplicación se ejecutará en `http://localhost:8080` por defecto

## Rutas

- **GET** /users: Obtiene todos los usuarios almacenados en la base de datos.
- **GET** /users/:id: Obtiene un usuario específico por su ID.
- **POST** /users: Crea un nuevo usuario con los datos proporcionados.
- **PATCH** /users/:id: Actualiza los datos de un usuario existente.
- **DELETE** /users/:id: Elimina un usuario específico por su ID.
