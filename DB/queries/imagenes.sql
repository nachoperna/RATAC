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
