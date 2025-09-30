-- name: ListarObrasPorMedico :many
SELECT ap.obra_social_nombre
FROM atiende_por ap
WHERE ap.id_medico = $1
ORDER BY ap.obra_social_nombre;

-- name: ListarMedicosPorObra :many
SELECT m.id_medico, u.DNI, u.nombre, u.apellido, m.nro_matricula, m.especialidad
FROM atiende_por ap
JOIN medicos m ON m.id_medico = ap.id_medico
JOIN usuarios u ON u.id = m.id_medico
WHERE ap.obra_social_nombre = $1
ORDER BY u.apellido, u.nombre;

-- name: CrearRelacionMedicoObra :one
INSERT INTO atiende_por (id_medico, obra_social_nombre)
VALUES ($1, $2)
RETURNING id_medico, obra_social_nombre;

-- name: ActualizarObraSocialDeMedico :exec
UPDATE atiende_por
SET obra_social_nombre = $2
WHERE id_medico = $1
  AND obra_social_nombre = $3;

-- name: EliminarRelacionMedicoObra :exec
DELETE FROM atiende_por
WHERE id_medico = $1
  AND obra_social_nombre = $2;
