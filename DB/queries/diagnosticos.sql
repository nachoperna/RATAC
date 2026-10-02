-- name: CreateDiagnostico :one
INSERT INTO Diagnosticos (
    diagnostico, grado, protocolo
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetDiagnosticosByProtocolo :many
SELECT * FROM Diagnosticos
WHERE protocolo = $1;

-- name: GetDiagnosticosByDiagnostico :many
SELECT * FROM Diagnosticos
WHERE diagnostico = $1;

-- name: GetDiagnosticosByDiagnosticoAndGrado :many
SELECT * FROM Diagnosticos
WHERE diagnostico = $1 AND grado = $2;

-- name: DeleteDiagnosticoByProtocolo :exec
DELETE FROM Diagnosticos
WHERE protocolo = $1;

-- name: DeleteDiagnosticoById :exec
DELETE FROM Diagnosticos
WHERE id = $1;