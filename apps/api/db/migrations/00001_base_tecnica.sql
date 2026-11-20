-- Línea base técnica: comprueba el motor exigido; no crea entidades de negocio.
-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  IF current_setting('server_version_num')::integer / 10000 <> 17 THEN
    RAISE EXCEPTION 'La base técnica requiere PostgreSQL 17';
  END IF;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- La validación no modificó el esquema; solo se revierte el registro de Goose.
SELECT current_database();
