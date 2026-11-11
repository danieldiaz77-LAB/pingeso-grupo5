const formateador = new Intl.DateTimeFormat("es-CL", {
  dateStyle: "medium",
  timeStyle: "short",
  timeZone: "America/Santiago",
});

export function formatearFecha(valor: string): string {
  const fecha = new Date(valor);
  if (Number.isNaN(fecha.getTime())) {
    return "Fecha inválida";
  }
  return formateador.format(fecha);
}

export function validarNombre(valor: string): string | undefined {
  const nombre = valor.trim();
  if (!nombre) {
    return "El nombre es obligatorio.";
  }
  if ([...nombre].length > 120) {
    return "El nombre no puede superar los 120 caracteres.";
  }
  return undefined;
}
