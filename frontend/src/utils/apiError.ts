import axios from 'axios'

// Traducciones de los errores del backend (services/errors.go) a mensajes
// para el usuario. El backend devuelve { "error": "<mensaje>" } y el mensaje
// contiene el texto del error centinela, p. ej. "event 3: event has no remaining capacity".
const ERROR_TRANSLATIONS: Array<[string, string]> = [
  // Compra
  ['no remaining capacity', 'No quedan entradas disponibles para este evento.'],
  ['event has been cancelled', 'El evento fue cancelado.'],
  ['sales are not yet open', 'La venta para este evento todavía no está abierta.'],
  ['presale code is required', 'Necesitás un código de preventa para comprar.'],
  ['invalid presale code', 'El código de preventa no es válido.'],
  ['resource not found', 'No se encontró el evento.'],
  // Entradas (cancelar / transferir)
  ['ticket is already cancelled', 'La entrada ya estaba cancelada.'],
  ['already been transferred', 'La entrada ya fue transferida.'],
  ['cannot transfer to the same user', 'No podés transferirte una entrada a vos mismo.'],
  // El backend responde "ticket not found" también cuando el email destino no existe.
  ['ticket not found', 'No existe un usuario registrado con ese correo, o la entrada ya no está disponible.'],
]

// Devuelve un mensaje en castellano para el error de una request a la API.
// Si el error no es conocido, devuelve `fallback`: nunca muestra el texto
// crudo del backend para no exponer detalles internos.
export function getApiErrorMessage(err: unknown, fallback: string): string {
  if (!axios.isAxiosError<{ error?: string }>(err)) return fallback

  const backendMessage = err.response?.data?.error
  if (typeof backendMessage !== 'string') return fallback

  const match = ERROR_TRANSLATIONS.find(([text]) => backendMessage.includes(text))
  return match ? match[1] : fallback
}
