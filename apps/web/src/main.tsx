import { createRoot } from 'react-dom/client'
import { App } from './app'
const contenedor = document.getElementById('root')
if (!contenedor) throw new Error('Falta el contenedor de la aplicación')
createRoot(contenedor).render(<App />)
