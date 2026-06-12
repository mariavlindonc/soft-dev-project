import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { useEvents } from '../hooks/useEvents'
import EventCard from '../components/events/EventCard'

const CATEGORIES = [
  { value: 'aire libre', label: 'Aire Libre' },
  { value: 'en salon', label: 'En Salón' },
  { value: 'grupos emergentes', label: 'Grupos Emergentes' },
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
                <EventCard key={event.id} event={event} />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
}
