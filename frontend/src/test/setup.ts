// Configuración común de los tests (Vitest + React Testing Library).
import '@testing-library/jest-dom/vitest'
import { afterEach } from 'vitest'
import { cleanup } from '@testing-library/react'

// Desmonta los componentes renderizados después de cada test.
afterEach(() => {
  cleanup()
})
