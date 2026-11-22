import { createRoot } from 'react-dom/client'
import { App } from './app'
import './estilos.css'

const contenedor = document.getElementById('root')
if (!contenedor) throw new Error('Falta el contenedor de la aplicación')
if (
  import.meta.env.DEV &&
  typeof window !== 'undefined' &&
  (import.meta.env.VITE_USE_MOCKS === 'true' ||
    new URLSearchParams(window.location.search).has('mocks'))
) {
  const { worker } = await import('./mocks/browser')
  await worker.start({ onUnhandledRequest: 'bypass' })
}
createRoot(contenedor).render(<App />)
