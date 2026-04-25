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
