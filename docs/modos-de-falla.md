# Modos de falla técnicos — Hito 0

Este esqueleto no define el servicio ni reemplaza el contrato funcional de Manuel.
Las pruebas deben fallar cuando no exista PostgreSQL real; no usan una base simulada.

| Falla                                              | Protección                                             | Verificación prevista                                        |
| -------------------------------------------------- | ------------------------------------------------------ | ------------------------------------------------------------ |
| Configuración incompleta o puerto distinto de 8080 | Validación al arrancar                                 | Proceso real termina con error                               |
| PostgreSQL inaccesible                             | `/readyz` hace ping con timeout; `/healthz` sigue vivo | E2E con conexión inválida y concurrente                      |
| Filtración del error SQL o credenciales            | Error público genérico RFC 9457                        | E2E compara cuerpo sin detalles internos                     |
| Contrato publicado en producción                   | Registro de docs condicionado por API_ENV              | E2E exige 404 en prod                                        |
| Contrato o tipos desactualizados                   | Regeneración y estado Git, incluidos archivos nuevos   | Script CI falla ante diferencia                              |
| Primer PR sin contrato anterior                    | Inspección del árbol de la base existente              | Se informa primera línea base; errores Git siguen bloqueando |
| Regresión incompatible                             | oasdiff contra SHA base del PR o SHA anterior del push | Prueba con eliminación de operación                          |
| Sonar sin configurar o sin cobertura               | Preflight explícito y Quality Gate bloqueante          | Job falla; no se omite el análisis                           |
| API expuesta sin migrar                            | Compose espera salida exitosa de migrate               | Compose completo y migración PostgreSQL real aprobados       |
| Proxy web sin backend                              | Timeout y respuesta 502 genérica                       | Prueba del servidor SPA                                      |

Pendiente del equipo: modos de falla de negocio, concurrencia sobre datos, matriz por rol,
idempotencia, validaciones y pruebas de mutación de esas protecciones. El e2e técnico no
acredita esos requisitos funcionales.

## Límites de archivos y arranque

- Una ruta HTTP no debe consultar archivos fuera del bundle ni seguir enlaces simbólicos.
  El servidor carga un catálogo inmutable al iniciar, rechaza enlaces y sirve solo ese catálogo.
  Se prueba un enlace a un archivo externo, traversal codificado y acceso legítimo a assets.
- Configuración inválida, puerto ocupado y healthcheck no saludable deben fallar con salida
  no cero. Los e2e ejecutan los binarios reales y comprueban errores y cierres.
- El generador OpenAPI debe producir el contrato esperado; se ejecuta instrumentado y su
  salida se compara con el contrato versionado, sin excluirlo del alcance de cobertura.
