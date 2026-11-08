import { spawnSync } from 'node:child_process'
import { mkdirSync, writeFileSync, statSync } from 'node:fs'
function ejecutar(cmd, args) {
  const r = spawnSync(cmd, args, { encoding: 'utf8' })
  if (r.error) throw r.error
  if (r.status !== 0) throw new Error(`${cmd} falló (${r.status}): ${r.stderr}`)
  return r.stdout
}
if (process.argv[2] === 'generated') {
  for (const archivo of ['contracts/openapi.yaml', 'apps/web/src/api/schema.d.ts']) {
    if (!statSync(archivo).size) throw new Error(`Generado vacío: ${archivo}`)
  }
  const rutas = ['contracts', 'apps/web/src/api', 'apps/api/internal/repository']
  ejecutar('git', ['diff', '--exit-code', 'HEAD', '--', ...rutas])
  const nuevos = ejecutar('git', [
    'ls-files',
    '--others',
    '--exclude-standard',
    '--',
    ...rutas,
  ])
  if (nuevos.trim())
    throw new Error(
      `Generados sin seguimiento; revísalos e inclúyelos en el commit:\n${nuevos}`,
    )
} else if (process.argv[2] === 'compatible') {
  const base = process.env.CONTRACT_BASE_SHA
  if (!base || !/^[a-f0-9]{40}$/.test(base) || /^0+$/.test(base))
    throw new Error('Falta un SHA base válido para comparar contrato')
  // Un fallo al resolver el commit nunca significa contrato ausente.
  ejecutar('git', ['cat-file', '-e', `${base}^{commit}`])
  const ruta = ejecutar('git', [
    'ls-tree',
    '--name-only',
    base,
    '--',
    'contracts/openapi.yaml',
  ]).trim()
  if (!ruta) {
    console.log(
      'Primera línea base: no existe contracts/openapi.yaml en el commit base. El contrato actual se genera y valida; todavía no hay versión anterior para oasdiff.',
    )
  } else {
    mkdirSync('.tmp', { recursive: true })
    writeFileSync(
      '.tmp/openapi-base.yaml',
      ejecutar('git', ['show', `${base}:contracts/openapi.yaml`]),
    )
    const version = process.env.OASDIFF_VERSION
    if (version !== 'v1.35.0') throw new Error('OASDIFF_VERSION debe ser v1.35.0')
    console.log(
      ejecutar('go', [
        'run',
        `github.com/oasdiff/oasdiff@${version}`,
        'breaking',
        '.tmp/openapi-base.yaml',
        'contracts/openapi.yaml',
        '--fail-on',
        'ERR',
      ]),
    )
  }
} else throw new Error('Tarea de contrato desconocida')
