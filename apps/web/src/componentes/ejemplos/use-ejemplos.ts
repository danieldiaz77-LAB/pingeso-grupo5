import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  crearEjemplo,
  listarEjemplos,
  type CrearEjemplo,
} from "../../api/ejemplos";

const claveEjemplos = ["ejemplos"] as const;

export function useListarEjemplos() {
  return useQuery({
    queryKey: claveEjemplos,
    queryFn: listarEjemplos,
  });
}

export function useCrearEjemplo() {
  const clienteConsultas = useQueryClient();

  return useMutation({
    mutationFn: (entrada: CrearEjemplo) => crearEjemplo(entrada),
    onSuccess: () => clienteConsultas.invalidateQueries({ queryKey: claveEjemplos }),
  });
}
