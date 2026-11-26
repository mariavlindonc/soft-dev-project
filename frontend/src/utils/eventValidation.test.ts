import { describe, it, expect } from 'vitest'
import { validateEventForm } from './eventValidation'
import type { CreateEventData } from '../types/admin'

// Convierte una fecha al formato del input datetime-local ("YYYY-MM-DDTHH:mm").
function toInput(date: Date): string {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

const DAY = 24 * 60 * 60 * 1000
const inDays = (days: number) => toInput(new Date(Date.now() + days * DAY))

// Formulario válido: cada test cambia solo lo que quiere probar.
function validForm(overrides: Partial<CreateEventData> = {}): CreateEventData {
  return {
    title: 'Recital de prueba',
    description: '',
    image_url: '',
    category: 'musica',
    location: 'Córdoba',
    event_date: inDays(30),
    duration_minutes: 120,
    capacity: 100,
    price: 5000,
    presale_active: false,
    presale_code: '',
    presale_start_date: '',
    general_sale_date: '',
    ...overrides,
  }
}

const presale = (overrides: Partial<CreateEventData> = {}) =>
  validForm({
    presale_active: true,
    presale_code: 'PREVENTA2026',
    presale_start_date: inDays(1),
    general_sale_date: inDays(10),
    ...overrides,
  })

describe('validateEventForm', () => {
  it('un formulario válido no tiene errores', () => {
    expect(validateEventForm(validForm(), false, false)).toEqual({})
  })

  it('el título es obligatorio (los espacios no cuentan)', () => {
    expect(validateEventForm(validForm({ title: '   ' }), false, false).title).toBeDefined()
  })

  it('la fecha del evento es obligatoria', () => {
    expect(validateEventForm(validForm({ event_date: '' }), false, false).event_date).toBeDefined()
  })

  it('al crear, la fecha del evento tiene que ser futura', () => {
    const errors = validateEventForm(validForm({ event_date: inDays(-1) }), false, false)
    expect(errors.event_date).toBe('La fecha del evento tiene que ser futura')
  })

  it('al editar, se puede guardar un evento con fecha pasada', () => {
    expect(validateEventForm(validForm({ event_date: inDays(-1) }), true, false)).toEqual({})
  })

  it('la capacidad tiene que ser un entero mayor o igual a 1', () => {
    expect(validateEventForm(validForm({ capacity: 0 }), false, false).capacity).toBeDefined()
    expect(validateEventForm(validForm({ capacity: 2.5 }), false, false).capacity).toBeDefined()
  })

  it('el precio no puede ser negativo, pero 0 (evento gratis) es válido', () => {
    expect(validateEventForm(validForm({ price: -1 }), false, false).price).toBeDefined()
    expect(validateEventForm(validForm({ price: 0 }), false, false)).toEqual({})
  })

  it('la duración es opcional, pero si viene tiene que ser al menos 1 minuto', () => {
    expect(validateEventForm(validForm({ duration_minutes: 0 }), false, false).duration_minutes).toBeDefined()
    expect(validateEventForm(validForm({ duration_minutes: undefined }), false, false)).toEqual({})
  })

  it('la URL de imagen tiene que empezar con http:// o https://', () => {
    expect(validateEventForm(validForm({ image_url: 'foto.jpg' }), false, false).image_url).toBeDefined()
    expect(validateEventForm(validForm({ image_url: 'https://cdn.ejemplo.com/a.jpg' }), false, false)).toEqual({})
  })

  describe('con preventa activa', () => {
    it('una preventa bien configurada no tiene errores', () => {
      expect(validateEventForm(presale(), false, false)).toEqual({})
    })

    it('el código y las dos fechas son obligatorios', () => {
      const errors = validateEventForm(
        presale({ presale_code: '', presale_start_date: '', general_sale_date: '' }),
        false,
        false,
      )
      expect(errors.presale_code).toBeDefined()
      expect(errors.presale_start_date).toBeDefined()
      expect(errors.general_sale_date).toBeDefined()
    })

    it('la preventa tiene que empezar antes de la venta general', () => {
      const errors = validateEventForm(presale({ presale_start_date: inDays(12) }), false, false)
      expect(errors.presale_start_date).toBe('La preventa tiene que empezar antes de la venta general')
    })

    it('la venta general tiene que empezar antes del evento', () => {
      const errors = validateEventForm(presale({ general_sale_date: inDays(40) }), false, false)
      expect(errors.general_sale_date).toBe('La venta general tiene que empezar antes del evento')
    })

    it('al editar un evento que ya tenía código, no hace falta reingresarlo', () => {
      expect(validateEventForm(presale({ presale_code: '' }), true, true)).toEqual({})
    })
  })
})
