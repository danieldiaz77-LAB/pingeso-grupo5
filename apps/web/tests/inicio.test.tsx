import { afterEach, expect, test, vi } from 'vitest'
const { render, createRoot } = vi.hoisted(() => {
  const render = vi.fn()
  return { render, createRoot: vi.fn(() => ({ render })) }
})
vi.mock('react-dom/client', () => ({ createRoot }))
afterEach(() => {
  vi.unstubAllGlobals()
  vi.clearAllMocks()
  vi.resetModules()
})
test('monta React en el contenedor existente', async () => {
  const contenedor = { id: 'root' }
  vi.stubGlobal('document', { getElementById: () => contenedor })
  await import('../src/main')
  expect(createRoot).toHaveBeenCalledWith(contenedor)
  expect(render).toHaveBeenCalledOnce()
})
test('informa un HTML sin contenedor en vez de continuar sin interfaz', async () => {
  vi.stubGlobal('document', { getElementById: () => null })
  await expect(import('../src/main')).rejects.toThrow('Falta el contenedor')
  expect(createRoot).not.toHaveBeenCalled()
})
