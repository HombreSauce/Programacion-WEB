-- name: ObtenerPaciente :one
SELECT u.id, u.DNI, u.nombre, u.apellido, u.sexo, u.fecha_nacimiento, u.email, u.telefono, p.obra_social, p.nro_afiliado
FROM usuarios u JOIN pacientes p ON p.id_paciente = u.id
WHERE id_paciente = $1;

-- name: ListarPacientes :many
SELECT u.id, u.DNI, u.nombre, u.apellido, u.sexo, u.fecha_nacimiento, u.email, u.telefono, p.obra_social, p.nro_afiliado
FROM usuarios u JOIN pacientes p ON p.id_paciente = u.id
ORDER BY u.apellido, u.nombre;

-- name: CrearPaciente :one
INSERT INTO pacientes (id_paciente, obra_social, nro_afiliado)
VALUES ($1, $2, $3)
RETURNING obra_social, nro_afiliado;

-- name: ActualizarPaciente :exec
UPDATE pacientes
SET obra_social = $2, nro_afiliado = $3
WHERE id_paciente = $1;

-- name: EliminarPaciente :exec
DELETE FROM pacientes
WHERE id_paciente = $1;