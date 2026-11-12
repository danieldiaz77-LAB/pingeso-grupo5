-- +goose Up
CREATE TABLE ejemplos (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre text NOT NULL,
    descripcion text NOT NULL DEFAULT '',
    creado_en timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE ejemplos;
