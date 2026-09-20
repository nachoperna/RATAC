-- name: GetSolicitudes :many
SELECT email, nombre_lab, ciudad_origen
FROM Usuarios
WHERE rol = 'solicitante';

-- name: GetVeterinarios :many
SELECT matricula, nombre
FROM Veterinarios
WHERE email_lab LIKE $1;

-- name: CreateSolicitud :one
INSERT INTO Usuarios 
      (id, email, nombre_lab, contraseña_hash, rol, ciudad_origen, activo, fecha_creacion, fecha_modificacion)
      VALUES (default, $1, $2, $3, 'solicitante'::roles, $4, true, default, default)
RETURNING *;

-- name: SetVetinarios :exec
INSERT INTO Veterinarios 
      (matricula, nombre, email_lab)
      VALUES ($1, $2, $3);
