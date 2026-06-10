import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import type { Event } from '../../types'
import { getEvents } from '../../api/events'
import { formatPrice } from '../../utils/format'

export default function FeaturedEvents() {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getEvents()
      .then((data) => setEvents(data.slice(0, 4)))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <section className="featured-section">
        <div className="section-header">
          <h2>Eventos Destacados</h2>
          <Link to="/events" className="btn btn-outline">Ver todos</Link>
        </div>
        <div className="event-grid__grid">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="event-card skeleton">
              <div className="event-card__image-wrapper" />
            </div>
          ))}
        </div>
      </section>
    )
  }

  if (events.length === 0) return null

  return (
    <section className="featured-section">
      <div className="section-header">
        <h2>Eventos Destacados</h2>
        <Link to="/events" className="btn btn-outline">Ver todos</Link>
      </div>
      <div className="event-grid__grid">
        {events.map((event) => (
          <Link key={event.id} to={`/events/${event.id}`} className="event-card">
            <div className="event-card__image-wrapper">
              <div
                className="event-card__image"
                style={{ backgroundImage: event.image_url ? `url(${event.image_url})` : undefined }}
              />
              <button type="button" className="event-card__quick-action">Entradas</button>
            </div>
            <div className="event-card__info">
              <span className="event-card__category">{event.category ?? 'General'}</span>
              <h3 className="event-card__title">{event.title}</h3>
              <p className="event-card__meta">
                {new Date(event.event_date).toLocaleDateString('es-ES', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
                {event.location ? ` · ${event.location}` : ''}
              </p>
              <span className="event-card__price">{formatPrice(event.price)}</span>
            </div>
          </Link>
        ))}
      </div>
    </section>
  )
}
