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
        <div className="event-grid">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="event-card skeleton" />
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
      <div className="event-grid">
        {events.map((event) => (
          <Link key={event.id} to={`/events/${event.id}`} className="event-card">
            <div
              className="event-card-image"
              style={{ backgroundImage: event.image_url ? `url(${event.image_url})` : undefined }}
            />
            <div className="event-card-body">
              <span className="event-category">{event.category ?? 'General'}</span>
              <h3>{event.title}</h3>
              <p className="event-date">
                {new Date(event.event_date).toLocaleDateString('es-ES', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </p>
              <p className="event-location">{event.location}</p>
              <span className="event-price">{formatPrice(event.price)}</span>
            </div>
          </Link>
        ))}
      </div>
    </section>
  )
}
