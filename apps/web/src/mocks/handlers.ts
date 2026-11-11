import { delay, http, HttpResponse } from "msw";

import type { CrearEjemplo, Ejemplo } from "../api/ejemplos";
import { ejemplosIniciales } from "./datos";

let ejemplos: Ejemplo[] = [...ejemplosIniciales];

const origen = typeof window === "undefined" ? "http://localhost" : window.location.origin;
const rutaEjemplos = `${origen}/api/v1/ejemplos`;

export const handlers = [
  http.get(rutaEjemplos, async () => {
    await delay(250);
    return HttpResponse.json({ ejemplos });
  }),

  http.post(rutaEjemplos, async ({ request }) => {
    const entrada = (await request.json()) as CrearEjemplo;
    const nombre = (entrada.nombre ?? "").trim();

    if (!nombre) {
      return HttpResponse.json(
        {
          type: "about:blank",
          title: "Solicitud inválida",
          status: 400,
          detail: "El nombre es obligatorio.",
        },
        {
          status: 400,
          headers: { "Content-Type": "application/problem+json" },
        },
      );
    }

    const nuevo: Ejemplo = {
      id: crypto.randomUUID(),
      nombre,
      descripcion: (entrada.descripcion ?? "").trim(),
      creadoEn: new Date().toISOString(),
    };
    ejemplos = [...ejemplos, nuevo];

    await delay(200);
    return HttpResponse.json(nuevo, { status: 201 });
  }),
];
