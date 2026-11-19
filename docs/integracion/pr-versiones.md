# Borrador de PR: actualización de versiones por seguridad (Hito 0)

Estado: **borrador para aprobación del jefe. No se ha cambiado ninguna versión.**
Fecha de los escaneos: 19/11/2026. PR separado de la integración, como exige el DAS, que además debe actualizar su tabla de versiones.

## Por qué

Con las versiones actuales el pipeline no puede quedar verde. Los jobs `seguridad` e `imagen` fallan por avisos HIGH/CRITICAL que ya tienen corrección publicada:

| Escaneo (integración local, 19/11)        | Resultado                                                                 | Evidencia en `docs/evidencias/`                      |
| ----------------------------------------- | ------------------------------------------------------------------------- | ---------------------------------------------------- |
| govulncheck v1.9.0, API                   | 21 alcanzables (19 stdlib, 2 x/text)                                      | govulncheck-api-integracion-2026-11-19.txt           |
| govulncheck v1.9.0, binario Goose v3.24.1 | 32 alcanzables (7 módulos + stdlib)                                       | govulncheck-goose-v3.24.1-integracion-2026-11-19.txt |
| pnpm audit --audit-level high             | 1 high, 6 moderate, 1 low (Vite, Vitest)                                  | pnpm-audit-integracion-2026-11-19.\*                 |
| Trivy fs                                  | 2 HIGH (x/text)                                                           | trivy-fs-integracion-2026-11-19.\*                   |
| Trivy imagen api / web / migrate          | stdlib en las 3; x/text en api; grafo de Goose en migrate. SO base limpio | imagenes-integracion-20261119151924/                 |
| Hadolint, Dockle, Gitleaks                | Aprobados                                                                 | hadolint-…, imagenes-…/salida.txt, gitleaks-…        |

Los paquetes del sistema operativo de las imágenes base (Debian 13.6 distroless, Alpine 3.23.7) no tienen hallazgos. Todo lo que hay que corregir son versiones de Go y de módulos.

## Cambios propuestos

| #   | Componente                                    | Actual → propuesto                                                                                                                        | Corrige                                                                                                                                                                          | Dónde se fija                                                                                                                                                                                                                                                                                                                  | Riesgo de compatibilidad                                                                                                                                                                                                                   |
| --- | --------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1   | Go (toolchain y stdlib)                       | 1.26.5 → **1.26.11**                                                                                                                      | 19 avisos de stdlib en la API y 13 CVE HIGH en las 3 imágenes. El máximo requerido es 1.26.11 por GO-2026-6859 / CVE-2026-97054; **1.26.10, la propuesta anterior, ya no basta** | `apps/api/go.mod` (`toolchain go1.26.11`; `go 1.26.0` no cambia); `.github/workflows/ci.yml` (`go-version`, 2 jobs); `apps/api/Dockerfile`, `apps/api/Dockerfile.migrate` y `apps/web/Dockerfile` (`golang:1.26.11-alpine3.23@sha256:<digest por verificar>`); README y `docs/pipeline-hito0.md`; DAS; Go portátil en `.tools` | **Bajo.** Es un parche de la misma rama 1.26 y no cambia el lenguaje. El digest de la imagen oficial hay que verificarlo al aprobar (no se inventa). Hay que reconstruir los binarios: actualizar solo el Go local no corrige las imágenes |
| 2   | golang.org/x/text (API, indirecta)            | v0.37.0 → **v0.41.0**                                                                                                                     | GO-2026-6629 / CVE-2026-56851 (fix 0.41.0) y GO-2026-5970 / CVE-2026-56852 (fix 0.39.0); los 2 HIGH de Trivy fs y de la imagen api                                               | `apps/api/go.mod`, `apps/api/go.sum` (`go get golang.org/x/text@v0.41.0` + `go mod tidy`)                                                                                                                                                                                                                                      | **Bajo-medio.** Llega por pgx/huma. Hay que revisar si MVS sube otras transitivas compartidas (x/sync, x/sys) y registrarlas                                                                                                               |
| 3   | Vite                                          | 8.1.5 → **8.1.6**                                                                                                                         | GHSA-587h-7wmj-h5pm (HIGH) y 4 avisos más del rango ≥8.0.0 ≤8.1.5                                                                                                                | `pnpm-workspace.yaml` (catalog), `pnpm-lock.yaml`                                                                                                                                                                                                                                                                              | **Bajo.** Parche; según la revisión previa acepta Node 22.22.0. Solo es herramienta de dev/build: la imagen web sirve los estáticos con el servidor Go                                                                                     |
| 4   | Vitest                                        | 4.1.10 → **4.1.12**                                                                                                                       | GHSA-82fw-gwwq-j7x9 (moderate)                                                                                                                                                   | `pnpm-workspace.yaml`, `pnpm-lock.yaml`                                                                                                                                                                                                                                                                                        | **Medio en la resolución.** Los overrides de §9 (inquirer, @types/node, tldts…) y `engine-strict` pueden chocar al regenerar el lock; hay que revisar el diff del lock                                                                     |
| 5   | @vitest/coverage-v8 (arrastra @vitest/mocker) | 4.1.10 → **4.1.12**                                                                                                                       | @vitest/mocker: GHSA-82fw-gwwq-j7x9 y GHSA-qxg6-89mg-rqv8                                                                                                                        | `pnpm-workspace.yaml`, `pnpm-lock.yaml`                                                                                                                                                                                                                                                                                        | Igual que el #4: se sube alineado con Vitest                                                                                                                                                                                               |
| 6   | Grafo del migrador Goose                      | Binario v3.24.1 con su go.sum publicado → **mismo Goose v3.24.1 compilado desde un módulo de herramientas propio** con mínimos corregidos | Los 32 de govulncheck Goose y los hallazgos de Trivy en migrate (tabla siguiente)                                                                                                | Nuevo `tools/goose/go.mod` y `go.sum`; `apps/api/Dockerfile.migrate`; `scripts/tareas.mjs:81`; `Makefile:2` (`GOOSE_VERSION`); DAS                                                                                                                                                                                             | **Medio-alto.** Goose 3.24.1 no se probó con esos drivers más nuevos. Puede no compilar o fallar con ClickHouse/YDB/gRPC. Hay que compilarlo, probarlo con el e2e real y con Compose, y revisar el grafo final                             |

### Mínimos para el grafo de Goose (#6)

Se toma la mayor versión de corrección que piden govulncheck y Trivy.

| Módulo                        | Actual | Mínimo  | Avisos                                                                                                                                     |
| ----------------------------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| github.com/jackc/pgx/v5       | 5.7.1  | 5.9.2   | CVE-2026-33815, CVE-2026-33816 (CRITICAL), GO-2026-4771, GO-2026-5004                                                                      |
| google.golang.org/grpc        | 1.62.1 | 1.83.2  | CVE-2026-33186 (CRITICAL), CVE-2026-84304, CVE-2026-84445, GHSA-hrxh-6v49-42gf, GO-2026-6061/6348/6443, GO-2026-4762                       |
| golang.org/x/net              | 0.33.0 | 0.62.0  | 7 CVE HIGH en la imagen (incluye CVE-2026-78661); govulncheck: GO-2026-5026/5942/6603…6617                                                 |
| golang.org/x/crypto           | 0.31.0 | 0.55.0  | 12 CVE HIGH en la imagen (CVE-2025-22869, CVE-2025-47913, CVE-2026-39828…46597, CVE-2026-56854); govulncheck no los marca como alcanzables |
| golang.org/x/text             | 0.21.0 | 0.41.0  | CVE-2026-56851, CVE-2026-56852                                                                                                             |
| github.com/klauspost/compress | 1.17.7 | 1.20.2  | CVE-2026-63209 (fix 1.18.7), GO-2026-6865 (fix 1.20.2)                                                                                     |
| github.com/golang-jwt/jwt/v4  | 4.5.1  | 4.5.2   | CVE-2025-30204 / GO-2025-3553                                                                                                              |
| github.com/andybalholm/brotli | 1.1.1  | 1.2.7   | GO-2026-6888 (solo govulncheck)                                                                                                            |
| github.com/ClickHouse/ch-go   | 0.61.5 | 0.65.0  | GO-2025-3603 (solo govulncheck)                                                                                                            |
| stdlib                        | 1.26.5 | 1.26.11 | Mismo cambio que el #1                                                                                                                     |

Alternativas para el #6 que conviene evaluar antes de aprobar, sin elegir ninguna aquí:

- **A, recomendada para empezar: módulo de herramientas.** Fija Goose v3.24.1 y esos mínimos. Mantiene la versión aprobada en el DAS y solo cambia sus dependencias.
- **B: una versión más nueva de Goose** cuyo `go.mod` ya cumpla los mínimos. Hay que identificar y verificar la versión concreta; aquí no se propone ninguna sin comprobarla.
- **C: compilar Goose solo con el driver de PostgreSQL**, mediante build tags de exclusión de drivers si la versión los soporta (verificar). Quitaría grpc, ch-go, brotli y otras dependencias del binario y reduciría la superficie, aunque añade una configuración de build propia.

## Fuera de alcance

- No se toca pgx de la API (ya está en v5.10.0), ni jwt v5.3.1 / keyfunc v3.9.1, que no tienen avisos.
- No se cambian Node 22.22.0, pnpm 10.13.1, TypeScript 6.0.3 ni las imágenes distroless/alpine de runtime, que no tienen hallazgos.
- No se suben otras dependencias por sugerencia de consola (por ejemplo, el aviso de Trivy 0.76.2).
- Ningún aviso se silencia con ignores ni exclusiones.

## Validación después de aplicar (en el PR)

1. `go mod tidy` y `pnpm install --frozen-lockfile` tras regenerar los locks; revisar en el diff de `go.mod`, `go.sum` y `pnpm-lock.yaml` cada transitiva que cambie.
2. `pnpm verify`, `pnpm codegen` (doble hash), `test:mutation`, `check:generated` y `check:contract`.
3. e2e con binario real y PostgreSQL 17.11, `check:coverage`, Playwright y Compose integrado (los mismos scripts de esta integración).
4. govulncheck (API y Goose), `pnpm audit --audit-level high`, Trivy fs, rebuild de las 3 imágenes, y Trivy de imagen y Dockle.
5. El PR solo se da por bueno si los jobs `calidad`, `seguridad` e `imagen` quedan verdes en CI. No se promete verde antes de ejecutarlos.

## Lo que necesitamos del jefe

- Aprobar cada fila (#1–#6) o indicar cuáles se posponen.
- Para el #6, elegir A, B o C, o pedir que se evalúen primero.
- Actualizar el DAS (Go 1.26.11, x/text, Vite, Vitest y la forma de compilar Goose) o autorizar que el PR lo haga.

Referencia histórica: `docs/actualizaciones-pendientes.md` (proponía Go 1.26.10; queda superado por el #1).
