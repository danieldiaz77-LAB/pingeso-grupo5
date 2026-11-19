# Prompt para continuar PINGESO con Claude

Copia el bloque siguiente en una conversación nueva. Dale acceso al proyecto local; si no puede leer el disco, adjunta continuidad-claude.md y acuerda cómo compartir los cambios sin publicar ni exponer secretos.

---

Continúa el proyecto PINGESO Grupo 5 desde la integración local existente. Quiero trabajo concreto hasta dejar una base integrada y comprobable; no otro diagnóstico ni reconstruir desde cero.

Repositorio: https://github.com/danieldiaz77-LAB/pingeso-grupo5
PR #1: https://github.com/danieldiaz77-LAB/pingeso-grupo5/pull/1 (mantener en borrador).

Raíz original:
C:\Users\danie\OneDrive\Desktop\ULTIMOS SEMESTRES\Semestre 2_2026\Proyecto Ingeniería de Software (Pingeso)\Proyecto_PINGESO\pingeso-grupo5

Worktree activo: .tmp/worktree-integracion-pipeline-plantilla, dentro de esa raíz.
Rama: integration/pipeline-plantilla-local.

Lee primero el documento canónico completo:
.tmp/worktree-integracion-pipeline-plantilla/docs/integracion/continuidad-claude.md
Después lee AGENTS.md aplicables y docs/referencias/DAS — Desarrollo de Microservicios.md dentro del worktree. Contrasta cualquier adjunto más antiguo con disco y con el documento: no deshagas avances posteriores.

El merge ya está iniciado: HEAD c445c03450692a0c7dd302b07edbbdf72ab20395 y MERGE_HEAD 7d2b7dcc216eec847b78b7d7aeb4731aec1cd94d al traspaso. Hay cambios locales y posiblemente conflictos de índice; inspecciónalos. NO vuelvas a ejecutar preparar-integracion-local.ps1 ni a aplicar el parche de preparación. Conserva historia de Daniel y Benja y todo cambio no confirmado.

Último resultado manual: typecheck aprobado y frontend 20/20 tests en 6 archivos; cobertura 84,93% líneas, 85,89% sentencias, 71,95% ramas y 73,91% funciones. No equivale al 80% de cobertura NUEVA Sonar. Ya pasaron codegen/frozen install, comprobaciones Go y build en estados documentados; no hay verify final completo, e2e PostgreSQL/Compose integrado ni seguridad actual acreditados.

Pendientes prioritarios: integrar auth/config/middleware de Benja y pruebas sin perder salud pública ni controles, resolver dependencias mínimas sin actualizaciones no aprobadas, revisar índice/conflictos por archivo, comprobar reproducibilidad de generación, e2e binario real con PostgreSQL propio, Playwright, Compose aislado y escaneos API/frontend/Goose por separado. Docker Desktop Linux está disponible según última salida. Seguridad e imagen históricas no están verdes.

Las herramientas exactas se cargan desde el worktree con:
. ..\..\scripts\entorno-local.ps1
Node 22.22.0, pnpm 10.13.1, Go 1.26.5. Vitest necesitó almacén virtual corto por fallo reproducido con rutas largas: C:\Users\danie\AppData\Local\pnpm-pingeso-4086dfaa. Verifica node_modules/.modules.yaml antes de reinstalar; no cambies versiones ni vuelvas a la ruta larga inadvertidamente.

El contexto terminal de Codex falló con CreateProcessAsUserW 1260, pero mi PowerShell manual funciona. Si tu sesión funciona, ejecuta las operaciones autorizadas; si reproduce el bloqueo, no lo eludas ni repitas intentos equivalentes. Prepara cambios por archivos y dame comandos PowerShell 5.1 por etapas cortas, indicando directorio y esperando mis salidas. No afirmes haber ejecutado nada sin evidencia.

No publiques, no hagas push, commit automático, merge en main ni auto-merge. No cambies protecciones, Quality Gate, umbrales ni excepciones. Mantén cobertura nueva mínima 80%. Actualizaciones aprobadas requieren PR separado según DAS. No recursos Docker ajenos, prune, BD compartida ni truncados fuera del recurso propio. No pedir ni imprimir secretos.

Daniel se ocupa del pipeline, Benja de la plantilla y Manuel del contrato funcional. No desarrolles la parte de Manuel, ni inventes servicios/endpoints a partir de Ideas (2).pdf. Conserva solo mecanismo Go/Huma cmd/openapi → contracts/openapi.yaml → schema.d.ts y documenta cómo incorporar su entrega cuando publique rama/PR/SHA.

Primero comprueba estado actual y registra el plan concreto de los pendientes, luego avanza con las tareas independientes de aprobaciones. Mantén actualizado docs/integracion/continuidad-claude.md en cada avance: fecha, SHA/estado probado, cambios, comandos/resultados, decisiones/bloqueos y siguiente paso. Al finalizar, entrega tabla aprobado/fallido/no ejecutado y pendientes por responsable. No cierres Hito 0 mientras seguridad e imagen sigan pendientes.
