INSERT INTO ejemplos (id, nombre, descripcion, creado_en) VALUES
    ('11111111-1111-4111-8111-111111111111', 'Primer ejemplo', 'Fila de demostración creada por el seed.', '2026-01-01T12:00:00Z'),
    ('22222222-2222-4222-8222-222222222222', 'Segundo ejemplo', 'Otra fila para verificar el listado.', '2026-01-02T12:00:00Z'),
    ('33333333-3333-4333-8333-333333333333', 'Tercer ejemplo', 'Tercera fila del prototipo.', '2026-01-03T12:00:00Z')
ON CONFLICT (id) DO NOTHING;
