-- name: CreateGradoOncologico :one
INSERT INTO Grado_oncologico(
      Caracteristica, Muestra_analizada, Puntaje, 
      Descripciones_microscopicas_Descripcion, Descripciones_microscopicas_Pacientes_Protocolo
) VALUES (
      $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetGradoOncologico :one
SELECT * FROM Grado_oncologico
WHERE Descripciones_microscopicas_Pacientes_Protocolo = $1
AND Descripciones_microscopicas_Descripcion = $2;

-- name: UpdateGradoOncologico :one
UPDATE Grado_oncologico
SET 
    Caracteristica = $3,
    Muestra_analizada = $4,
    Puntaje = $5
WHERE Descripciones_microscopicas_Pacientes_Protocolo = $1
AND Descripciones_microscopicas_Descripcion = $2
RETURNING *;
