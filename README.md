## Proyecto GO Web API

### Descripción
Este proyecto implementa un servicio web a nivel local utilizando microservicios y aplicaciones con comunicación HTTP 1.1 y el formato de intercambio JSON. La API proporciona operaciones básicas para administrar usuarios, como crear, leer, actualizar y eliminar (CRUD).

### Tecnologías Utilizadas
- Go (Golang)
- HTTP 1.1
- JSON

### Estructura del Proyecto
El proyecto sigue una arquitectura por capas para una mejor organización y mantenimiento del código:

1. **Capa Repositorio:** Contiene las funciones para interactuar con la base de datos o almacenamiento persistente.
2. **Capa Servicio:** Implementa la lógica de negocio y actúa como intermediario entre la capa de repositorio y la capa de controlador.
3. **Capa Controlador:** Define los manejadores HTTP y gestiona las solicitudes entrantes.
