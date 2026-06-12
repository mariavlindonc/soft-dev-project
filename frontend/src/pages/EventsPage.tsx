import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { useEvents } from '../hooks/useEvents'
import { formatDate, formatPrice } from '../utils/format'

const CATEGORIES = [
  { value: 'concierto', label: 'Conciertos' },
  { value: 'teatro', label: 'Teatro' },
  { value: 'deporte', label: 'Deportes' },
  { value: 'conferencia', label: 'Conferencias' },
  { value: 'familiar', label: 'Familiares' },
  { value: 'feria', label: 'Ferias' },
]

export default function EventsPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const { events, loading, error, updateFilters } = useEvents()

  const [search, setSearch] = useState(searchParams.get('q') ?? '')
  const [category, setCategory] = useState(searchParams.get('category') ?? '')
  const [debouncedSearch, setDebouncedSearch] = useState(search)

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedSearch(search), 300)
    return () => clearTimeout(timer)
  }, [search])

  useEffect(() => {
    const params: Record<string, string> = {}
    if (debouncedSearch) params.q = debouncedSearch
    if (category) params.category = category
    setSearchParams(params, { replace: true })
    updateFilters({ category: category || undefined })
  }, [debouncedSearch, category, setSearchParams, updateFilters])

  const filtered = events.filter((e) => {
    if (!debouncedSearch) return true
    const q = debouncedSearch.toLowerCase()
    return (
      e.title.toLowerCase().includes(q) ||
      (e.description?.toLowerCase().includes(q) ?? false) ||
      (e.location?.toLowerCase().includes(q) ?? false)
    )
  })

  return (
    <div className="page">
      <div className="events-page-header">
        <h1>Eventos</h1>
      </div>

      <div className="events-filters">
        <input
          type="text"
          placeholder="Buscar eventos…"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <select value={category} onChange={(e) => setCategory(e.target.value)}>
          <option value="">Todas las categorías</option>
          {CATEGORIES.map((c) => (
            <option key={c.value} value={c.value}>{c.label}</option>
          ))}
        </select>
      </div>

      {loading && (
        <div className="event-grid">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="event-card skeleton" />
          ))}
        </div>
      )}

      {error && <div className="alert alert-error">{error}</div>}

      {!loading && !error && (
        <>
          <p className="events-count">
            {filtered.length} {filtered.length === 1 ? 'evento encontrado' : 'eventos encontrados'}
          </p>

          {filtered.length === 0 ? (
            <div className="events-empty">
              <p>No se encontraron eventos</p>
              <span className="form-footer">
                Probá con otros filtros o <Link to="/events">limpiá la búsqueda</Link>
              </span>
            </div>
          ) : (
            <div className="event-grid">
              {filtered.map((event) => (
                <Link key={event.id} to={`/events/${event.id}`} className="event-card">
                  <div
                    className="event-card-image"
                    style={{ backgroundImage: event.image_url ? `url(${event.image_url})` : undefined }}
                  />
                  <div className="event-card-body">
                    <span className="event-category">{event.category ?? 'General'}</span>
                    <h3>{event.title}</h3>
                    <p className="event-date">{formatDate(event.event_date)}</p>
                    <p className="event-location">{event.location}</p>
                    <span className="event-price">{formatPrice(event.price)}</span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
}
