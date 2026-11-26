import type { CreateEventData } from '../types/admin'

export type FormErrors = Partial<Record<keyof CreateEventData, string>>

// Validaciones preventivas del formulario de eventos. Replican las reglas del
// backend (EventService.Create / Update y validatePresaleConfig) para avisar
// antes de enviar la request.
// hasExistingCode: al editar un evento que ya tenía preventa, el código es
// opcional (si se deja vacío, el backend conserva el actual).
export function validateEventForm(form: CreateEventData, isEdit: boolean, hasExistingCode: boolean): FormErrors {
  const errors: FormErrors = {}

  if (!form.title.trim()) {
    errors.title = 'El título es obligatorio'
  }

  const eventDate = form.event_date ? new Date(form.event_date) : null
  if (!eventDate || isNaN(eventDate.getTime())) {
    errors.event_date = 'La fecha del evento es obligatoria'
  } else if (!isEdit && eventDate.getTime() <= Date.now()) {
    errors.event_date = 'La fecha del evento tiene que ser futura'
  }

  if (!Number.isInteger(form.capacity) || form.capacity < 1) {
    errors.capacity = 'La capacidad tiene que ser un número entero mayor o igual a 1'
  }

  if (Number.isNaN(form.price) || form.price < 0) {
    errors.price = 'El precio no puede ser negativo'
  }

  if (
    form.duration_minutes != null &&
    (!Number.isInteger(form.duration_minutes) || form.duration_minutes < 1)
  ) {
    errors.duration_minutes = 'La duración tiene que ser un número entero de al menos 1 minuto'
  }

  const imageUrl = form.image_url?.trim()
  if (imageUrl && !imageUrl.startsWith('http://') && !imageUrl.startsWith('https://')) {
    errors.image_url = 'La URL tiene que empezar con http:// o https://'
  }

  if (form.presale_active) {
    if (!hasExistingCode && !form.presale_code?.trim()) {
      errors.presale_code = 'El código de preventa es obligatorio'
    }

    const presaleStart = form.presale_start_date ? new Date(form.presale_start_date) : null
    const generalSale = form.general_sale_date ? new Date(form.general_sale_date) : null

    if (!presaleStart) {
      errors.presale_start_date = 'La fecha de inicio de preventa es obligatoria'
    }
    if (!generalSale) {
      errors.general_sale_date = 'La fecha de venta general es obligatoria'
    }
    if (presaleStart && generalSale && presaleStart >= generalSale) {
      errors.presale_start_date = 'La preventa tiene que empezar antes de la venta general'
    }
    if (generalSale && eventDate && generalSale >= eventDate) {
      errors.general_sale_date = 'La venta general tiene que empezar antes del evento'
    }
  }

  return errors
}
