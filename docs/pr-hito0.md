Se incorpora el pipeline de integración continua del grupo 5 con los checks `calidad`,
`seguridad` e `imagen`: contrato OpenAPI, pruebas, cobertura, SonarQube y controles de
dependencias, secretos y contenedores.

La base técnica mínima en Go y React incluye PostgreSQL, migraciones y endpoints de salud.
No contiene funcionalidades de negocio. Benja debe revisar la base y Manuel completar el
contrato funcional. El servidor SPA carga archivos regulares mediante os.Root al arrancar
y sirve un catálogo inmutable; rechaza enlaces simbólicos y traversal. El bundle requiere
reinicio para reflejar cambios. Se mantienen assets, fallback SPA, HEAD y proxy API.

Validación local:

- Migración y nueve pruebas e2e con PostgreSQL real aprobadas, incluidos errores de arranque,
  healthchecks, puerto ocupado, contrato generado, bundles inválidos y comportamiento HTTP.
- La reproducción anterior expuso un fixture externo mediante un enlace simbólico; el
  binario corregido rechazó el bundle sin exponerlo.
- Cobertura local Go: 171/180 sentencias (95%), con rutas/formato y lector oficial verificados.
- `pnpm verify` y Gitleaks aprobados; actionlint, Hadolint, Dockle y pruebas negativas según
  las evidencias registradas. Las imágenes API, web y migraciones se construyeron.
- Compose completo aprobado: PostgreSQL saludable, migraciones terminadas, backend y
  frontend saludables; healthz, readyz, HTML y proxy API comprobados. Puertos provisionales.
- Los escaneos locales anteriores de imágenes mantienen hallazgos Trivy; no se actualizaron
  versiones pendientes de aprobación ni se desactivaron controles.

[Resultado remoto de 8fad783](https://github.com/danieldiaz77-LAB/pingeso-grupo5/actions/runs/42739207546):

- `calidad`: **aprobado**; verify, contrato, mutación, e2e PostgreSQL/navegador y Sonar pasaron.
  Cobertura nueva **95,1%**, Quality Gate **OK**, ratings A, duplicación 0%; cero incidencias
  abiertas. S6549 dejó de reportarse sin exclusiones ni resolución manual.
- `seguridad`: **fallido** por govulncheck (19 avisos: 17 stdlib y 2 x/text), audit (8 avisos)
  y Trivy filesystem (2 HIGH). Gitleaks y Hadolint pasaron. GO-2026-6604 se refiere a creación
  de directorios mediante junctions en Windows; se documentó el alcance del aviso adicional.
- `imagen`: **omitido**, porque requiere los dos checks anteriores exitosos. Construcciones
  y escaneos locales son evidencia independiente, no aprobación del job en CI.

Pendientes:

- Aprobar las actualizaciones documentadas, aplicarlas y repetir los controles de seguridad.
  La propuesta Go 1.26.10 ya incluye el fix del aviso adicional; no se aplicó.
- Revisar con el administrador la duplicación: el gate permite exactamente 3% y el DAS exige
  menos de 3%. El 80% de cobertura es más estricto que el DAS y se mantiene sin cambios.
- Comprobar el ajuste efectivo de cambios pequeños. Automatic Analysis está desactivado;
  el PR compara contra main y la API informa previous_version heredado para código nuevo.
  El endpoint específico de New Code no está disponible; no acredita otra configuración.
- Confirmar nombre y número del servicio y la asignación oficial de puertos.
- Obtener revisión de Benja sobre la base técnica, completar el contrato con Manuel y
  obtener la aprobación de otro integrante.

El PR permanece en borrador hasta resolver los bloqueos y comprobar los tres checks.
**El Hito 0 no está cerrado.** La cobertura local no acredita por sí sola el Quality Gate.
No se solicita revisión formal todavía.

Evidencias: [validación](https://github.com/danieldiaz77-LAB/pingeso-grupo5/blob/feat/pipeline-hito0/docs/evidencias/validacion-local.md),
[actualizaciones propuestas](https://github.com/danieldiaz77-LAB/pingeso-grupo5/blob/feat/pipeline-hito0/docs/actualizaciones-pendientes.md),
[Sonar/DAS](https://github.com/danieldiaz77-LAB/pingeso-grupo5/blob/feat/pipeline-hito0/docs/sonar-das.md)
y [Sonar del PR](https://sonarcloud.io/dashboard?id=danieldiaz77-LAB_pingeso-grupo5&pullRequest=1).
