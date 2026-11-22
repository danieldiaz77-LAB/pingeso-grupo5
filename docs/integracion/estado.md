> **Actualización 19/11/2026:** este documento conserva información histórica de preparación. El worktree y el merge ya existen; contrato/lock fueron regenerados y hay validaciones posteriores. Para el estado vigente y próximos pasos, leer [continuidad-claude.md](continuidad-claude.md). No volver a ejecutar el script de preparación ni interpretar las tablas antiguas como estado actual.

**Actualización:** consultar [continuación y bloqueo actual](continuacion.md). El script ya no añade conflictos al índice automáticamente.

# Estado de la preparación local

**No es una integración Git terminada ni una base validada.** Se preparó código real
en una copia independiente. No se crearon rama/worktree, no se hizo fetch ni merge y no
se ejecutaron pruebas de aplicación, Docker o Sonar después del bloqueo de Windows.

La primera llamada `git status --short --branch` sí mostró el árbol limpio en
`feat/pipeline-hito0`. La siguiente llamada fue bloqueada antes de ejecutar comandos,
con error Windows 1260. Se continuó solo con lectura/escritura de archivos y APIs de GitHub.

Fuentes verificadas:

- Pipeline: `c445c03450692a0c7dd302b07edbbdf72ab20395`.
- Plantilla de Benja: `7d2b7dcc216eec847b78b7d7aeb4731aec1cd94d`.
- Rama prevista (NO creada): `integration/pipeline-plantilla-local`.
- Copia preparada: `.tmp/integracion-pipeline-plantilla-2026-11-17`.
- Fuentes originales de Benja: carpeta hermana con sufijo `-fuentes`.
- Script futuro: `.tmp/preparar-integracion-local.ps1`.
- Parche respecto al pipeline y manifiesto SHA-256: archivos hermanos de esta copia.
- No hay commit de integración; el manifiesto identifica el estado preparado.

## Resoluciones implementadas en archivos

| Solapamiento      | Resolución preparada                                                                                                             |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| package.json raíz | Conserva verify con formato, lint, arquitectura, tipos, Go, tests y build; añade dev y sqlc                                      |
| API y arranque    | Mantiene arranque, errores, límites HTTP y registro del pipeline; conecta service/store/repository de ejemplos                   |
| Salud             | Mantiene respuestas y readiness del pipeline; registra ejemplos en la misma composición Huma usada por generador                 |
| Migraciones       | Conserva 00001_base_tecnica.sql y añade 202611080001_crear_ejemplos.sql, queries, seed y sqlc                                    |
| Generador         | Solo cmd/openapi; salida contracts/openapi.yaml y schema.d.ts; sin contrato paralelo manual                                      |
| Frontend          | Solo app.tsx en minúsculas; formulario/listado/hooks de Benja; URL base por props y cliente *.api.ts                             |
| Mocks             | MSW solo en desarrollo; si falla su arranque no se cambia silenciosamente a la API real                                          |
| Tests frontend    | Conserva tests/e2e del pipeline, traslada tests de Benja a tests y añade recorrido Playwright con MSW                            |
| Tests API         | Conserva todos los e2e del pipeline, adapta HTTP de Benja a esquemas efímeros y añade e2e instrumentado de crear/listar/rechazos |
| Sonar             | Clasifica *_test.go como pruebas; mantiene exclusión existente de repository generado; no excluye nuevo código productivo        |
| Docker y SPA      | Dockerfiles, Compose base, hooks, servidor SPA y pruebas previas se conservan byte por byte                                      |
| CI                | Conserva checks y pasos; BD efímera identificada; runner rechaza skip y ausencia de escenarios requeridos                        |

Estas son resoluciones de contenido; no se ejecutó un merge para obtener una lista
de conflictos Git. El script materializa el merge y estas resoluciones cuando se habilite
la terminal. Exige ramas remotas en los SHA indicados y un destino nuevo; no hace commit/push.

## Bloqueos explícitos

1. **Contrato y tipos aún sin regenerar:** los archivos presentes son de la base de salud.
   Los consumidores nuevos requieren los tipos de ejemplos; no se afirma typecheck aprobado.
   Ejecutar pnpm codegen y versionar los resultados antes de validar CI.
2. **Lockfile pendiente:** package.json incluye dependencias del catalog ya fijado, pero
   no se ejecutó pnpm para resolver sus importers/peers. Frozen-lockfile todavía debe fallar.
3. **Formato Go y compilación pendientes:** se prepararon fuentes Go pero no se ejecutó gofmt.
4. **Auth/config/middleware de Benja diferidos:** sus originales están conservados fuera de
   la copia activa. Incorporarlos exige revisar keyfunc/jwt y transitivas, probar salud pública
   y unificar la configuración sin romper pruebas de arranque. No se incorporaron ni omitieron
   silenciosamente sus tests de configuración: permanecen pendientes junto a ese bloque.
5. **Ningún porcentaje de cobertura de la rama anterior acredita esta copia.**
   El reporte coverage.out anterior no se copió como resultado de integración.
6. Esta preparación no está lista para producción ni cierra Hito 0/Hito 1.

No se adoptaron aumentos de versiones ni todo el go.mod de Benja. Las capas de ejemplos
y los tests adaptados usan Huma/pgx ya existentes; los cambios de grafo diferidos están
enumerados en dependencias.json. No se sustituyen e2e reales por los de httptest.

## Responsables

| Responsable                | Pendiente concreto                                                                                                                     |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| Daniel                     | Materializar worktree, revisar diff, resolver lock aprobado, regenerar/versionar contrato, ejecutar runbook y registrar SHA/resultados |
| Benja + Daniel             | Revisar adaptación de ejemplos y acordar incorporación de auth/config; confirmar uso de URL base y mocks                               |
| Manuel                     | Publicar contrato funcional con SHA, requests/responses, errores, permisos y mocks acordados; seguir contrato-manuel.md                |
| Jefe                       | Aprobar versiones en PR separado con tabla DAS; resolver diferencias de grafo antes de adoptarlas                                      |
| Administrador Windows      | Habilitar terminal y herramientas aprobadas; el agente no intentó eludir la directiva                                                  |
| Administrador GitHub/Sonar | Acreditar protección efectiva de main y ajustes pendientes de Sonar, sin cambiar umbrales aquí                                         |

## Remoto (solo lectura)

main devolvió protected=false; rulesets devolvió []; la protección detallada devolvió
403 Resource not accessible by integration. No se puede certificar protección efectiva.
No se cambió el PR #1 ni se solicitó revisión; continúa sin acciones de publicación,
merge o auto-merge por esta preparación.

## Cómo continuar

Una vez habilitada la terminal, revisar el script y ejecutar desde la raíz original:
`& .\.tmp\preparar-integracion-local.ps1`.
El script deja un worktree con merge sin commit. Revisar `git diff --cached` y
seguir [validación](validacion.md), [dependencias](seguridad.md) y
[incorporación de Manuel](contrato-manuel.md). No ejecutar la copia .tmp como si
fuera ya un repositorio independiente: aún hereda el .git del directorio superior.

El parche adjunto es una alternativa para inspección o aplicación sobre una rama nueva
basada en el SHA del pipeline; no conserva por sí solo la ascendencia de Benja.
No aplicar el parche encima del merge resuelto por el script.
