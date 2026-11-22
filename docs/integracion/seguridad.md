# Decisiones de dependencias y seguridad

No se aplicaron las actualizaciones pendientes. `go.mod`, `go.sum`,
`pnpm-workspace.yaml` y `pnpm-lock.yaml` conservan la base del pipeline.

## Diferencias verificadas en las dos fuentes

`dependencias.json` enumera los módulos añadidos por Benja y cada cambio común:
puddle/v2 2.2.2 → 2.2.3, x/sync 0.20.0 → 0.24.0 y x/text 0.37.0 → 0.43.0.
Su grafo también añade x/crypto 0.58.0 y compress 1.20.1, entre otros.
No se ha adoptado ese grafo ni probado su compatibilidad.

El backend integrado aquí usa las dependencias ya existentes (Huma y pgx).
Las pruebas HTTP adaptadas usan pgx contra la base efímera del runner; ya no requieren
importar Goose ni Testcontainers dentro del módulo API. Goose sigue aplicando
migraciones antes del e2e, con su versión independiente fijada.

Autenticación/configuración/middleware de Benja NO están activados en esta copia:
requieren incorporar keyfunc/jwt y revisar su resolución transitiva. Sus fuentes
originales se conservan en la carpeta hermana `*-fuentes`.
Además, su wrapper de autenticación abarca salud; antes de adoptarlo hay que conservar
healthz/readyz públicos y comprobar fallos de JWKS. La copia no está lista para producción.

Frontend: se añaden al package.json únicamente paquetes ya fijados en el catalog del
pipeline. Se conservan @types/react 19.2.2 y @types/react-dom 19.2.2; no se importan
19.2.8/19.2.3 ni @types/node 22.20.5 de Benja. No se incorpora react-router porque
la app de demostración no necesita rutas adicionales.

**Lockfile pendiente:** no se pudo ejecutar pnpm. El lock heredado no contiene los
nuevos importers de la plantilla: `pnpm install --frozen-lockfile` debe fallar hasta
resolverlo. No se fabricó un lock combinando snapshots con peers incompatibles.
En una terminal habilitada, generar un lock candidato en el worktree, revisar todas
las versiones transitivas y detenerse si exige una actualización no aprobada.
Los locks de ambas fuentes están disponibles para comparación, no para reemplazo ciego.

## Propuesta separada por alcance

Las versiones siguientes son propuestas históricas de
`docs/actualizaciones-pendientes.md`, NO una comprobación nueva de avisos o fixes.

| Alcance         | Base / propuesta histórica                                                                                     | Validación tras aprobación                                                                |
| --------------- | -------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| API y toolchain | Go 1.26.5 → 1.26.10; x/text 0.37.0 → 0.41.0. Contrastar con el 0.43.0 de Benja antes de resolver el grafo      | go mod verify, go vet, tests, e2e, govulncheck v1.9.0, reconstruir API y SPA              |
| Frontend        | Vite 8.1.5 → 8.1.6; Vitest/coverage-v8/mocker 4.1.10 → 4.1.12                                                  | frozen-lockfile, verify, navegador, pnpm audit, Trivy fs                                  |
| Goose/migrate   | Mantener Goose v3.24.1; proponer módulo de herramientas independiente y reproducible para corregir transitivas | compilar migrador desde ese módulo, aplicar migraciones en BD nueva, Trivy image y Dockle |

Para Goose: la tabla histórica propone jwt/v4 4.5.2, pgx/v5 5.9.0,
compress 1.18.7, x/crypto 0.55.0, x/net 0.62.0, x/text 0.41.0 y grpc 1.83.2.
Revalidar avisos, disponibilidad y compatibilidad antes del PR de actualización.
Modificar solamente apps/api/go.mod no modifica el binario instalado con
`go install github.com/pressly/goose/v3/cmd/goose@v3.24.1`.

Fuentes y avisos previos: docs/actualizaciones-pendientes.md y sus JSON de evidencias.
No hay un escaneo nuevo sobre esta preparación; los recuentos anteriores no se reutilizan.
El jefe debe aprobar y registrar las versiones en el DAS mediante un PR separado.
No cambiar exclusiones, severidades, Quality Gate ni el mínimo de cobertura nueva de 80%.
