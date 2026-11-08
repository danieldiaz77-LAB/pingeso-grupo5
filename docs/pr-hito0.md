Se incorpora el pipeline de integración continua del grupo 5 con los checks `calidad`,
`seguridad` e `imagen`. Incluye validación del contrato OpenAPI, pruebas, cobertura, análisis
con SonarQube y controles de dependencias, secretos y contenedores.

Para ejecutar los controles se agregó una base técnica mínima en Go y React, con PostgreSQL,
migraciones y endpoints de salud. No incluye funcionalidades de negocio. La base queda
pendiente de revisión por Benja y el contrato funcional, de completar por Manuel.

Validación local:

- Migración y cuatro pruebas con PostgreSQL real aprobadas. Se corrigió el enlace `$schema`
  que Huma añadía en producción hacia una ruta deshabilitada.
- Cobertura Go: 104 de 141 sentencias (73,76%). Formato, rutas y lector oficial Go comprobados.
- `pnpm verify`, actionlint, Hadolint y Gitleaks aprobados según las evidencias registradas.
- Imágenes de API, web y migraciones construidas.
- Dockle aprobado al umbral `warn`, con INFO y SKIP nativos documentados.
- Las pruebas negativas comprobaron los controles documentados de contrato, secretos,
  arquitectura frontend y mutación de readiness.
- Trivy bloquea las imágenes con 14/12/41 avisos HIGH respectivamente. Los 18 avisos
  alcanzables de govulncheck corresponden a 16 de stdlib y 2 de x/text, no a 18 dependencias.
  `pnpm audit` bloquea por Vite. No se modificaron las versiones pendientes de aprobación.

Pendientes:

- Aprobar y aplicar las actualizaciones de seguridad documentadas y repetir los controles.
- Resolver los hallazgos de Trivy y govulncheck.
- Resolver el hallazgo S6549 del servidor SPA y elevar la cobertura nueva medida por Sonar
  (72,5%) hasta aprobar el Quality Gate. Se conserva 80% de
  cobertura; duplicación <=3% permite el límite que el DAS <3% prohíbe. No se cambió el gate.
  New Code, Automatic Analysis y cambios pequeños requieren comprobación separada.
- Validar el arranque completo con `docker compose up`.
- Confirmar nombre y número del servicio.
- Obtener la revisión de la base técnica y la aprobación de otro integrante.

El PR permanece en borrador hasta resolver los bloqueos y comprobar los tres checks.
**El Hito 0 no está cerrado.** El porcentaje de cobertura local no acredita por sí solo el
cumplimiento del Quality Gate. No se solicita revisión formal todavía.

Evidencias: [validación local](https://github.com/danieldiaz77-LAB/pingeso-grupo5/blob/feat/pipeline-hito0/docs/evidencias/validacion-local.md),
[actualizaciones propuestas](https://github.com/danieldiaz77-LAB/pingeso-grupo5/blob/feat/pipeline-hito0/docs/actualizaciones-pendientes.md)
y [comparación Sonar/DAS](https://github.com/danieldiaz77-LAB/pingeso-grupo5/blob/feat/pipeline-hito0/docs/sonar-das.md).

Resultado de la primera ejecución remota (commit `0f5d93a0eafdea8e2281e53e4abfd69488db9a99`):

- `calidad`: **fallido** en Quality Gate. Verify, contrato, mutación, e2e PostgreSQL,
  navegador y scanner Sonar pasaron. Sonar leyó cobertura Go y LCOV; fallaron cobertura
  nueva 72,5% <80% y security rating C >A, con S6549 en `apps/api/internal/spa/servidor.go:43`.
- `seguridad`: **fallido** por govulncheck, audit y Trivy. Gitleaks y Hadolint pasaron.
- `imagen`: **omitido**, porque requiere calidad y seguridad exitosos; las construcciones
  y escaneos locales anteriores se mantienen como evidencia independiente.

[Ejecución comprobada](https://github.com/danieldiaz77-LAB/pingeso-grupo5/actions/runs/42489503785) · [Sonar del PR](https://sonarcloud.io/dashboard?id=danieldiaz77-LAB_pingeso-grupo5&pullRequest=1).
