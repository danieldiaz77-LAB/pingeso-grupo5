import { renderToStaticMarkup } from 'react-dom/server'
import { expect, test } from 'vitest'
import { ListaEjemplos } from '../src/componentes/ejemplos/lista-ejemplos'

test('explica el estado vacío', () => {
  expect(renderToStaticMarkup(<ListaEjemplos ejemplos={[]} />)).toContain(
    'Todavía no hay ejemplos',
  )
})

test('presenta los datos y la alternativa para una descripción vacía', () => {
  const html = renderToStaticMarkup(
    <ListaEjemplos
      ejemplos={[
        {
          id: 'a',
          nombre: 'Primero',
          descripcion: 'Detalle visible',
          creadoEn: '2026-11-01T12:00:00Z',
        },
        { id: 'b', nombre: 'Segundo', descripcion: '', creadoEn: '2026-11-02T12:00:00Z' },
      ]}
    />,
  )
  expect(html).toContain('Primero')
  expect(html).toContain('Segundo')
  expect(html).toContain('Detalle visible')
  expect(html).toContain('Sin descripción.')
  expect(html.match(/<li /g)).toHaveLength(2)
  expect(html).toContain('dateTime="2026-11-01T12:00:00Z"')
})
