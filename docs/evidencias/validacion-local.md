# Validación local — 8 de noviembre de 2026

Rama `feat/pipeline-hito0`, HEAD `ed6f74d33d696eb4cafe807291f3f891451928b1`.
Registro de validación local anterior al commit y al PR; sin modificación de rulesets. Se preservó el DAS original sin editarlo.
Las herramientas portátiles están ignoradas en `.tools/`.

| Comprobación                    | Resultado observado                                                                  |
| ------------------------------- | ------------------------------------------------------------------------------------ |
| Versiones y SHA oficiales       | Verificados; dos SHA del ejemplo DAS corregidos conservando versiones                |
| Digests                         | Índices oficiales verificados; conservado distroless histórico del DAS               |
| pnpm install --frozen-lockfile  | Aprobado; advertencia peer TS6/openapi-typescript documentada                        |
| Formato, lint y tipos           | Aprobados en la verificación integrada final (pnpm verify, código 0)                 |
| Go mod verify, vet, test, build | Aprobados; las pruebas de salud aisladas no sustituyen PostgreSQL real               |
| Vitest                          | 3 pruebas aprobadas; LCOV real, 100% líneas instrumentadas (4/4) de esta base mínima |
| Playwright Chromium             | 1 e2e de navegador aprobado                                                          |
| Contrato y tipos                | Generados desde Huma y openapi-typescript; TypeScript aprobado                       |
| Primera base sin contrato       | Informa primera línea base y termina correctamente                                   |
| SHA base inexistente            | Rechazado con error Git 128; no se confunde con ausencia de contrato                 |
| Generados nuevos                | Rechazados como se espera: aún no están commiteados                                  |
| oasdiff negativo                | Eliminación temporal de /readyz rechazada, código 1; contrato real intacto           |
| Mutación readiness              | Protección eliminada solo en copia temporal; prueba detecta 200 en vez de 503        |
| actionlint 1.7.7                | Workflow válido                                                                      |
| Hadolint 2.15.1                 | Aprobado al umbral error; dos INFO DL3066 por USER nonroot del DAS                   |
| Gitleaks 8.30.1                 | Sin secretos en historial ni archivos versionables escaneados                        |
| Gitleaks staged negativo        | Fixture sintético con formato válido bloqueado en repo temporal                      |
| Reglas frontend negativas       | localStorage y fetch fuera de *.api.ts rechazados                                    |
| Hooks                           | Instalados pre-commit y pre-push; no se creó un commit real                          |
| pnpm audit                      | FALLA: 1 HIGH, 6 moderate, 1 low; Vite 8.1.5 requiere aprobación para cambiar        |
| govulncheck                     | FALLA: 16 avisos stdlib + 2 x/text alcanzables; salida adjunta                       |
| Trivy filesystem                | FALLA: 2 HIGH con fix en x/text; sin HIGH/CRITICAL de configuración en Dockerfiles   |
| Sonar público                   | Proyecto/organización y Sonar way confirmados; condiciones adjuntas                  |
| Sonar sin token local           | Preflight falla expresamente; no se simula un análisis aprobado                      |

## Validación Docker posterior

Docker Desktop **4.62.0** y Docker Engine **29.2.1**, Linux/amd64 sobre WSL2, operativos.
El problema inicial de pipe/timeout ya no bloquea. Se usó una red exclusiva y PostgreSQL
17.11 con datos efímeros en tmpfs, sin puertos publicados. La limpieza comprueba etiquetas
y solo elimina contenedores/red de esta ejecución: no se borraron volúmenes ni recursos
ajenos. Las imágenes locales y cachés de build quedan disponibles para revisión.

Primer e2e con BD real: migración aprobada; falló `TestSaludReal` porque Huma añadía
`$schema` con el host de la petición pese a tener las rutas de esquema deshabilitadas.
No fue una filtración de la cadena PostgreSQL. Se desactivó el hook predeterminado de enlaces
solo en producción y se añadió la comprobación explícita de ausencia de `$schema`; se
conservaron las comprobaciones de errores sanitizados. Evidencia histórica:
`docker-primer-intento.json`. La segunda ejecución usa los archivos corregidos; sus hashes
SHA256 dentro del contenedor coincidieron con el árbol local.

**Segunda ejecución aprobada (código 0):** migración Goose a versión 1 y cuatro pruebas
`TestReadiness`, `TestSaludReal`, `TestServidorWebReal`, `TestSPA`. Resumen sin logs de
peticiones en `e2e-docker.txt`; códigos y tiempos en `docker-validacion.json`.

Se produjo `coverage.out` auténtico desde contadores de los procesos instrumentados:
**104/141 sentencias (73,76%)**, cuatro archivos. `pnpm check:coverage` terminó 0 y comprobó
formato Go, existencia de fuentes y correspondencia con `sonar.sources` y reportPaths.
`go tool cover -func` también leyó el reporte (código 0, 73,8% redondeado).
`cmd/openapi/main.go` no tiene cobertura en este e2e y permanece en el alcance de Sonar;
no se añadió exclusión para elevar el porcentaje. `cobertura-go.json` conserva el resumen
y SHA256 del reporte, que queda ignorado por Git. Su compatibilidad de formato/rutas no
prueba ingestión remota. El 73,76% de sentencias Go no es la métrica de cobertura de código
nuevo de Sonar, ni acredita aprobar su 80%; no se rebajó el umbral.

| Imagen      | Build                            | Trivy HIGH/CRITICAL con fix | Dockle al umbral warn | Secretos en imagen |
| ----------- | -------------------------------- | --------------------------- | --------------------- | ------------------ |
| API         | 0, aprobado tras corrección Huma | 1, bloqueado: 14 HIGH       | 0                     | 0                  |
| Web SPA     | 0                                | 1, bloqueado: 12 HIGH       | 0                     | 0                  |
| Migraciones | 0                                | 1, bloqueado: 41 HIGH       | 0                     | 0                  |

`imagenes-validacion.json` conserva tags locales, IDs, avisos y códigos de salida finales.
Son recuentos por imagen, con repeticiones de stdlib entre ellas. Dockle informa de forma
nativa SKIP DKL-LI-0001 en distroless y mensajes INFO de Content Trust; migrate también
informa INFO de HEALTHCHECK por ser un proceso de una sola ejecución. No se agregaron
opciones de exclusión. Código 0 no significa que los controles nativamente SKIP se ejecutaran.

No ejecutados: `docker compose up` completo, scanner Sonar autenticado y GitHub Actions.
El e2e sí utiliza PostgreSQL real y migraciones; no equivale a probar la orquestación Compose.
Los valores de los secrets de GitHub no se leyeron; su disponibilidad se comprobará al ejecutar CI.

## Sonar y seguridad

[Comparación exacta del gate](../sonar-das.md): 80% es más estricto que el 60% del DAS y se
conserva; duplicación <=3% permite exactamente 3%, que el DAS prohíbe. Ratings A de bugs y
vulnerabilidades equivalen a cero, hotspots=100%, y mantenibilidad A es una condición extra.
No se cambió el gate. New Code, Automatic Analysis Off y tolerancia de cambios pequeños
requieren comprobación aparte; no se afirma que estén configurados ni que falten.

[Propuesta para aprobación](../actualizaciones-pendientes.md): **18 avisos alcanzables de
govulncheck = 16 de stdlib + 2 de x/text**, no 18 dependencias. Se conservan versiones DAS;
Go, módulos de Goose, x/text, Vite y Vitest se detallan con fuentes y archivos afectados.
`pnpm-audit.json` reproduce el resultado real actual, código 1 (1 HIGH, 6 moderate, 1 low).

## Revisión antes del commit

Los checks se llaman exactamente **calidad**, **seguridad**, **imagen**, tanto en ID como en
`name`; actionlint valida el workflow. En esta validación local todavía no existían checks remotos.
Revisión de archivos nuevos y diff: no hay `.tools`, `.tmp`, node_modules, binarios,
coverage.out ni `.env` reales en el conjunto versionable. Gitleaks de ese conjunto terminó
0, sin secretos; las credenciales de prueba son efímeras y explícitas. Las evidencias públicas
se limitaron a resultados técnicos, sin tokens ni rutas personales.

Benja debe revisar la base mínima: `apps/api/cmd/{api,spa-server,openapi}`, `apps/api/internal`,
`apps/api/e2e` y goldens, `apps/api/db/migrations/00001_base_tecnica.sql`, los módulos Go,
`apps/web/src`, tests y configuración del frontend, los tres Dockerfiles, `docker-compose.yaml`,
`.env.example`, `Makefile` y el contrato `contracts/openapi.yaml`/tipos generados. Es base
infraestructural: no sustituye dominio, autorización, esquema/queries sqlc ni seed de negocio.
Nombre y N del servicio siguen pendientes; no se deducen del grupo 5.

`revision-entrega.json` registra branch, HEAD, nombres de checks y comparación de locks:
Go mod/sum, catalog pnpm y lockfile no cambiaron respecto de la entrega anterior.
La captura corresponde al estado anterior al commit, staging y push; no se modificaron gate ni rulesets.

Verificación integrada de cierre: `pnpm verify` código 0 (formato, lint, tipos, Go mod/vet/test,
Vitest y builds). Mutación readiness y actionlint repetidos tras la corrección: código 0.
Una ejecución intermedia de verify detectó formato en el JSON mientras el scanner aún lo
escribía; se formateó la evidencia terminada y se repitió verify completo con éxito.

Fallos de entorno resueltos durante la validación: PowerShell/CMD requerían ejecución con
permisos del entorno por error Windows 1260; una copia recursiva nativa de Node terminó con
error Windows 3221226505. Se sustituyó por lectura/escritura explícita y la mutación pasó.
El primer fixture de secreto era de longitud incorrecta y no fue detectado; el fixture final
usa el formato correcto y el scanner lo bloqueó. No se presenta aquel intento como aprobado.

Los archivos `govulncheck.txt`, `trivy-fs.json`, `oasdiff-negativo.txt`, `versiones.json`,
`imagenes.json` y `sonar-gate.json` contienen evidencia local/consultas oficiales, no evidencia
de un pipeline remoto verde. Los resultados finales resumidos se guardan en `comprobaciones.json`.

## Correcciones y Compose — 9 de noviembre de 2026

S6549: se confirmó el flujo URL.Path → os.Stat → FileServer. Con un enlace simbólico
preexistente dentro del bundle hacia un archivo externo, el binario anterior devolvió el
fixture privado con HTTP 200. El corregido rechaza ese bundle al iniciar (salida 1), sin
exponer el contenido. Evidencia: `spa-regresion.json`. El cambio carga archivos regulares
mediante os.Root al arrancar y sirve un catálogo inmutable; una petición no consulta el
filesystem. El bundle debe ser pequeño, confiable y desplegarse completo: modificar archivos
requiere reiniciar. Se mantienen assets, fallback SPA, proxy, GET/HEAD y cabeceras de seguridad.

Compose completo aprobado: proyecto aislado `p5-compose-2bac0ae9`, PostgreSQL saludable,
migración terminada con código 0, backend/frontend saludables; healthz, readyz, HTML y
`/api/openapi.yaml` a través del frontend dieron 200 con contenido comprobado. Puertos
provisionales 51781 (web) y 51782 (backend), no asignación oficial. Se usó API_ENV=dev para
verificar el contrato a través del proxy. El perfil opcional Mailpit no se activó.

El primer intento no pudo crear red: Docker agotó sus pools automáticos. Se comprobó que
las subredes 10.253.240.0/24 y 10.253.241.0/24 no coincidían con redes Docker ni rutas locales
específicas y se usaron overrides temporales exclusivamente para validación. No se borró
ninguna red o volumen ajeno. `compose down --volumes` se ejecutó solo sobre el nombre único
creado por esta prueba. Evidencia: `compose-correcciones.json`.

Sonar antes del cambio: 21 líneas sin cobertura en cmd/api, 14 en cmd/spa-server, 4 en
cmd/openapi y 2 en internal/spa. El frontend no reportaba líneas ni condiciones pendientes.
Se añadieron escenarios de procesos reales para configuración inválida, puertos ocupados,
healthchecks exitosos/no saludables/inaccesibles, bundle inválido, traversal, HEAD/POST,
assets ausentes e inmutabilidad. El generador OpenAPI se instrumenta y su salida se compara
con el contrato versionado; no se excluyó código ni se modificó el umbral.

La migración y las nueve pruebas e2e terminaron con código 0 sobre PostgreSQL real.
Cobertura Go: 171/180 sentencias (95%); las cinco fuentes de producción están presentes,
sin archivos omitidos. El validador de rutas/formato pasó. Esta cifra local no sustituye
la cobertura nueva medida por Sonar. `pnpm verify` y Gitleaks sobre apps/api pasaron.
Evidencias: `docker-correcciones.json` y `cobertura-correcciones.json`.

## Validación remota de las correcciones

[Actions del commit 8fad783](https://github.com/danieldiaz77-LAB/pingeso-grupo5/actions/runs/42739207546): **calidad aprobado**, **seguridad fallido**,
**imagen omitido**. Calidad ejecutó verify, mutación de readiness, contrato, e2e PostgreSQL,
e2e navegador, scanner y Quality Gate con éxito. Sonar atribuye el análisis al SHA
`8fad7834f28042a808082ffbb3ff6270770c7b4e`: cobertura nueva **95,1%**, gate **OK**, ratings A,
duplicación 0% y hotspots revisados 100%. No hay incidencias abiertas; S6549 dejó de
reportarse sin exclusiones ni resolución manual. `ignoredConditions=false` en este análisis.

Seguridad sigue fallando en govulncheck (19 avisos: 17 stdlib + 2 x/text), audit (8 avisos)
y Trivy filesystem (2 HIGH). Gitleaks y Hadolint pasaron. El aviso adicional GO-2026-6604
se explica en actualizaciones pendientes; la propuesta de Go 1.26.10 lo cubre. No se cambiaron
versiones, gate ni dependencias entre jobs. Imagen **no está validado en CI**.
La limpieza de los recursos locales exclusivos terminó con éxito. El lector oficial
`go tool cover` confirmó 95% local. Evidencia remota: `remoto-correcciones.json`.
