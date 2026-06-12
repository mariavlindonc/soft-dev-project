import { useMemo } from 'react'
import { Link } from 'react-router-dom'
import type { Event } from '../../types'
import { formatDateTime, formatPrice } from '../../utils/format'
import { getRandomEventImage } from '../../data/eventImages'

interface EventCardProps {
  event: Event
}

export default function EventCard({ event }: EventCardProps) {
  const image = useMemo(getRandomEventImage, [])

  const statusLabel: Record<string, string> = {
    active: 'Activo',
    presale: 'Preventa',
    sold_out: 'Agotado',
    cancelled: 'Cancelado',
  }

  const showStatus = event.status === 'sold_out' || event.status === 'cancelled'

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
        {showStatus && (
          <span className={`event-card__status-badge event-card__status-badge--${event.status}`}>
            {statusLabel[event.status]}
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
