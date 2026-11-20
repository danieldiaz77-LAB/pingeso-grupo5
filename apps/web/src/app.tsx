import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState } from 'react'
import { Ejemplos } from './componentes/ejemplos'

export function App() {
  const [cliente] = useState(
    () =>
      new QueryClient({
        defaultOptions: { queries: { retry: 1, refetchOnWindowFocus: false } },
      }),
  )
  return (
    <QueryClientProvider client={cliente}>
      <main className="mx-auto min-h-screen max-w-5xl space-y-8 px-4 py-10">
        <header>
          <h1>Entorno técnico PINGESO</h1>
          <p>Contrato funcional pendiente de definición por el equipo.</p>
          <p>Los ejemplos siguientes pertenecen únicamente a la plantilla técnica.</p>
        </header>
        <Ejemplos baseUrl={import.meta.env.VITE_API_URL?.trim() ?? ''} />
      </main>
    </QueryClientProvider>
  )
}
