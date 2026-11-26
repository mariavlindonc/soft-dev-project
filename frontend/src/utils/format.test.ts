import { describe, it, expect } from 'vitest'
import { formatPrice, formatDate, formatTime } from './format'

describe('formatPrice', () => {
  it('muestra "Gratis" para los eventos con precio 0', () => {
    expect(formatPrice(0)).toBe('Gratis')
  })

  it('formatea el precio con separador de miles argentino', () => {
    expect(formatPrice(25000)).toBe('$ 25.000')
  })
})

describe('formatDate y formatTime', () => {
  it('muestra la fecha en castellano', () => {
    expect(formatDate('2026-08-15T20:00:00')).toBe('15 de agosto de 2026')
  })

  it('muestra la hora con dos dígitos', () => {
    expect(formatTime('2026-08-15T20:05:00')).toBe('20:05')
  })
})
