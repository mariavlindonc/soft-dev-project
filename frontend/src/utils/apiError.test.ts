import { describe, it, expect } from 'vitest'
import { AxiosError } from 'axios'
import { getApiErrorMessage } from './apiError'

const FALLBACK = 'Error genérico'

// Simula el error que devuelve axios cuando el backend responde { error: "..." }.
function apiError(backendMessage: unknown) {
  return Object.assign(new AxiosError('Request failed'), {
    response: { data: { error: backendMessage } },
  })
}

describe('getApiErrorMessage', () => {
  it.each([
    ['event 3: event has no remaining capacity', 'No quedan entradas disponibles para este evento.'],
    ['event 3: event has been cancelled', 'El evento fue cancelado.'],
    ['event 3: sales are not yet open for this event', 'La venta para este evento todavía no está abierta.'],
    ['event 3: presale code is required', 'Necesitás un código de preventa para comprar.'],
    ['event 3: invalid presale code', 'El código de preventa no es válido.'],
    ['event 99: resource not found', 'No se encontró el evento.'],
    ['ticket 5: ticket is already cancelled', 'La entrada ya estaba cancelada.'],
    ['ticket 5: ticket has already been transferred', 'La entrada ya fue transferida.'],
    ['ticket 5: cannot transfer to the same user', 'No podés transferirte una entrada a vos mismo.'],
  ])('traduce "%s"', (backendMessage, expected) => {
    expect(getApiErrorMessage(apiError(backendMessage), FALLBACK)).toBe(expected)
  })

  it('devuelve el mensaje genérico si el error del backend no es conocido', () => {
    expect(getApiErrorMessage(apiError('some internal sql error'), FALLBACK)).toBe(FALLBACK)
  })

  it('nunca muestra el texto crudo del backend', () => {
    const message = getApiErrorMessage(apiError('pq: connection refused at 10.0.0.1'), FALLBACK)
    expect(message).not.toContain('10.0.0.1')
  })

  it('devuelve el mensaje genérico si no hubo respuesta del servidor (backend caído)', () => {
    expect(getApiErrorMessage(new AxiosError('Network Error'), FALLBACK)).toBe(FALLBACK)
  })

  it('devuelve el mensaje genérico si el error no viene de axios', () => {
    expect(getApiErrorMessage(new Error('boom'), FALLBACK)).toBe(FALLBACK)
    expect(getApiErrorMessage('texto suelto', FALLBACK)).toBe(FALLBACK)
  })

  it('devuelve el mensaje genérico si el campo error no es texto', () => {
    expect(getApiErrorMessage(apiError({ code: 500 }), FALLBACK)).toBe(FALLBACK)
  })
})
