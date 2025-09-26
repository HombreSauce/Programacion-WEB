-- name: ObtenerUsuario :one
SELECT id, DNI, nombre, apellido, email, telefono, fecha_nacimiento, sexo, rol
FROM usuarios
WHERE id = $1;

-- name: ListarUsuarios :many
SELECT id, DNI, nombre, apellido, email, telefono, fecha_nacimiento, sexo, rol
FROM usuarios
ORDER BY nombre, apellido;

-- name: CrearUsuario :one
INSERT INTO usuarios (DNI, nombre, apellido, email, telefono, fecha_nacimiento, sexo, rol)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, DNI, nombre, apellido, email, telefono, fecha_nacimiento, sexo, rol;

-- name: ActualizarUsuario :exec
UPDATE usuarios
SET DNI = $2, nombre = $3, apellido = $4, email = $5, telefono = $6, fecha_nacimiento = $7, sexo = $8, rol = $9
WHERE id = $1;

-- name: EliminarUsuario :exec
DELETE FROM usuarios
WHERE id = $1;