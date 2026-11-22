import ts from 'typescript'
import { readdirSync, readFileSync } from 'node:fs'
import { join, resolve, relative } from 'node:path'
function archivos(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((e) =>
    e.isDirectory() ? archivos(join(dir, e.name)) : [join(dir, e.name)],
  )
}
let errores = 0
for (const ruta of archivos('apps/web/src').filter((p) => /\.[jt]sx?$/.test(p))) {
  const source = ts.createSourceFile(
    ruta,
    readFileSync(ruta, 'utf8'),
    ts.ScriptTarget.Latest,
    true,
  )
  const fallo = (mensaje) => {
    console.error(`${ruta}: ${mensaje}`)
    errores++
  }
  function visitar(nodo) {
    if (ts.isIdentifier(nodo) && ['localStorage', 'sessionStorage'].includes(nodo.text))
      fallo('almacenamiento de sesión prohibido')
    if (
      ts.isStringLiteral(nodo) &&
      ['localStorage', 'sessionStorage'].includes(nodo.text) &&
      ts.isElementAccessExpression(nodo.parent)
    )
      fallo('almacenamiento de sesión prohibido')
    if (ts.isCallExpression(nodo)) {
      const expresion = nodo.expression
      if (
        (ts.isIdentifier(expresion) && expresion.text === 'fetch') ||
        (ts.isPropertyAccessExpression(expresion) && expresion.name.text === 'fetch')
      ) {
        if (!ruta.endsWith('.api.ts')) fallo('fetch solo puede vivir en *.api.ts')
      }
    }
    if (ts.isImportDeclaration(nodo) && ts.isStringLiteral(nodo.moduleSpecifier)) {
      const origen = ruta.replaceAll('\\', '/').match(/\/(componentes|features)\/([^/]+)/)
      const destino = relative(
        process.cwd(),
        resolve(ruta, '..', nodo.moduleSpecifier.text),
      ).replaceAll('\\', '/')
      const otro = destino.match(/\/(componentes|features)\/([^/]+)/)
      if (origen && otro && (origen[1] !== otro[1] || origen[2] !== otro[2]))
        fallo('import entre microcomponentes o features')
    }
    ts.forEachChild(nodo, visitar)
  }
  visitar(source)
}
if (errores) process.exit(1)
