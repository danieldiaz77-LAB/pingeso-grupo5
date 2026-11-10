# Quality Gate actual frente al DAS

Consulta de solo lectura realizada el 8 de noviembre de 2026. Proyecto
`danieldiaz77-LAB_pingeso-grupo5`, organización `danieldiaz77-lab`. Está asignado el gate
integrado **Sonar way**, id 9. Evidencia: `evidencias/comparacion-externa.json`.
No se modificó el gate, el método de análisis ni la definición de código nuevo.

Las condiciones siguientes se aplican a **código nuevo**, no a todo el historial.

| Control                 | DAS                            | Condición de fallo actual (API)      | Qué admite el gate            | Comparación                                                             |
| ----------------------- | ------------------------------ | ------------------------------------ | ----------------------------- | ----------------------------------------------------------------------- |
| Bugs nuevos             | 0                              | new_reliability_rating > 1           | Rating A = 0 bugs             | Equivalente para esta métrica                                           |
| Vulnerabilidades nuevas | 0                              | new_security_rating > 1              | Rating A = 0 vulnerabilidades | Equivalente para esta métrica                                           |
| Hotspots revisados      | 100%                           | new_security_hotspots_reviewed < 100 | 100%                          | Equivalente                                                             |
| Cobertura nueva         | >=60%                          | new_coverage < 80                    | >=80%                         | Más estricto: rechaza 60% a menos de 80%, aunque el DAS lo admita       |
| Duplicación nueva       | <3%                            | new_duplicated_lines_density > 3     | <=3%                          | Más permisivo en el límite: acepta exactamente 3%, prohibido por el DAS |
| Mantenibilidad          | Sin umbral numérico específico | new_maintainability_rating > 1       | Rating A de deuda técnica     | Condición adicional; no sustituye las anteriores                        |

La equivalencia de ratings se basa en las claves efectivamente devueltas por la API y en las
[definiciones oficiales de SonarQube Cloud](https://docs.sonarsource.com/sonarqube-cloud/managing-your-projects/metric-definitions):
security_rating A significa cero vulnerabilidades y reliability_rating A significa cero
bugs. No se confunden esas claves con otras métricas de severidad de software quality.

## Decisión propuesta, sin aplicar

Conservar el **80%**. No hay motivo para rebajarlo a 60%: cumplir 80% también satisface el
mínimo del DAS. El informe de Go mide sentencias cubiertas, mientras Sonar usa su propia
medida de líneas/condiciones y el alcance de código nuevo; no deben igualarse los porcentajes.

Corregir por separado duplicación para que **3% falle**. Antes de editar o copiar el gate,
el administrador debe comprobar los operadores y precisión que admite la interfaz; si solo
permite fallo por `>` será necesario acordar un umbral estrictamente menor que 3 con esa
precisión. No se cambia automáticamente a otro número ni se modifica el gate integrado.
La tabla anterior es la diferencia concreta que debe aprobarse antes de escribir en Sonar.

## Código nuevo y cambios pequeños: pendientes distintos

- La consulta `api/new_code_periods/show` devolvió `Unknown url`. Ese endpoint no está
  disponible en esta instancia de Cloud; **no demuestra que falte configuración**.
- La consulta de `sonar.qualitygate.ignoreSmallChanges` devolvió `settings: []`. No acredita
  el valor efectivo ni que esté deshabilitada la excepción para cambios pequeños.
- El responsable del proyecto debe revisar y registrar en la interfaz la definición de New Code de main
  (versión anterior, días u otra opción admitida) y la base de comparación de los PR.
  Debe ser la definición acordada antes del primer análisis; no cambiarla para ocultar hallazgos.
- Comprobar que cobertura/duplicación se evalúen también en cambios pequeños según el DAS.
- Confirmar **Administration → Analysis Method → Automatic Analysis: Off**. Cobertura Go
  requiere análisis por CI. El workflow ya prepara reportes y ejecuta scan + Quality Gate.
- Un proyecto público accesible y un gate configurado no son evidencia de análisis aprobado.
  Sin token local no se ejecuta un análisis autenticado. La disponibilidad de los secrets GitHub se comprobará mediante CI; no se leyó su contenido.

Fuentes: [código nuevo](https://docs.sonarsource.com/sonarqube-cloud/standards/about-new-code),
[cobertura Go y análisis por CI](https://docs.sonarsource.com/sonarqube-cloud/analyzing-source-code/test-coverage/go-test-coverage).

## Comprobación del 9 de noviembre de 2026

Consultas de solo lectura: `settings/values` informa `sonar.autoscan.enabled=false`:
el proyecto utiliza análisis por CI y tiene Automatic Analysis desactivado. La misma API
informa `sonar.leak.period=previous_version`, heredado de la instancia; `project_pull_requests/list`
confirma que el PR 1 tiene `main` como base/target. El endpoint específico `new_code_periods/show`
sigue sin estar disponible (404), por lo que no se atribuye a esa API una definición adicional.

No se obtiene valor efectivo de `sonar.qualitygate.ignoreSmallChanges`. El resultado del gate
anterior indicó `ignoredConditions=false`, que acredita la evaluación de ese PR, no la
configuración global para futuros cambios pequeños. Este último ajuste requiere revisión
administrativa. El gate asignado sigue siendo Sonar way; no se editó ninguna condición.
