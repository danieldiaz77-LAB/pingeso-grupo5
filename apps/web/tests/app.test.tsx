import { expect, test } from 'vitest'
import { renderToStaticMarkup } from 'react-dom/server'
import { App } from '../src/app'
test('presenta el alcance técnico sin atribuir funcionalidades de negocio', () => {
  const html = renderToStaticMarkup(<App />)
  expect(html).toContain('<h1>Entorno técnico PINGESO</h1>')
  expect(html).toContain('Contrato funcional pendiente')
})
