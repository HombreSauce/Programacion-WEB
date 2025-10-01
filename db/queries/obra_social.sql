-- name: CrearObraSocial :one
INSERT INTO obra_social (nombre)
VALUES ($1) 
RETURNING nombre;

-- name: EliminarObraSocial :execrows
DELETE FROM obra_social
WHERE nombre = $1;