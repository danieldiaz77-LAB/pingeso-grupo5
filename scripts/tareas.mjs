import { spawnSync } from 'node:child_process'
import { mkdirSync, writeFileSync, rmSync, readdirSync } from 'node:fs'
import { resolve } from 'node:path'

const raiz = process.cwd()
const api = resolve('apps/api')
const pnpm = process.platform === 'win32' ? 'pnpm.cmd' : 'pnpm'
function ejecutar(cmd, args, cwd = raiz, capturar = false, informe) {
  const r = spawnSync(cmd, args, {
    cwd,
    stdio: capturar ? 'pipe' : 'inherit',
    encoding: 'utf8',
    shell: process.platform === 'win32' && cmd.endsWith('.cmd'),
  })
  if (informe) writeFileSync(informe, r.stdout ?? '')
  if (r.error) throw r.error
  if (r.status !== 0) {
    if (capturar) console.error(r.stdout, r.stderr)
    process.exit(r.status ?? 1)
  }
  return r.stdout
}
mkdirSync('.tmp', { recursive: true })
switch (process.argv[2]) {
  case 'verify-go': {
    const cambios = ejecutar('gofmt', ['-l', '.'], api, true)
    if (cambios.trim()) throw new Error(`Ejecuta gofmt: ${cambios}`)
    ejecutar('go', ['mod', 'verify'], api)
    ejecutar('go', ['vet', './...'], api)
    ejecutar('go', ['test', './...'], api)
    break
  }
  case 'test-api':
    ejecutar('go', ['test', './...'], api)
    break
  case 'build-api':
    for (const cmd of ['api', 'spa-server'])
      ejecutar(
        'go',
        [
          'build',
          '-trimpath',
          '-buildvcs=false',
          '-o',
          resolve(`.tmp/${cmd}${process.platform === 'win32' ? '.exe' : ''}`),
          `./cmd/${cmd}`,
        ],
        api,
      )
    break
  case 'sqlc':
    ejecutar('go', ['run', 'github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0', 'generate'], api)
    break
  case 'codegen':
    ejecutar('go', ['run', 'github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0', 'generate'], api)
    mkdirSync('contracts', { recursive: true })
    writeFileSync(
      'contracts/openapi.yaml',
      ejecutar('go', ['run', './cmd/openapi'], api, true),
    )
    ejecutar(pnpm, ['--dir', 'apps/web', 'codegen'])
    break
  case 'e2e': {
    if (process.platform === 'win32')
      throw new Error(
        'El e2e instrumentado requiere Linux/WSL para cerrar con señales y volcar cobertura',
      )
    if (!process.env.TEST_DATABASE_URL)
      throw new Error('TEST_DATABASE_URL obligatorio: PostgreSQL real')
    const nombre = process.env.TEST_DATABASE_NAME
    const url = new URL(process.env.TEST_DATABASE_URL)
    if (!nombre?.startsWith('pingeso_integracion_') || url.pathname !== '/' + nombre)
      throw new Error('Solo se permiten bases efímeras propias pingeso_integracion_*')
    process.env.GOOSE_DRIVER = 'postgres'
    process.env.GOOSE_DBSTRING = process.env.TEST_DATABASE_URL
    const herramientas = resolve('.tmp/herramientas-e2e')
    mkdirSync(herramientas, { recursive: true })
    const gobinAnterior = process.env.GOBIN
    try {
      process.env.GOBIN = herramientas
      ejecutar('go', ['install', 'github.com/pressly/goose/v3/cmd/goose@v3.24.1'], api)
    } finally {
      if (gobinAnterior === undefined) delete process.env.GOBIN
      else process.env.GOBIN = gobinAnterior
    }
    process.env.GOOSE_E2E_BINARY = resolve(herramientas, 'goose')
    ejecutar(process.env.GOOSE_E2E_BINARY, ['-dir', 'db/migrations', 'up'], api)

    const cobertura = resolve('.tmp/cobertura-e2e')
    rmSync(cobertura, { recursive: true, force: true })
    mkdirSync(cobertura)
    const binario = resolve(`.tmp/api-e2e${process.platform === 'win32' ? '.exe' : ''}`)
    ejecutar(
      'go',
      [
        'build',
        '-cover',
        '-coverpkg=./cmd/api,./internal/...',
        '-o',
        binario,
        './cmd/api',
      ],
      api,
    )
    const web = resolve(`.tmp/web-e2e${process.platform === 'win32' ? '.exe' : ''}`)
    ejecutar(
      'go',
      [
        'build',
        '-cover',
        '-coverpkg=./cmd/spa-server,./internal/spa',
        '-o',
        web,
        './cmd/spa-server',
      ],
      api,
    )
    const openapi = resolve('.tmp/openapi-e2e')
    ejecutar('go', ['build', '-cover', '-o', openapi, './cmd/openapi'], api)
    process.env.OPENAPI_E2E_BINARY = openapi
    process.env.WEB_E2E_BINARY = web
    process.env.API_E2E_BINARY = binario
    process.env.GOCOVERDIR = cobertura
    const reporte = ejecutar(
      'go',
      ['test', '-tags=e2e', '-count=1', '-json', './e2e'],
      api,
      true,
      '.tmp/e2e.jsonl',
    )
    writeFileSync('.tmp/e2e.jsonl', reporte)
    const eventos = reporte
      .trim()
      .split(/\r?\n/)
      .filter(Boolean)
      .map((linea) => JSON.parse(linea))
    if (eventos.some((evento) => evento.Action === 'skip'))
      throw new Error('Una prueba requerida se omitió; no se permite aprobar por omisión')
    for (const nombre of [
      'TestSaludReal',
      'TestServidorWebReal',
      'TestContratoDelBinario',
      'TestEjemplosBinarioReal',
      'TestAutenticacionBinarioReal',
      'TestRechazosDeArranque',
      'TestListarEjemplosContraLaBaseDeDatos',
      'TestCrearEjemploYVerificarloEnElListado',
      'TestCrearEjemploRechazaNombreEnBlanco',
      'TestSaludYDisponibilidad',
      'TestAislamientoDeEsquemas',
    ]) {
      if (!eventos.some((evento) => evento.Test === nombre && evento.Action === 'pass'))
        throw new Error('No se acreditó ejecución de ' + nombre)
    }
    if (!readdirSync(cobertura).some((x) => x.startsWith('covcounters.')))
      throw new Error('E2E sin contadores de cobertura')
    ejecutar(
      'go',
      ['tool', 'covdata', 'textfmt', '-i', cobertura, '-o', resolve('coverage.out')],
      api,
    )
    break
  }
  default:
    throw new Error('Tarea desconocida')
}
