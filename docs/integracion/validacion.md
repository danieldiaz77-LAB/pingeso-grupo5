> **Actualización 19/11/2026:** este documento conserva información histórica de preparación. El worktree y el merge ya existen; contrato/lock fueron regenerados y hay validaciones posteriores. Para el estado vigente y próximos pasos, leer [continuidad-claude.md](continuidad-claude.md). No volver a ejecutar el script de preparación ni interpretar las tablas antiguas como estado actual.

# Registro y comandos de validación

Fecha: 17 de noviembre de 2026 (Chile).
Estado probado: copia de archivos; no hay SHA de integración ni rama creada.
Bases: pipeline c445c03450692a0c7dd302b07edbbdf72ab20395 y
plantilla 7d2b7dcc216eec847b78b7d7aeb4731aec1cd94d.
El manifiesto SHA-256 que acompaña la copia identifica exactamente los archivos preparados.

| Operación                                                | Resultado de esta sesión                                                                           |
| -------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| git status --short --branch                              | Ejecutado: feat/pipeline-hito0…origin/feat/pipeline-hito0, sin cambios reportados                  |
| Llamada de terminal que incluía git fetch origin         | Bloqueada antes de crear el proceso: CreateProcessAsUserW, error 1260, directiva de grupo          |
| Git fetch, worktree, merge, commits                      | No ejecutados                                                                                      |
| Lectura de DAS e instrucciones                           | Leídos archivos del proyecto y búsqueda en directorios superiores; no se encontró AGENTS.md        |
| Lectura de GitHub                                        | Tres ramas y SHA verificados; main protected=false, rulesets=[], protección detallada 403          |
| Preparación de código                                    | Archivos escritos en carpeta aislada .tmp; no se modificaron los archivos de trabajo originales    |
| Formato/sintaxis JS, TS, JSON, YAML, CSS, Markdown       | Prettier standalone local, mediante API en memoria, sin procesos; detalle en informe-estatico.json |
| pnpm install, pnpm verify, gofmt, go test, sqlc, codegen | No ejecutados por bloqueo de terminal; no se afirma compilación correcta                           |
| E2E proceso real, navegador, mutación, cobertura         | Preparados/conservados, NO ejecutados                                                              |
| Compose completo, Trivy, govulncheck, audit, Sonar       | NO ejecutados sobre esta copia                                                                     |

El administrador debe habilitar la creación de procesos de la terminal PowerShell
usada por el agente (WindowsPowerShell\\v1.0\\powershell.exe) y las herramientas
aprobadas Git, Go, Node/pnpm y Docker. No se probó otra shell ni lanzamiento alternativo
de procesos para sortear la directiva. El error es de Windows, no una solicitud de
aprobación de cambios del proyecto.

## Orden al recuperar la terminal

1. Materializar la copia con el script de preparación adjunto, que comprueba hashes,
   SHA remotos y destino nuevo. Deja la integración local sin commit ni push.
2. Resolver el lockfile y decisiones de dependencias según seguridad.md.
   No aceptar cambios transitivos de versión sin revisar/aprobar.
3. Ejecutar `pnpm install --frozen-lockfile`, `pnpm codegen`,
   `gofmt -w apps/api`, `pnpm format:write`, `pnpm verify`.
   Revisar los generados y registrar un commit local de integración antes de
   `pnpm check:generated`, ya que compara con HEAD.
4. Ejecutar los comandos siguientes en Linux/WSL **habilitado por el administrador**,
   con Go 1.26.5, Node 22.22.0 y pnpm 10.13.1.
   No se propone WSL como evasión de una política de Windows.
5. Guardar comando, salida, exit code, SHA y `git status --porcelain`.
   Si el árbol cambia después del commit, registrar también el diff.

## Base y Compose exclusivamente propios

Desde la raíz del worktree, con terminal habilitada. Cada invocación exige un proyecto nuevo.
Las credenciales son aleatorias y solo para ese entorno efímero. Los puertos son dinámicos.

```bash
set -e
mkdir -p .tmp
id_integracion="$(date +%Y%m%d%H%M%S)-$(openssl rand -hex 4)"
proyecto_integracion="pingeso-integracion-$id_integracion"
nombre_bd="pingeso_integracion_$(openssl rand -hex 8)"
clave_bd="$(openssl rand -hex 24)"
archivo_env=".tmp/$proyecto_integracion.env"
printf 'POSTGRES_USER=prueba\nPOSTGRES_PASSWORD=%s\nPOSTGRES_DB=%s\nDATABASE_URL=postgres://prueba:%s@postgres:5432/%s?sslmode=disable\nAPI_ENV=test\nFRONTEND_URL=http://localhost\nWEB_PORT=0\nBACKEND_PORT=0\n' "$clave_bd" "$nombre_bd" "$clave_bd" "$nombre_bd" > "$archivo_env"
chmod 600 "$archivo_env"
compose_integracion() {
  docker compose --env-file "$archivo_env" -p "$proyecto_integracion" -f docker-compose.yaml -f docker-compose.integracion.yaml "$@"
}
compose_integracion up -d --wait postgres
direccion_bd="$(compose_integracion port postgres 5432)"
export TEST_DATABASE_NAME="$nombre_bd"
export TEST_DATABASE_URL="postgres://prueba:$clave_bd@$direccion_bd/$nombre_bd?sslmode=disable"
pnpm test:mutation
pnpm test:e2e:api
pnpm check:coverage
pnpm test:e2e:web
compose_integracion up -d --build --wait
compose_integracion ps -a
direccion_api="$(compose_integracion port backend 8080)"
direccion_web="$(compose_integracion port frontend 8080)"
curl --fail "http://$direccion_api/healthz"
curl --fail "http://$direccion_api/readyz"
curl --fail "http://$direccion_api/api/v1/ejemplos"
curl --fail "http://$direccion_web/"
curl --fail "http://$direccion_web/api/v1/ejemplos"
compose_integracion logs --no-color migrate
```

Detenerse al primer comando fallido (usar `set -e` en una sesión/script de validación).
Verificar que migrate terminó con código 0 y que backend/frontend/postgres están saludables.
Completar el recorrido crear/listar en navegador contra Compose; la prueba Playwright con
MSW no acredita persistencia real. Mailpit no forma parte de esta validación por defecto.

Limpieza opcional, solo del proyecto nuevo cuyo nombre se conservó en esta misma sesión:
`compose_integracion down --volumes`. No ejecutar prune ni borrar otros proyectos.
La guarda del runner exige nombre de BD coincidente; el aislamiento real lo proporciona
el PostgreSQL creado por esta invocación. No apuntar TEST_DATABASE_URL a una BD existente.

Sonar remoto y el mínimo de cobertura nueva de 80% se mantienen pendientes de medición.
Conservar los checks calidad, seguridad e imagen; el Hito 0 continúa abierto.
