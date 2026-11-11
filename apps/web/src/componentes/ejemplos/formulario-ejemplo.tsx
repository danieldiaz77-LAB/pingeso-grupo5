import { Button } from "@base-ui/react/button";
import { useForm } from "react-hook-form";

import type { CrearEjemplo } from "../../api/ejemplos";
import { validarNombre } from "./utilidades";

interface Props {
  enviando: boolean;
  error?: string;
  alCrear: (entrada: CrearEjemplo) => Promise<void>;
}

export function FormularioEjemplo({ enviando, error, alCrear }: Props) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CrearEjemplo>({
    defaultValues: {
      nombre: "",
      descripcion: "",
    },
  });

  const enviar = handleSubmit(async (valores) => {
    try {
      await alCrear({
        nombre: valores.nombre.trim(),
        descripcion: valores.descripcion?.trim() ?? "",
      });
      reset();
    } catch {
      // El error se muestra desde el estado de la mutación.
    }
  });

  return (
    <form onSubmit={enviar} className="space-y-4" noValidate>
      <div>
        <label
          htmlFor="nombre"
          className="mb-1 block text-sm font-medium text-slate-700"
        >
          Nombre
        </label>
        <input
          id="nombre"
          type="text"
          autoComplete="off"
          aria-invalid={errors.nombre ? true : undefined}
          aria-describedby={errors.nombre ? "error-nombre" : undefined}
          className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm outline-none transition focus:border-sky-500 focus:ring-2 focus:ring-sky-200"
          {...register("nombre", { validate: validarNombre })}
        />
        {errors.nombre ? (
          <p id="error-nombre" className="mt-1 text-sm text-red-600">
            {errors.nombre.message}
          </p>
        ) : null}
      </div>

      <div>
        <label
          htmlFor="descripcion"
          className="mb-1 block text-sm font-medium text-slate-700"
        >
          Descripción
        </label>
        <textarea
          id="descripcion"
          rows={4}
          aria-invalid={errors.descripcion ? true : undefined}
          aria-describedby={errors.descripcion ? "error-descripcion" : undefined}
          className="w-full resize-y rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm outline-none transition focus:border-sky-500 focus:ring-2 focus:ring-sky-200"
          {...register("descripcion", {
            maxLength: {
              value: 500,
              message: "La descripción no puede superar los 500 caracteres.",
            },
          })}
        />
        {errors.descripcion ? (
          <p id="error-descripcion" className="mt-1 text-sm text-red-600">
            {errors.descripcion.message}
          </p>
        ) : null}
      </div>

      {error ? (
        <p role="alert" className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">
          {error}
        </p>
      ) : null}

      <Button
        type="submit"
        disabled={enviando}
        className="inline-flex w-full items-center justify-center rounded-lg bg-sky-700 px-4 py-2 text-sm font-semibold text-white transition hover:bg-sky-800 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {enviando ? "Creando..." : "Crear ejemplo"}
      </Button>
    </form>
  );
}
