-- name: CreateDescripcionMicroscopica :one
INSERT INTO Descripciones_microscopicas (
    Descripcion, Diagnostico, Pacientes_Protocolo
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetDescripcionByPaciente :many
SELECT * FROM Descripciones_microscopicas
WHERE Pacientes_Protocolo = $1;

-- name: UpdateDescripcionMicroscopica :one
UPDATE Descripciones_microscopicas
SET Diagnostico = $3
WHERE Descripcion = $1 AND Pacientes_Protocolo = $2
RETURNING *;

-- name: DeleteDescripcionMicroscopica :exec
DELETE FROM Descripciones_microscopicas
WHERE Descripcion = $1 AND Pacientes_Protocolo = $2;
