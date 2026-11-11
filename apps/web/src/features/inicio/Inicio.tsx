import { Ejemplos } from "../../componentes/ejemplos";

export function Inicio() {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-5xl flex-col gap-8 px-4 py-10 sm:px-6">
      <header className="rounded-2xl bg-slate-900 px-6 py-8 text-white shadow-lg">
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-sky-300">
          PINGESO · Grupo 5
        </p>
        <h1 className="mt-3 text-3xl font-bold sm:text-4xl">
          Plantilla técnica del microservicio
        </h1>
        <p className="mt-3 max-w-2xl text-sm text-slate-300 sm:text-base">
          Demostración ejecutable de la arquitectura del DAS: React, TypeScript, Vite,
          Tailwind, React Query, React Hook Form, Go, Huma, PostgreSQL, Goose y sqlc.
        </p>
      </header>

      <Ejemplos />

      <footer className="mt-auto border-t border-slate-200 pt-4 text-xs text-slate-500">
        Prototipo de la plantilla. La lógica de ranking y rentabilidad se implementa en
        una etapa posterior.
      </footer>
    </main>
  );
}
