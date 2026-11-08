# Pipeline de calidad y seguridad — Hito 0

Estado: implementación preparada para revisión en PR en borrador. **Hito no cerrado.**
Las evidencias locales no sustituyen la ejecución de Actions ni el Quality Gate remoto.

## Alcance y referencia

La referencia es el [DAS](referencias/DAS%20—%20Desarrollo%20de%20Microservicios.md), conservado
sin modificaciones. El repositorio comenzó con README y el DAS local sin seguimiento.
Los controles se implementan en `feat/pipeline-hito0` sobre el commit inicial
`ed6f74d33d696eb4cafe807291f3f891451928b1`, con destino `main` del repositorio
`danieldiaz77-LAB/pingeso-grupo5`. Las protecciones y los valores de secrets no se modifican
como parte de esta entrega.

## Archivos principales

- `.github/workflows/ci.yml`: eventos, tres nombres exactos, mínimos permisos y dependencias.
- `scripts/tareas.mjs`: mismos comandos locales/CI, compilación, contrato y cobertura real.
- `scripts/contrato.mjs`: comparación con la base correcta y primera línea base explícita.
- `sonar-project.properties` y `scripts/sonar.mjs`: fuentes, pruebas, reportes y bloqueo.
- `lefthook.yml`, `scripts/arquitectura.mjs`: controles anticipados repetidos por CI.
- `pnpm-lock.yaml`, `apps/api/go.sum` y `docs/evidencias/`: reproducibilidad y verificación.
- Dockerfiles y Compose: imágenes fijadas, no root y orden de arranque.

## Controles

| Job       | Comandos y criterio                                                                                                                     |
| --------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| calidad   | `pnpm verify`: Prettier, oxlint sin warnings, restricciones de arquitectura, TypeScript, gofmt, go mod verify, go vet, pruebas y builds |
| calidad   | `pnpm codegen` y `pnpm check:generated`: contrato y tipos generados; bloquea diferencias respecto de HEAD y archivos nuevos             |
| calidad   | `pnpm check:contract`: oasdiff v1.35.0, falla en ERR                                                                                    |
| calidad   | `pnpm test:e2e:api`: binarios instrumentados, PostgreSQL real y golden; `go tool covdata textfmt` produce coverage.out                  |
| calidad   | Vitest genera `apps/web/coverage/lcov.info`; Playwright ejecuta Chromium real                                                           |
| calidad   | Sonar Cloud y Quality Gate bloqueante, sin saltar secrets ausentes                                                                      |
| seguridad | Gitleaks cero hallazgos, govulncheck v1.9.0 y pnpm audit HIGH                                                                           |
| seguridad | Hadolint nivel error y Trivy filesystem vuln,misconfig,secret HIGH/CRITICAL con fix                                                     |
| imagen    | Construye api, web y migrate; Trivy vuln,secret y Dockle warn en las tres                                                               |

`imagen` necesita `calidad` y `seguridad`. Los escáneres independientes posteriores siguen ejecutándose aunque falle otro;
`if: !cancelled()` no elimina ni convierte en éxito las fallas previas. Se deshabilitan
comentarios automáticos de Gitleaks para conservar contents:read.
Ningún job usa matrices, filtros de rutas,
continue-on-error ni tolerancia de fallos. Los PR anteriores se cancelan al recibir cambios;
los push a main no se cancelan. No se publica ninguna imagen en Hito 0.

## Instalación y comandos

Requiere Go 1.26.5, Node 22.22.0, pnpm 10.13.1, Git, Gitleaks 8.30.1 y Docker Linux activo.
Las copias portátiles descargadas para esta revisión están en `.tools/` (ignoradas).
En esta computadora, cargar `. ./scripts/entorno-local.ps1` en PowerShell. En Linux/CI las
Actions preparan las herramientas. No es necesario cambiar el Node global del usuario.

```sh
pnpm install --frozen-lockfile
pnpm exec lefthook install
pnpm verify
pnpm codegen
pnpm check:generated
pnpm --dir apps/web exec playwright install chromium
pnpm test:e2e:web
# Linux; PostgreSQL de pruebas desechable, nunca producción:
TEST_DATABASE_URL='postgres://usuario:clave@localhost:5432/prueba?sslmode=disable' pnpm test:e2e:api
pnpm audit --audit-level high
cd apps/api
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

Los hooks requieren instalación explícita porque `.npmrc` deshabilita scripts de instalación.
El e2e termina procesos con señales POSIX para volcar cobertura: ejecutarlo en Linux/WSL/CI.
No reemplaza la prueba funcional de dificultad media/alta que hará el equipo.

Copiar `.env.example` a `.env`, completar todos los `replace-with-*`, acordar nombre y N,
y luego `docker compose up -d --build`. No asumir N=5. `WEB_PORT=51N0`, `BACKEND_PORT=51N1`.
`DATABASE_URL` debe apuntar a `postgres` dentro de Compose. No se publican puertos de BD.
`AUTH_JWKS_URL` queda documentado para la base funcional: salud es pública y esta base no
implementa endpoints autenticados. No se inventa un servicio de identidad.

## Primer contrato

La base de PR es `github.event.pull_request.base.sha`; en push es `github.event.before`.
Checkout obtiene el historial completo. El script primero exige que el SHA resuelva a un
commit; después inspecciona su árbol. Solo un árbol válido sin `contracts/openapi.yaml`
se considera primera línea base. Errores de Git, SHA inexistente o SHA cero fallan.

El primer contrato se genera desde Huma y openapi-typescript lo lee/valida para producir
tipos. Se informa expresamente que no existe una versión previa para comparar. En cambios
posteriores oasdiff es obligatorio. En push se compara con el commit previo, nunca consigo
mismo. Los generados deben incluirse en el commit: `check:generated` falla deliberadamente
mientras sean archivos nuevos sin seguimiento. No hacer `git add` automático para taparlo.

## SonarQube Cloud

Claves del proyecto (respetar mayúsculas):

- `sonar.projectKey=danieldiaz77-LAB_pingeso-grupo5`
- `sonar.organization=danieldiaz77-lab`
- secret `SONAR_HOST_URL=https://sonarcloud.io/` y secret `SONAR_TOKEN`.

El token no se guarda ni se imprime. Las fuentes Go están en cmd/internal y las pruebas en
e2e; frontend src y tests/e2e están separados. Solo se excluyen los dos directorios de código
generado permitidos por DAS. No hay exclusiones de cobertura para elevar artificialmente
el porcentaje. La cobertura de Go viene de binarios reales, no de un archivo inventado.

La API pública confirmó el proyecto y la organización. También confirmó que está asignado
**Sonar way**, con cobertura nueva >=80%, hotspots revisados=100%, duplicación <=3% y
ratings de seguridad, fiabilidad y mantenibilidad A (evidencia en sonar-gate.json).
No se modificó ese gate ni se rebajó su 80%. Su límite de duplicación acepta exactamente 3%,
mientras el DAS dice menor que 3%; el jefe debe resolver esa diferencia ; los ratings A consultados equivalen a cero bugs y vulnerabilidades. No se acredita todavía cumplimiento exacto del DAS.

Configuración externa pendiente de comprobar:

1. Desactivar Automatic Analysis y usar análisis por CI.
2. Confirmar que main es rama principal y que el plan admite análisis de PR.
3. Definir New Code antes del primer análisis (PR: diferencias frente a la rama destino;
   main: definición acordada, por ejemplo versión anterior). Registrar cuál se eligió.
4. Gate sobre código nuevo: bugs=0, vulnerabilidades=0, hotspots revisados=100%,
   duplicación **menor que 3%**, cobertura DAS **al menos 60%**, conservando el gate actual de **80%**. Desactivar la tolerancia de
   cobertura/duplicación para cambios pequeños si la instancia la tiene; no rebajar el DAS.
5. Confirmar secrets con `gh secret list --repo danieldiaz77-LAB/pingeso-grupo5` cuando gh
   esté instalado y autenticado. Solo listar nombres.

Si Cloud exige una línea base de main antes de analizar el primer PR, registrar primero un
análisis de la rama principal desde una copia limpia del main real con un scanner oficial,
el proyecto correcto y su Quality Gate; o resolver el arranque con administración de Sonar.
No etiquetar código de la feature como main, no omitir el gate y no desactivar protecciones.
La solución depende del error real que entregue Cloud; todavía no se ha observado ese error.
El escáner local previo al PR queda pendiente de un token en el entorno (sin incluirlo en archivos versionados).

## Versiones y discrepancias verificadas

`docs/evidencias/versiones.json` contiene consultas a npm, proxy Go y API GitHub oficiales.
`docs/evidencias/imagenes.json` registra índices multi-arquitectura de los registros.

- pnpm/action-setup **v6.1.0**: se usa el SHA oficial
  `ea17c68df8912ef543352723c149a84f56e3d413`, distinto del ejemplo DAS.
- trivy-action **v0.36.0**: se usa el SHA oficial
  `ed142fd0673e97e23eac54620cfb913e5ce36c25`, distinto del ejemplo DAS.
- Distroless: el tag flotante ahora apunta a otro digest, pero se comprobó HTTP 200 e índice
  de seis plataformas para el digest histórico del DAS y se conserva ese digest exacto.
- `openapi-typescript 7.13.0` declara peer TypeScript `^5.x`; el DAS fija `6.0.3`.
  Instalación avisa. Generación y comprobación de tipos se ejecutan con esas versiones;
  no se sustituyen ni se oculta la advertencia. Solicitar aclaración al jefe.
- No se incorporan River, correo, JWT ni componentes de negocio que este esqueleto no usa.
  Las versiones opcionales del DAS siguen siendo la referencia al incorporarlos.
- La BD de CI usa un service PostgreSQL real con TEST_DATABASE_URL, alternativa permitida
  por la sección de pruebas; no se añade Testcontainers como dependencia sin uso.

## Bloqueo de seguridad confirmado: Vite

Se conserva **Vite 8.1.5** según el DAS, pendiente de aprobación del jefe para actualizarlo.
`pnpm audit --audit-level high` terminó con código 1: 8 avisos (1 HIGH, 6 moderate, 1 low).
El bloqueo HIGH es [GHSA-587h-7wmj-h5pm](https://github.com/advisories/GHSA-587h-7wmj-h5pm),
CVE-2026-112942. La fuente oficial confirma que el rango de Vite `>=8.0.0 <=8.1.5` está
expuesto y la corrección de esa serie es **8.1.6**. El origen es launch-editor en Windows.
No se ignora por correr el CI en Linux: el umbral de audit del DAS sigue siendo obligatorio.

Propuesta mínima para aprobación: cambiar solo Vite a 8.1.6 en el catalog, regenerar el
lockfile, actualizar la tabla del DAS por el procedimiento del jefe y ejecutar todo el CI.
No se aplicó. Hasta resolverlo `seguridad` falla y `imagen` no puede ejecutarse remotamente.
Las comprobaciones locales independientes de calidad pueden continuar.

## Responsables y límites

**Benja:** revisar todo `apps/api/`, `apps/web/`, Dockerfiles, Compose y Makefile. Confirmar
nombre y N; revisar la migración técnica que valida PostgreSQL 17 (sin entidades), incorporar esquema, sqlc.yaml/queries, migraciones de negocio y seed dev cuando se
conozca el dominio. `make generate` queda preparado con sqlc fijado, pero no se invoca en CI
porque no existe esquema ni consultas: no se inventan tablas. No hay seed ficticio. La mutación automatizada elimina el ping de readiness solo en una
copia temporal y exige un fallo concreto de la prueba aislada; el e2e real sigue siendo obligatorio.
Revisar validación completa de configuración, autenticación/JWKS y estructura de negocio
Handler → Service → Store al agregar funcionalidades; salud es infraestructura.

**Manuel:** contrato funcional, structs/validaciones y mocks antes de lógica; regenerar YAML
y tipos. El contrato de salud inicial no sustituye esa entrega.

**Todos:** escenarios funcionales e2e con golden, matriz de permisos por rol, idempotencia,
reintentos, propagación JWT/request-id y pruebas con fallas deliberadas. Probar fechas también
en UTC y Pacific/Auckland cuando exista lógica temporal; aquí no hay cálculos de fechas.

**Daniel / jefe:** aprobación de las actualizaciones detalladas en `actualizaciones-pendientes.md`, aclaración del peer TS, verificación de secrets,
Automatic Analysis, New Code y gate, ejecución remota y revisión del compañero.

Publicación, SBOM CycloneDX, firma Cosign 3.1.4, GHCR y paquete npm corresponden a releases.
No se habilitan permisos packages:write/id-token:write en este CI. Una futura release debe
usar los artefactos ya escaneados, adjuntar contrato y SBOM y firmar; no reconstruirlos.

## Entrega posterior a revisión local

Procedimiento de publicación, tras revisar el diff y los secretos:

```sh
git status --short
git add README.md .github .gitignore .gitattributes .dockerignore .editorconfig .npmrc .nvmrc .prettierrc.json .prettierignore .env.example package.json pnpm-workspace.yaml pnpm-lock.yaml lefthook.yml sonar-project.properties Makefile docker-compose.yaml apps contracts scripts docs
git diff --cached --check
gitleaks protect --staged --redact --no-banner
git commit -m "ci: configurar pipeline de calidad y seguridad del hito 0"
git push -u origin feat/pipeline-hito0
gh pr create --repo danieldiaz77-LAB/pingeso-grupo5 --base main --head feat/pipeline-hito0 --draft --title "ci: configurar controles de calidad y seguridad del Hito 0" --body-file docs/pr-hito0.md
```

Con Vite sin aprobación el PR puede servir para revisión, pero no estará listo para integrar.
No solicitar merge hasta los tres checks verdes, Sonar comprobado, conversaciones resueltas
y aprobación de otro integrante. Añadir al README enlaces a la ejecución y al análisis
concretos cuando existan; una página genérica de Actions no es evidencia de ejecución verde.

## Bloqueos adicionales confirmados en Go

`go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...` terminó con código 1
(el ejecutable interno informó `exit status 3`). Reportó **18 vulnerabilidades alcanzables**
en la biblioteca estándar Go 1.26.5 y `golang.org/x/text v0.37.0`, además de seis avisos no
alcanzables. Salida completa en [govulncheck.txt](evidencias/govulncheck.txt).

Las correcciones máximas indicadas para los avisos alcanzables son Go **1.26.10** y
x/text **v0.41.0**. Se verificaron los avisos en la base oficial Go:
[GO-2026-6809](https://pkg.go.dev/vuln/GO-2026-6809) y
[GO-2026-6629](https://pkg.go.dev/vuln/GO-2026-6629).
Son propuestas para revisión del jefe, no actualizaciones aplicadas. Actualizar Go también
requiere verificar y fijar nuevos digests de build, repetir builds y escanear runtimes.

Trivy v0.76.1 confirmó dos HIGH con fix en x/text: CVE-2026-56851 (0.41.0) y
CVE-2026-56852 (0.39.0). No se agregaron excepciones ni se bajaron umbrales.
Aunque se apruebe Vite, estos hallazgos también deben resolverse antes de un check seguridad
verde. Trivy excluye dependencias de desarrollo por defecto; pnpm audit revisa esas
dependencias y detecta Vite. Se conservan ambos controles.

La comparación exacta del gate está en [sonar-das.md](sonar-das.md). La validación Docker y
los bloqueos actuales se registran en [validacion-local.md](evidencias/validacion-local.md);
la propuesta de versiones completa está en [actualizaciones-pendientes.md](actualizaciones-pendientes.md).
