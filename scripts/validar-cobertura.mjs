import { readFileSync, readdirSync, existsSync } from 'node:fs'
import { resolve, relative } from 'node:path'

// Comprueba el reporte real sin alterar rutas, contadores ni porcentajes.
const propiedades = Object.fromEntries(
  readFileSync('sonar-project.properties', 'utf8')
    .split(/\r?\n/)
    .filter((linea) => linea && !linea.startsWith('#'))
    .map((linea) => {
      const indice = linea.indexOf('=')
      return [linea.slice(0, indice), linea.slice(indice + 1)]
    }),
)
const ruta = propiedades['sonar.go.coverage.reportPaths']
const fuentes = propiedades['sonar.sources'].split(',')
const modulo = readFileSync('apps/api/go.mod', 'utf8').match(/^module (.+)$/m)?.[1]
if (!modulo || !ruta) throw new Error('Falta módulo Go o ruta Sonar')
const lineas = readFileSync(ruta, 'utf8').trimEnd().split(/\r?\n/)
if (!/^mode: (set|count|atomic)$/.test(lineas.shift()))
  throw new Error('Formato de cobertura Go inválido')
const archivos = new Set()
let sentencias = 0
let cubiertas = 0
for (const linea of lineas) {
  const dato = linea.match(/^(.+\.go):(\d+)\.(\d+),(\d+)\.(\d+) (\d+) (\d+)$/)
  if (!dato) throw new Error(`Bloque de cobertura inválido: ${linea}`)
  const nombre = dato[1].startsWith(modulo + '/')
    ? 'apps/api/' + dato[1].slice(modulo.length + 1)
    : dato[1]
  const local = relative(process.cwd(), resolve(nombre)).replaceAll('\\', '/')
  if (
    local.startsWith('../') ||
    !fuentes.some((fuente) => local.startsWith(fuente + '/'))
  )
    throw new Error(`Cobertura fuera de fuentes Sonar: ${local}`)
  if (!existsSync(local)) throw new Error(`No se encuentra la fuente: ${local}`)
  const contenido = readFileSync(local, 'utf8').split(/\r?\n/)
  if (
    Number(dato[2]) < 1 ||
    Number(dato[4]) < Number(dato[2]) ||
    Number(dato[4]) > contenido.length
  )
    throw new Error(`Líneas fuera del archivo: ${local}`)
  archivos.add(local)
  sentencias += Number(dato[6])
  if (Number(dato[7]) > 0) cubiertas += Number(dato[6])
}
if (!archivos.size || !sentencias || !cubiertas)
  throw new Error('Reporte vacío o sin ejecución cubierta')
function listar(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((e) =>
    e.isDirectory() ? listar(dir + '/' + e.name) : [dir + '/' + e.name],
  )
}
const sinReporte = fuentes
  .flatMap(listar)
  .filter((f) => f.endsWith('.go') && !archivos.has(f))
console.log(
  JSON.stringify(
    {
      ruta,
      archivos: [...archivos],
      sentencias,
      cubiertas,
      porcentajeSentencias: Number(((100 * cubiertas) / sentencias).toFixed(2)),
      fuentesGoSinReporte: sinReporte,
      alcance:
        'Formato y rutas compatibles; no acredita ingestión remota ni cobertura de código nuevo de Sonar',
    },
    null,
    2,
  ),
)
