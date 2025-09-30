-- name: ObtenerTurnoPorId :one
SELECT id_turno, id_medico, id_paciente, estado, fecha, hora
FROM turnos
WHERE id_turno = $1;

-- name: ObtenerListaTurnosDePaciente :many
SELECT id_turno, id_medico, id_paciente, estado, fecha, hora
FROM turnos
WHERE id_paciente = $1
  AND estado = 'programado'
  AND (fecha > CURRENT_DATE OR (fecha = CURRENT_DATE AND hora >= CURRENT_TIME))
ORDER BY fecha, hora;

-- name: ObtenerHistorialTurnosDePaciente :many
SELECT id_turno, id_medico, id_paciente, estado, fecha, hora
FROM turnos
WHERE id_paciente = $1
ORDER BY fecha DESC, hora DESC;

-- name: ObtenerListaTurnosPorMedico :many
SELECT id_turno, id_medico, id_paciente, estado, fecha, hora
FROM turnos
WHERE id_medico = $1
  AND estado = 'programado'
  AND (fecha > CURRENT_DATE OR (fecha = CURRENT_DATE AND hora >= CURRENT_TIME))
ORDER BY fecha, hora;

-- name: ObtenerHistorialTurnosPorMedico :many
SELECT id_turno, id_medico, id_paciente, estado, fecha, hora
FROM turnos
WHERE id_medico = $1
ORDER BY fecha DESC, hora DESC;

-- name: CrearTurno :one
INSERT INTO turnos (id_medico, id_paciente, fecha, hora)
VALUES ($1, $2, $3, $4)
RETURNING id_turno, id_medico, id_paciente, estado, fecha, hora;

-- name: ActualizarTurnoDatos :exec
UPDATE turnos
SET id_medico = $2,
    id_paciente = $3,
    fecha = $4,
    hora = $5
WHERE id_turno = $1;

-- name: CancelarTurno :exec
UPDATE turnos
SET estado = 'cancelado'
WHERE id_turno = $1
  AND estado <> 'cancelado';

-- name: CambiarEstadoTurno :exec
UPDATE turnos
SET estado = $2
WHERE id_turno = $1;

-- name: AtenderTurno :exec
UPDATE turnos
SET estado = 'atendido'
WHERE id_turno = $1
  AND estado = 'programado';
