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

-- name: CreateImagen :one
INSERT INTO Imagenes (
    Ruta, Descripciones_microscopicas_Descripcion, Descripciones_microscopicas_Pacientes_Protocolo
) VALUES (
    $1, $2, $3
)
RETURNING *; 

-- name: GetImagenesByDescripcion :many
SELECT * FROM Imagenes
WHERE Descripciones_microscopicas_Descripcion = $1 
AND Descripciones_microscopicas_Pacientes_Protocolo = $2;

-- name: DeleteImagen :exec
DELETE FROM Imagenes
WHERE Ruta = $1;
-- name: GetGradoOncologico :one
SELECT * FROM Grado_oncologico
WHERE Descripciones_microscopicas_Pacientes_Protocolo = $1
AND Descripciones_microscopicas_Descripcion = $2;

-- name: UpdateGradoOncologico :one
UPDATE Grado_oncologico
SET 
    Caracteristica = $3,
    Muestra_analizada = $4,
    Puntaje = $5,
    Total = $6
WHERE Descripciones_microscopicas_Pacientes_Protocolo = $1
AND Descripciones_microscopicas_Descripcion = $2
RETURNING *;