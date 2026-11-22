> **Actualización 19/11/2026:** este documento conserva información histórica de preparación. El worktree y el merge ya existen; contrato/lock fueron regenerados y hay validaciones posteriores. Para el estado vigente y próximos pasos, leer [continuidad-claude.md](continuidad-claude.md). No volver a ejecutar el script de preparación ni interpretar las tablas antiguas como estado actual.

# Incorporación del contrato funcional

Esta preparación integra únicamente ejemplos genéricos. No define recursos, permisos,
requests ni responses funcionales de Manuel. No acredita entrega del Hito 1.

## Flujo único

- Fuente: structs y registro de operaciones Huma en `apps/api/internal/api/`.
- Composición compartida: `api.Nueva`, usada por API real y `apps/api/cmd/openapi/main.go`.
- Generador único: `cmd/openapi`. No se incorpora `cmd/contrato`.
- Contrato único: `contracts/openapi.yaml`.
- Tipos únicos: `apps/web/src/api/schema.d.ts`, generados con openapi-typescript.
- No se incorpora el esquema manual `esquema.ts` ni la salida `contrato.ts`.
- `pnpm codegen` ejecuta sqlc, Huma y openapi-typescript.
- `pnpm check:generated` exige seguimiento Git, archivos no vacíos y ausencia de
  diferencias contra HEAD, incluidos archivos nuevos. CI regenera antes de comprobar.
- `pnpm check:contract` compara contra `CONTRACT_BASE_SHA` con oasdiff v1.35.0.

**Bloqueo de esta copia:** no se pudo ejecutar Go ni pnpm. Los dos generados conservados
son los de la base técnica anterior y NO representan todavía los endpoints de ejemplos.
La guarda de contrato falla deliberadamente mientras no se regenere. No editar los
generados a mano ni publicarlos como contrato integrado.

## Cuando Manuel publique su entrega

1. Obtener URL de rama/PR y SHA exacto, con endpoints, permisos, errores y mocks acordados.
   Revisar sus cambios contra esta base; no copiar carpetas completas.
2. Incorporar requests/responses y validaciones en archivos de
   `apps/api/internal/api/`. Registrar sus operaciones en la composición común
   utilizada por servidor y generador. Conservar salud y ejemplos identificados como plantilla.
3. Incorporar handlers de respuestas simuladas acordadas para Hito 1. Mantener la
   separación Handler → Service → Store cuando se incorpore persistencia real.
4. Añadir mocks en `apps/web/src/mocks/`, importando tipos desde `schema.d.ts`
   o alias derivados de `paths`. Los clientes de cada microcomponente van en
   `*.api.ts`; reciben la URL base. No duplicar schemas escritos a mano.
5. Completar autenticación y matriz de permisos con lo acordado, sin deducirlos de
   los ejemplos. La variante auth/config de Benja permanece en la copia de fuentes,
   pendiente de resolver su grafo y probar salud pública/JWKS inválido.
6. En terminal habilitada y dependencias resueltas:
   `pnpm install --frozen-lockfile`, `pnpm codegen`,
   `gofmt -w apps/api`, `pnpm format:write`, `pnpm verify`.
   Revisar `git diff`, añadir generados y fuentes, y hacer un commit SOLO local
   revisable. Después ejecutar `pnpm check:generated` contra ese HEAD.
7. Fijar `CONTRACT_BASE_SHA` al commit previo de integración y
   `OASDIFF_VERSION=v1.35.0`; ejecutar `pnpm check:contract`.
   No usar automáticamente main inicial si ya existe un contrato integrado previo.
8. Ejecutar mutación, e2e instrumentados con base efímera propia, pruebas navegador,
   cobertura y escaneos. Añadir escenarios funcionales de Manuel y rechazos por rol.
   No reemplazar los e2e de salud/SPA ni presentar tests con mocks como PostgreSQL real.
9. Registrar SHA probado y resultados nuevos. La revisión, aprobación ajena y tres
   checks verdes siguen siendo condiciones posteriores; no publicar ni fusionar
   automáticamente por completar estos pasos.
