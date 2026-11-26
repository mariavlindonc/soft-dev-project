import { describe, it, expect } from 'vitest'
import { getEventImage } from './eventImages'

describe('getEventImage', () => {
  it('usa la imagen propia del evento si la URL es válida', () => {
    expect(getEventImage(1, 'https://cdn.ejemplo.com/foto.jpg')).toBe('https://cdn.ejemplo.com/foto.jpg')
  })

  it('ignora los espacios alrededor de la URL', () => {
    expect(getEventImage(1, '  https://cdn.ejemplo.com/foto.jpg  ')).toBe('https://cdn.ejemplo.com/foto.jpg')
  })

  it('usa una imagen de respaldo si el evento no tiene imagen', () => {
    const image = getEventImage(1, null)
    expect(image).toBeTruthy()
    expect(image).not.toBe('')
  })

  it('usa una imagen de respaldo si la URL no empieza con http', () => {
    expect(getEventImage(1, '/images/pena.jpg')).toBe(getEventImage(1, null))
  })

  it('el mismo evento siempre recibe la misma imagen de respaldo', () => {
    expect(getEventImage(7, null)).toBe(getEventImage(7, undefined))
    expect(getEventImage(7, '')).toBe(getEventImage(7, null))
  })

  it('reparte las imágenes según el id (id % cantidad de imágenes)', () => {
    // Hay 8 imágenes: el 3 y el 11 comparten foto, el 3 y el 4 no.
    expect(getEventImage(3)).toBe(getEventImage(11))
    expect(getEventImage(3)).not.toBe(getEventImage(4))
  })
})
