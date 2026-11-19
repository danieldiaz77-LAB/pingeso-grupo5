import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { afterAll, afterEach, beforeAll, expect, test } from 'vitest'
import { apiEjemplos } from '../src/componentes/ejemplos/ejemplos.api'
import { handlers, reiniciarEjemplos } from '../src/mocks/handlers'

const servidor = setupServer(...handlers)
let cliente: ReturnType<typeof apiEjemplos>
beforeAll(() => {
  servidor.listen({ onUnhandledRequest: 'error' })
  cliente = apiEjemplos('http://localhost')
})
afterEach(() => {
  servidor.resetHandlers()
  reiniciarEjemplos()
})
afterAll(() => servidor.close())

test('crea mediante el cliente tipado y encuentra el resultado en el listado', async () => {
  const antes = await cliente.listar()
  const creado = await cliente.crear({
    nombre: 'Integración cliente',
    descripcion: 'Prueba MSW',
  })
  const despues = await cliente.listar()
  expect(despues).toHaveLength(antes.length + 1)
  expect(despues).toContainEqual(creado)
})

test('envía credenciales y utiliza la URL base recibida', async () => {
  servidor.use(
    http.get('http://otra-api.local/api/v1/ejemplos', ({ request }) => {
      expect(request.credentials).toBe('include')
      return HttpResponse.json({ ejemplos: [] })
    }),
  )
  await expect(apiEjemplos('http://otra-api.local').listar()).resolves.toEqual([])
})

test('propaga el rechazo de un nombre vacío sin añadir un registro', async () => {
  const antes = await cliente.listar()
  await expect(cliente.crear({ nombre: '   ' })).rejects.toThrow()
  await expect(cliente.listar()).resolves.toEqual(antes)
})

test.each([
  [
    { detail: 'Servicio temporalmente no disponible', title: 'Error' },
    'Servicio temporalmente no disponible',
  ],
  [{ title: 'Acceso denegado' }, 'Acceso denegado'],
  [{ detail: 123, title: false }, 'La API respondió con estado 503.'],
])(
  'informa errores HTTP sin convertirlos en un listado vacío: %j',
  async (cuerpo, mensaje) => {
    servidor.use(
      http.get('http://localhost/api/v1/ejemplos', () =>
        HttpResponse.json(cuerpo, { status: 503 }),
      ),
    )
    await expect(cliente.listar()).rejects.toThrow(mensaje)
  },
)

test('rechaza una respuesta exitosa sin cuerpo', async () => {
  servidor.use(
    http.get(
      'http://localhost/api/v1/ejemplos',
      () => new HttpResponse(null, { status: 204 }),
    ),
  )
  await expect(cliente.listar()).rejects.toThrow('La API respondió con estado 204.')
})
