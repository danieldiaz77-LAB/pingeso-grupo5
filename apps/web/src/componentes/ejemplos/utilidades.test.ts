import { describe, expect, it } from "vitest";

import { formatearFecha, validarNombre } from "./utilidades";

describe("validarNombre", () => {
  it("acepta un nombre con contenido", () => {
    expect(validarNombre("Primer ejemplo")).toBeUndefined();
  });

  it("rechaza un nombre vacío o en blanco", () => {
    expect(validarNombre("")).toBe("El nombre es obligatorio.");
    expect(validarNombre("   ")).toBe("El nombre es obligatorio.");
  });

  it("rechaza un nombre de más de 120 caracteres", () => {
    expect(validarNombre("a".repeat(121))).toBe(
      "El nombre no puede superar los 120 caracteres.",
    );
  });
});

describe("formatearFecha", () => {
  it("formatea una fecha RFC 3339", () => {
    expect(formatearFecha("2026-01-01T12:00:00Z")).toContain("2026");
  });

  it("informa una fecha inválida", () => {
    expect(formatearFecha("no-es-una-fecha")).toBe("Fecha inválida");
  });
});
