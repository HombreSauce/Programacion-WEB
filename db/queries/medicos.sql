-- name: ObtenerMedico :one
SELECT u.ID, u.DNI, u.nombre, u.apellido, m.nro_matricula, m.especialidad
FROM usuarios u JOIN medicos m ON m.id_medico = u.id
WHERE id_medico = $1;

-- name: ListarMedicos :many
SELECT u.ID, u.DNI, u.nombre, u.apellido, m.nro_matricula, m.especialidad
FROM usuarios u JOIN medicos m ON m.id_medico = u.id
ORDER BY u.apellido, u.nombre;

-- name: CrearMedico :one
INSERT INTO medicos (id_medico, nro_matricula, especialidad)
VALUES ($1, $2, $3)
RETURNING nro_matricula, especialidad;

-- name: ActualizarMedico :exec
UPDATE medicos
SET nro_matricula = $2, especialidad = $3
WHERE id_medico = $1;

-- name: EliminarMedico :execrows
DELETE FROM medicos
WHERE id_medico = $1;