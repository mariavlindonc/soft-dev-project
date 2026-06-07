export function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleTimeString('es-ES', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatDateTime(dateStr: string): string {
  return `${formatDate(dateStr)}, ${formatTime(dateStr)}`
}

export function formatPrice(price: number): string {
  return price === 0 ? 'Gratis' : `$${price.toFixed(2)}`
}
