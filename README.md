# pingeso-grupo5

Microservicio del grupo 5 de PINGESO, semestre 2-2026. Base reutilizable con Go, React,
PostgreSQL y controles de calidad y seguridad.

**Hito 0 en revisión; no cerrado.** Nombre y número del servicio pendientes. El contrato
actual contiene únicamente salud técnica. Vite 8.1.5 tiene un HIGH confirmado y se conserva
por DAS hasta aprobación del jefe; el check seguridad debe bloquearlo. También hay vulnerabilidades alcanzables en Go y
x/text documentadas en la guía.

- [Guía del pipeline, comandos y pendientes](docs/pipeline-hito0.md)
- [Modos de falla técnicos](docs/modos-de-falla.md)
- [Validaciones locales](docs/evidencias/validacion-local.md)
- [Actualizaciones pendientes de aprobación](docs/actualizaciones-pendientes.md)
- [Comparación Sonar / DAS](docs/sonar-das.md)
- [DAS](docs/referencias/DAS%20—%20Desarrollo%20de%20Microservicios.md)
- [Contrato técnico generado](contracts/openapi.yaml)

Inicio: Go 1.26.5, Node 22.22.0 y pnpm 10.13.1. Ejecutar `pnpm install --frozen-lockfile`,
`pnpm exec lefthook install` y `pnpm verify`. Para Compose, completar `.env.example` como
`.env` y acordar los puertos 51N0/51N1. No asumir N=5.

[PR #1 en borrador](https://github.com/danieldiaz77-LAB/pingeso-grupo5/pull/1) y
[primera ejecución de Actions](https://github.com/danieldiaz77-LAB/pingeso-grupo5/actions/runs/42489503785): calidad y seguridad fallaron; imagen fue omitido.
El [análisis Sonar del PR](https://sonarcloud.io/dashboard?id=danieldiaz77-LAB_pingeso-grupo5&pullRequest=1) leyó los reportes y rechazó el gate por cobertura nueva
72,5% (mínimo 80%) y un hallazgo S6549 en el servidor SPA. El hito sigue abierto.

[Proyecto SonarQube Cloud verificado](https://sonarcloud.io/project/overview?id=danieldiaz77-LAB_pingeso-grupo5).
Su existencia no acredita aún un análisis aprobado.

El servidor SPA carga un catálogo de archivos regulares al iniciar y rechaza bundles con
enlaces simbólicos. Las correcciones incluyen pruebas de errores de arranque y de HTTP.
El arranque completo con Compose fue validado; los puertos utilizados son provisionales
hasta confirmar el número del servicio. Los resultados anteriores de Sonar se conservan
como antecedentes; el siguiente análisis debe comprobar la corrección y la cobertura nueva.
