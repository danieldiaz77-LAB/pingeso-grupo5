> **Actualización 19/11/2026:** este documento conserva información histórica de preparación. El worktree y el merge ya existen; contrato/lock fueron regenerados y hay validaciones posteriores. Para el estado vigente y próximos pasos, leer [continuidad-claude.md](continuidad-claude.md). No volver a ejecutar el script de preparación ni interpretar las tablas antiguas como estado actual.

# Continuación y desbloqueo de la validación

Fecha: 17 de noviembre de 2026, Chile. Esta integración NO está terminada.

## Bloqueo reproducido

Comando solicitado: `git fetch origin`.

Ejecutable empleado por la herramienta:
`C:\WINDOWS\System32\WindowsPowerShell\v1.0\powershell.exe`

Invocación (sin datos sensibles):

```powershell
powershell.exe -NoProfile -Command "try { [Console]::OutputEncoding=[System.Text.Encoding]::UTF8 } catch {}
git fetch origin"
```

Mensaje exacto relevante:
`CreateProcessAsUserW failed: 1260 (La directiva de grupo bloquea a este programa. Para obtener más información, póngase en contacto con el administrador del sistema.)`

La herramienta termina con `Windows error 1260` antes de crear el proceso.
No es un fallo de GitHub, de autenticación Git, del merge o de los tests.
`git status --short --branch` y una lectura de archivos sí funcionaron antes.
No se puede deducir de ello qué regla concreta provoca el bloqueo.

Hubo además un fallo independiente de la revisión automática del entorno:
`Automatic approval review failed: You've hit your usage limit`.
La operación de archivos de esa llamada NO se ejecutó. En la continuación el acceso
autorizado a archivos volvió a funcionar y se guardó el respaldo antes de editar.
Ese límite no explica el error 1260 de Windows. No se intentó otra shell, elevar
privilegios ni ejecutar procesos a través de Node para sortear la directiva.

## Paso concreto del usuario/administrador

1. Entregar al administrador el ejecutable, comando y error anteriores.
2. Pedir que identifique la regla de control de aplicaciones/directiva que impide
   crear PowerShell **en el contexto de ejecución del agente**, y que habilite ese
   uso y las herramientas aprobadas Git, Go, Node/pnpm y Docker. No desactivar
   controles globales ni usar ExecutionPolicy Bypass como sustituto.
3. Después del ajuste, reintentar `git fetch origin` desde esta misma sesión del
   agente. Que funcione en una terminal interactiva distinta no acredita por sí solo
   que el agente esté habilitado.
4. Solo entonces materializar con `.tmp/preparar-integracion-local.ps1`.
   El script comprueba hashes/commits y exige un destino nuevo.
   Si deja un worktree con conflictos, continuar en él: NO volver a crear otro ni
   borrar el existente. Revisar cada archivo y usar `git add -- ruta` únicamente
   después de comprobarlo. No usar `git add --all` para resolver conflictos.
5. Seguir validacion.md: resolver dependencias aprobadas, regenerar lock/contrato,
   comprobar instalación congelada/reproducibilidad, verify, e2e y Compose aislado.
   No publicar esta integración. El contrato funcional de Manuel permanece pendiente.

## Correcciones efectivamente escritas en esta continuación

- Se verificaron los 139 hashes del paquete previo y se guardó un respaldo completo.
- El script valida rutas, manifiesto y commits; rechaza rutas fuera del destino,
  enlaces/reparse points y conflictos no previstos.
- Se eliminó `git add --all`. El índice conserva los conflictos para revisión
  explícita; el script informa y termina con error si quedan entradas sin resolver.
- No hay commit, push, merge en main ni limpieza Docker en ese script.
- Los tests HTTP ya no clonan `public.ejemplos`.
  El runner prepara Goose v3.24.1 sin cambiar el grafo API y expone GOOSE_E2E_BINARY.
  Cada fixture crea un esquema aleatorio y construye una URL con search_path y timezone.
  Goose usa esa URL mediante variables de entorno; su tabla de versiones queda en
  el mismo esquema. Todas las conexiones del pool usan esos parámetros.
- Seed y handlers usan ese pool. Se comprueban current_schema y la presencia de las
  tablas ejemplos/goose_db_version dentro del esquema.
- La conexión administrativa solo crea y elimina el esquema propio, con identificador
  escapado. La limpieza cierra servidor y pool antes de DROP SCHEMA del fixture,
  luego cierra la conexión administrativa; no hay TRUNCATE ni borrado de otras tablas.
- Las migraciones actuales usan nombres no cualificados: se revisó su contenido.
  Las migraciones del proceso real siguen aplicándose en public, pero exclusivamente
  dentro de la base efímera de la ejecución, distinta del aislamiento por fixture.
- Se añadió TestAislamientoDeEsquemas y se exige su ejecución en el runner.
  Estas garantías de código aún requieren validación real con PostgreSQL.

## Validaciones de esta continuación

| Validación                                              | Estado                  | Alcance                                               |
| ------------------------------------------------------- | ----------------------- | ----------------------------------------------------- |
| git status --short --branch                             | Aprobado                | Rama original limpia al inicio                        |
| Hashes de la preparación anterior                       | Aprobado                | 139 archivos antes de editar                          |
| git fetch origin                                        | Fallido                 | Windows error 1260; no inició                         |
| Revisión del script                                     | Aprobado                | Revisión estática de cambios; no ejecución PowerShell |
| Formato/sintaxis de scripts JS/documentos modificados   | Ver informe de revisión | API de Prettier standalone; no CLI                    |
| Worktree, merge y resolución real de conflictos         | No ejecutado            | Terminal bloqueada                                    |
| Auth/config/middleware de Benja                         | No ejecutado            | Incorporación pendiente; no se adoptó su grafo        |
| OpenAPI/schema.d.ts/lockfile reproducibles              | No ejecutado            | Generados y lock todavía pendientes                   |
| gofmt, lint, tipos, compilación, tests backend/frontend | No ejecutado            | No afirmar éxito                                      |
| E2E proceso real, aislamiento PostgreSQL y navegador    | No ejecutado            | Código preparado                                      |
| Compose, migraciones, healthz/readyz, frontend/API      | No ejecutado            | Sin contenedores nuevos                               |
| govulncheck, audit, Trivy, Sonar/cobertura              | No ejecutado            | Sin resultados nuevos atribuibles a integración       |

Incorporado en archivos: ejemplos genéricos backend/frontend, migración/queries,
service/store/repository, mocks y adaptación de pruebas, además de los controles
del pipeline y el mecanismo único de contrato.
Pendiente: auth/config/middleware y sus pruebas, resolución de dependencias,
generados/lock, materialización Git y validación completa.
No se desarrolló el contrato funcional de Manuel ni se cambió el PR #1.
