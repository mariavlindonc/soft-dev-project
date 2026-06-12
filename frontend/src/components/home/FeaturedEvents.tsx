import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import type { Event } from '../../types'
import { getEvents } from '../../api/events'
import { mockEvents } from '../../data/mockEvents'
import EventCard from '../events/EventCard'

export default function FeaturedEvents() {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getEvents()
      .then((data) => setEvents(data.slice(0, 4)))
      .catch(() => setEvents(mockEvents.slice(0, 4)))
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
          <EventCard key={event.id} event={event} />
        ))}
      </div>
    </section>
  )
}
