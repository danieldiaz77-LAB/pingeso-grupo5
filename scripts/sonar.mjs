import { readFileSync, statSync } from 'node:fs'
for (const clave of ['SONAR_TOKEN', 'SONAR_HOST_URL']) {
  if (!process.env[clave]?.trim())
    throw new Error(
      `Bloqueo Sonar: falta el secret ${clave}; el Quality Gate es obligatorio`,
    )
}
const config = readFileSync('sonar-project.properties', 'utf8')
for (const propiedad of [
  'sonar.projectKey=danieldiaz77-LAB_pingeso-grupo5',
  'sonar.organization=danieldiaz77-lab',
]) {
  if (!config.includes(propiedad))
    throw new Error(`Configuración Sonar inesperada: ${propiedad}`)
}
for (const ruta of ['coverage.out', 'apps/web/coverage/lcov.info']) {
  if (statSync(ruta).size === 0) throw new Error(`Cobertura vacía: ${ruta}`)
}
console.log(
  'Configuración y reportes Sonar presentes; el servidor decidirá el Quality Gate',
)
