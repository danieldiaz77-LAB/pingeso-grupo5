import { readdirSync, mkdtempSync, readFileSync, writeFileSync, mkdirSync } from 'node:fs'
import { resolve, join } from 'node:path'
import { spawnSync } from 'node:child_process'
mkdirSync('.tmp', { recursive: true })
const copia = mkdtempSync(resolve('.tmp/mutacion-'))
function copiar(origen, destino) {
  mkdirSync(destino, { recursive: true })
  for (const entrada of readdirSync(origen, { withFileTypes: true })) {
    const desde = join(origen, entrada.name)
    const hasta = join(destino, entrada.name)
    if (entrada.isDirectory()) copiar(desde, hasta)
    else if (entrada.isFile()) writeFileSync(hasta, readFileSync(desde))
  }
}
copiar('apps/api', copia)
const archivo = join(copia, 'internal/api/salud.go')
const original = readFileSync(archivo, 'utf8')
const proteccion = /if err := base\.Ping\(ctx\); err != nil \{/
if (!proteccion.test(original)) throw new Error('No se encontró la protección a mutar')
writeFileSync(archivo, original.replace(proteccion, 'if base.Ping(ctx); false {'))
const resultado = spawnSync(
  'go',
  ['test', '-count=1', './e2e', '-run', '^TestReadiness$'],
  { cwd: copia, encoding: 'utf8' },
)
if (resultado.error) throw resultado.error
if (
  resultado.status !== 1 ||
  !resultado.stdout.includes('/readyz: esperado 503, recibido 200')
) {
  console.error(resultado.stdout, resultado.stderr)
  throw new Error(
    'La mutación no produjo el rechazo esperado; un error de compilación no cuenta',
  )
}
console.log(
  'Protección de readiness eliminada solo en copia temporal: prueba falla con 200 en vez de 503, como se esperaba',
)
