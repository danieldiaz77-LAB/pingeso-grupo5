# Migraciones

La migración 00001 comprueba PostgreSQL 17 sin crear tablas de negocio. Evita que Goose
falle por ausencia de archivos y establece su registro de versión. Benja debe revisar esta
línea base y agregar migraciones nuevas cuando se defina el esquema; no editar las desplegadas.
