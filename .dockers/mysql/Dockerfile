#imagen base que se utilizará como punto de partida para construir la nueva imagen
FROM mysql:5.7

#Puerto
EXPOSE 3306

# Copiar el script de inicialización a la carpeta de scripts de inicio de Docker
COPY ./init.sql /docker-entrypoint-initdb.d/
