import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import { App } from "./App";
import "./estilos.css";

const clienteConsultas = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

async function iniciar() {
  if (debeUsarMocks()) {
    try {
      const { worker } = await import("./mocks/browser");
      await worker.start({ onUnhandledRequest: "bypass" });
    } catch (error) {
      console.error("No se pudo iniciar MSW; se usará la API real.", error);
    }
  }

  const contenedor = document.getElementById("root");
  if (!contenedor) {
    throw new Error("No se encontró el contenedor #root.");
  }

  createRoot(contenedor).render(
    <StrictMode>
      <QueryClientProvider client={clienteConsultas}>
        <App />
      </QueryClientProvider>
    </StrictMode>,
  );
}

function debeUsarMocks(): boolean {
  const porVariable = import.meta.env.VITE_USE_MOCKS === "true";
  const porConsulta = new URLSearchParams(window.location.search).has("mocks");
  return porVariable || porConsulta;
}

void iniciar();
