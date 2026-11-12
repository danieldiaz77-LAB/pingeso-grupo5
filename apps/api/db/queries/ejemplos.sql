-- name: ListarEjemplos :many
SELECT id, nombre, descripcion, creado_en
FROM ejemplos
ORDER BY creado_en ASC, id ASC;

-- name: CrearEjemplo :one
INSERT INTO ejemplos (nombre, descripcion)
VALUES ($1, $2)
RETURNING id, nombre, descripcion, creado_en;
