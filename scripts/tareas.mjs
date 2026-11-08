import { spawnSync } from 'node:child_process'
import { mkdirSync, writeFileSync, rmSync, readdirSync } from 'node:fs'
import { resolve } from 'node:path'

const raiz = process.cwd()
const api = resolve('apps/api')
const pnpm = process.platform === 'win32' ? 'pnpm.cmd' : 'pnpm'
function ejecutar(cmd, args, cwd = raiz, capturar = false) {
  const r = spawnSync(cmd, args, {
    cwd,
    stdio: capturar ? 'pipe' : 'inherit',
    encoding: 'utf8',
    shell: process.platform === 'win32' && cmd.endsWith('.cmd'),
  })
  if (r.error) throw r.error
  if (r.status !== 0) {
    if (capturar) console.error(r.stderr)
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
  case 'codegen':
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
    process.env.GOOSE_DRIVER = 'postgres'
    process.env.GOOSE_DBSTRING = process.env.TEST_DATABASE_URL
    ejecutar(
      'go',
      [
        'run',
        'github.com/pressly/goose/v3/cmd/goose@v3.24.1',
        '-dir',
        'db/migrations',
        'up',
      ],
      api,
    )

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
    process.env.WEB_E2E_BINARY = web
    process.env.API_E2E_BINARY = binario
    process.env.GOCOVERDIR = cobertura
    ejecutar('go', ['test', '-tags=e2e', '-count=1', '-v', './e2e'], api)
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
