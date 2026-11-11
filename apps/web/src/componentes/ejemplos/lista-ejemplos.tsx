import type { Ejemplo } from "../../api/ejemplos";
import { formatearFecha } from "./utilidades";

interface Props {
  ejemplos: Ejemplo[];
}

export function ListaEjemplos({ ejemplos }: Props) {
  if (ejemplos.length === 0) {
    return (
      <p className="rounded-lg border border-dashed border-slate-300 px-4 py-6 text-center text-sm text-slate-500">
        Todavía no hay ejemplos. Crea el primero desde el formulario.
      </p>
    );
  }

  return (
    <ul className="space-y-3">
      {ejemplos.map((ejemplo) => (
        <li
          key={ejemplo.id}
          className="rounded-lg border border-slate-200 bg-white px-4 py-3 shadow-sm"
        >
          <div className="flex flex-col gap-1 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h3 className="text-base font-semibold text-slate-900">
                {ejemplo.nombre}
              </h3>
              <p className="mt-1 text-sm text-slate-600">
                {ejemplo.descripcion || "Sin descripción."}
              </p>
            </div>
            <time
              dateTime={ejemplo.creadoEn}
              className="shrink-0 text-xs text-slate-500"
            >
              {formatearFecha(ejemplo.creadoEn)}
            </time>
          </div>
        </li>
      ))}
    </ul>
  );
}
