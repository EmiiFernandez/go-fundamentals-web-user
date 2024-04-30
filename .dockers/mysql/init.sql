/*Iniciar BDD*/

SET @MYSQLDUMP_TEMP_LOG_BIN = @SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN = 0;
SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ '';

/*Creo BDD*/
CREATE DATABASE IF NOT EXISTS `go_course_users`;

/*Tabla usuarios*/
CREATE TABLE `go_course_users`.`users` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `first_name` VARCHAR(45) NULL,
    `last_name` VARCHAR(45) NULL,
    `email` VARCHAR(45) NULL,
    PRIMARY KEY(`id`)
);

/*Script de inicialización de una base de datos MySQL. 

SET @MYSQLDUMP_TEMP_LOG_BIN = @SESSION.SQL_LOG_BIN;: Esta sentencia guarda el valor actual de la configuración de los registros binarios de MySQL en una variable temporal llamada MYSQLDUMP_TEMP_LOG_BIN. En esencia, deshabilita temporalmente la escritura de registros binarios para el volcado.
SET @@SESSION.SQL_LOG_BIN = 0;: Se desactiva la escritura de registros binarios para la sesión actual, es decir que las operaciones de modificación de datos realizadas durante la ejecución de este script no se registrarán en el log binario.
SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ '';: Establece la variable global GTID_PURGED en una cadena vacía. Los GTID (Global Transaction ID) son identificadores únicos asignados a cada transacción en un servidor MySQL. Al establecer GTID_PURGED en una cadena vacía, se limpia cualquier información de GTID previamente purgada, lo que puede ser útil para evitar conflictos o asegurar que las transacciones futuras se registren correctamente.
*/