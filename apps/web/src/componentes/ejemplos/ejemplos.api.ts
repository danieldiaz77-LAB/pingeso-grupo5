import createClient from 'openapi-fetch'
import type { paths } from '../../api/schema'

export type Ejemplo =
  paths['/api/v1/ejemplos']['get']['responses'][200]['content']['application/json']['ejemplos'][number]
export type CrearEjemplo =
  paths['/api/v1/ejemplos']['post']['requestBody']['content']['application/json']

export function apiEjemplos(baseUrl: string) {
  const cliente = createClient<paths>({ baseUrl, credentials: 'include' })
  return {
    async listar(): Promise<Ejemplo[]> {
      const { data, error, response } = await cliente.GET('/api/v1/ejemplos')
      if (error || !data) throw new Error(mensajeDeError(error, response.status))
      return data.ejemplos
    },
    async crear(entrada: CrearEjemplo): Promise<Ejemplo> {
      const { data, error, response } = await cliente.POST('/api/v1/ejemplos', {
        body: entrada,
      })
      if (error || !data) throw new Error(mensajeDeError(error, response.status))
      return data
    },
  }
}
function mensajeDeError(error: unknown, estado: number): string {
  if (typeof error === 'object' && error !== null) {
    if ('detail' in error && typeof error.detail === 'string') return error.detail
    if ('title' in error && typeof error.title === 'string') return error.title
  }
  return `La API respondió con estado ${estado}.`
}
