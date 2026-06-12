import { useMemo } from 'react'
import { Link } from 'react-router-dom'
import type { Event } from '../../types'
import { formatDateTime, formatPrice } from '../../utils/format'
import { getRandomEventImage } from '../../data/eventImages'

interface EventCardProps {
  event: Event
}

type CardBadge = { key: string; label: string } | null

function getBadge(event: Event): CardBadge {
  if (event.status === 'cancelled') {
    return { key: 'cancelled', label: 'Cancelado' }
  }
  if (event.tickets_sold >= event.capacity) {
    return { key: 'sold_out', label: 'Agotado' }
  }
  if (event.presale_active && event.presale_start_date && event.general_sale_date) {
    const now = new Date()
    const presaleStart = new Date(event.presale_start_date)
    const generalSale = new Date(event.general_sale_date)
    if (now < presaleStart) {
      return { key: 'upcoming', label: 'Próximamente' }
    }
    if (now >= presaleStart && now < generalSale) {
      return { key: 'presale', label: 'Preventa' }
    }
  }
  return null
}

export default function EventCard({ event }: EventCardProps) {
  const image = useMemo(getRandomEventImage, [])
  const badge = getBadge(event)

  return (
    <Link to={`/events/${event.id}`} className="event-card">
      <div className="event-card__image-wrapper">
        <div
          className="event-card__image"
          style={{ backgroundImage: `url(${image})` }}
        />
        {event.category && (
          <span className="event-card__category-badge">{event.category}</span>
        )}
        {badge && (
          <span className={`event-card__status-badge event-card__status-badge--${badge.key}`}>
            {badge.label}
          </span>
        )}
      </div>
      <div className="event-card__info">
        <h3 className="event-card__title">{event.title}</h3>
        <p className="event-card__meta">
          {formatDateTime(event.event_date)}
          {event.location ? ` · ${event.location}` : ''}
        </p>
        <span className="event-card__price">{formatPrice(event.price)}</span>
      </div>
    </Link>
  )
}
