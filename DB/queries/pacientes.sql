-- name: CreatePaciente :one
INSERT INTO Pacientes (
    Protocolo, Fecha, Solicitante, Tecnica, Familia, 
    Especie, Raza, Edad, Paciente, Antecedentes, 
    Descripcion_macroscopica, Referencias_mastocitomas
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: GetPaciente :one
SELECT * FROM Pacientes
WHERE Protocolo = $1 LIMIT 1;

-- name: ListPacientes :many
SELECT * FROM Pacientes
ORDER BY Fecha DESC;

-- name: ListUltimosPacientes :many
SELECT * FROM Pacientes
ORDER BY Fecha DESC
LIMIT 3;

-- name: UpdatePaciente :one
UPDATE Pacientes
SET 
    Fecha = $2,
    Solicitante = $3,
    Tecnica = $4,
    Familia = $5,
    Especie = $6,
    Raza = $7,
    Edad = $8,
    Paciente = $9,
    Antecedentes = $10,
    Descripcion_macroscopica = $11,
    Referencias_mastocitomas = $12
WHERE Protocolo = $1
RETURNING *;

-- name: DeletePaciente :exec
DELETE FROM Pacientes
WHERE Protocolo = $1;

-- name: CountPacientes :one
SELECT count(*) FROM pacientes;

-- name: GetPacienteByNombre :many
SELECT * FROM Pacientes
WHERE Paciente ilike CONCAT('%', $1::VARCHAR, '%');
