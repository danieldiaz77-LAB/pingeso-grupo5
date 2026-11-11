import { Button } from "@base-ui/react/button";

import { FormularioEjemplo } from "./formulario-ejemplo";
import { ListaEjemplos } from "./lista-ejemplos";
import { useCrearEjemplo, useListarEjemplos } from "./use-ejemplos";

export function Ejemplos() {
  const consulta = useListarEjemplos();
  const mutacion = useCrearEjemplo();

  const ejemplos = consulta.data ?? [];
  const errorConsulta =
    consulta.error instanceof Error
      ? consulta.error.message
      : "No se pudo cargar el listado.";
  const errorMutacion =
    mutacion.error instanceof Error ? mutacion.error.message : undefined;

  return (
    <section className="rounded-2xl border border-slate-200 bg-slate-50 p-5 shadow-sm sm:p-6">
      <div className="mb-5 flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h2 className="text-xl font-semibold text-slate-900">Ejemplos</h2>
          <p className="mt-1 text-sm text-slate-600">
            Microcomponente conectado a la API de la plantilla.
          </p>
        </div>
        <span className="text-sm text-slate-500">
          {consulta.isPending ? "Cargando..." : `${ejemplos.length} registros`}
        </span>
      </div>

      <div className="grid gap-5 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.4fr)]">
        <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
          <h3 className="mb-3 text-base font-semibold text-slate-900">
            Nuevo ejemplo
          </h3>
          <FormularioEjemplo
            enviando={mutacion.isPending}
            error={errorMutacion}
            alCrear={async (entrada) => {
              await mutacion.mutateAsync(entrada);
            }}
          />
        </div>

        <div>
          <div className="mb-3 flex items-center justify-between">
            <h3 className="text-base font-semibold text-slate-900">Listado</h3>
            <Button
              type="button"
              onClick={() => void consulta.refetch()}
              disabled={consulta.isFetching}
              className="rounded-lg border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {consulta.isFetching ? "Actualizando..." : "Actualizar"}
            </Button>
          </div>

          {consulta.isPending ? (
            <p className="rounded-lg border border-slate-200 bg-white px-4 py-6 text-center text-sm text-slate-500">
              Cargando ejemplos...
            </p>
          ) : null}

          {consulta.isError ? (
            <p
              role="alert"
              className="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700"
            >
              {errorConsulta} Revisa que la API esté en ejecución.
            </p>
          ) : null}

          {consulta.isSuccess ? <ListaEjemplos ejemplos={ejemplos} /> : null}
        </div>
      </div>
    </section>
  );
}
