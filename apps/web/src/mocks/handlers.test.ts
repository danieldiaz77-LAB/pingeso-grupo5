import { setupServer } from "msw/node";
import { afterAll, afterEach, beforeAll, describe, expect, it } from "vitest";

import { handlers } from "./handlers";

const servidor = setupServer(...handlers);

beforeAll(() => servidor.listen({ onUnhandledRequest: "error" }));
afterEach(() => servidor.resetHandlers());
afterAll(() => servidor.close());

describe("handlers MSW", () => {
  it("lista los ejemplos iniciales", async () => {
    const respuesta = await fetch("http://localhost/api/v1/ejemplos");
    expect(respuesta.status).toBe(200);

    const cuerpo = (await respuesta.json()) as { ejemplos: unknown[] };
    expect(cuerpo.ejemplos.length).toBeGreaterThanOrEqual(2);
  });

  it("crea un ejemplo y lo devuelve", async () => {
    const respuesta = await fetch("http://localhost/api/v1/ejemplos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ nombre: "Ejemplo mock", descripcion: "Desde MSW." }),
    });
    expect(respuesta.status).toBe(201);

    const cuerpo = (await respuesta.json()) as { nombre: string; id: string };
    expect(cuerpo.nombre).toBe("Ejemplo mock");
    expect(cuerpo.id).not.toBe("");
  });

  it("rechaza un nombre en blanco", async () => {
    const respuesta = await fetch("http://localhost/api/v1/ejemplos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ nombre: "   " }),
    });
    expect(respuesta.status).toBe(400);
  });
});
