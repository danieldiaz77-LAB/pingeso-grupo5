# Autenticación JWT

- Salud sin token: healthz 200, readyz 200
- Ejemplos sin token: 401 problem+json
- Token vencido o mal formado: 401
- Token EdDSA válido del JWKS: 200
