import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { apiEjemplos, type CrearEjemplo } from './ejemplos.api'

export function useListarEjemplos(baseUrl: string) {
  return useQuery({
    queryKey: ['ejemplos', baseUrl],
    queryFn: () => apiEjemplos(baseUrl).listar(),
  })
}
export function useCrearEjemplo(baseUrl: string) {
  const clienteConsultas = useQueryClient()
  return useMutation({
    mutationFn: (entrada: CrearEjemplo) => apiEjemplos(baseUrl).crear(entrada),
    onSuccess: () =>
      clienteConsultas.invalidateQueries({ queryKey: ['ejemplos', baseUrl] }),
  })
}
