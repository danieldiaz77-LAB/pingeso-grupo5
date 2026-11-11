import { cliente } from "./cliente";
import type { components } from "./esquema";

export type Ejemplo = components["schemas"]["Ejemplo"];
export type CrearEjemplo = components["schemas"]["CrearEjemplo"];

export async function listarEjemplos(): Promise<Ejemplo[]> {
  const { data, error, response } = await cliente.GET("/api/v1/ejemplos");
  if (error) {
    throw new Error(mensajeDeError(error, response.status));
  }
  return data?.ejemplos ?? [];
}

export async function crearEjemplo(entrada: CrearEjemplo): Promise<Ejemplo> {
  const { data, error, response } = await cliente.POST("/api/v1/ejemplos", {
    body: entrada,
  });
  if (error || !data) {
    throw new Error(mensajeDeError(error, response.status));
  }
  return data;
}

function mensajeDeError(error: unknown, estado: number): string {
  if (esProblema(error)) {
    return error.detail ?? error.title;
  }
  return `La API respondió con estado ${estado}.`;
}

function esProblema(valor: unknown): valor is components["schemas"]["Problema"] {
  return typeof valor === "object" && valor !== null && "title" in valor;
}
