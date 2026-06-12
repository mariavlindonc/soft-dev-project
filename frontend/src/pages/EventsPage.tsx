import { useState, useEffect, useRef } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { useEvents } from '../hooks/useEvents'
import EventCard from '../components/events/EventCard'

const CATEGORIES = [
  { value: 'aire libre', label: 'Aire Libre' },
  { value: 'en salon', label: 'En Salón' },
  { value: 'musica', label: 'Música' },
  { value: 'teatro', label: 'Teatro' },
  { value: 'gastronomia', label: 'Gastronomía' },
]

export default function EventsPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const initialCategory = searchParams.get('category') ?? undefined
  const { events, loading, error, updateFilters } = useEvents(
    initialCategory ? { category: initialCategory } : undefined
  )

  const [search, setSearch] = useState(searchParams.get('q') ?? '')
  const [category, setCategory] = useState(initialCategory ?? '')
  const [debouncedSearch, setDebouncedSearch] = useState(search)

  const [dateFrom, setDateFrom] = useState(searchParams.get('date_from') ?? '')
  const [dateTo, setDateTo] = useState(searchParams.get('date_to') ?? '')
  const [minPrice, setMinPrice] = useState(searchParams.get('min_price') ?? '')
  const [maxPrice, setMaxPrice] = useState(searchParams.get('max_price') ?? '')

  const initialised = useRef(false)

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedSearch(search), 300)
    return () => clearTimeout(timer)
  }, [search])

  useEffect(() => {
    if (!initialised.current) {
      initialised.current = true
      return
    }
    const params: Record<string, string> = {}
    if (debouncedSearch) params.q = debouncedSearch
    if (category) params.category = category
    if (dateFrom) params.date_from = dateFrom
    if (dateTo) params.date_to = dateTo
    if (minPrice) params.min_price = minPrice
    if (maxPrice) params.max_price = maxPrice
    setSearchParams(params, { replace: true })

    const filters: Record<string, unknown> = {}
    if (category) filters.category = category
    if (dateFrom) filters.date_from = dateFrom
    if (dateTo) filters.date_to = dateTo
    if (minPrice) filters.min_price = Number(minPrice)
    if (maxPrice) filters.max_price = Number(maxPrice)
    updateFilters(filters as Parameters<typeof updateFilters>[0])
  }, [debouncedSearch, category, dateFrom, dateTo, minPrice, maxPrice, setSearchParams, updateFilters])

  const filtered = events.filter((e) => {
    const matchesCategory = !category || e.category === category
    if (!debouncedSearch && !dateFrom && !dateTo && !minPrice && !maxPrice) return matchesCategory

    const q = debouncedSearch.toLowerCase()
    const matchesSearch =
      !debouncedSearch ||
      e.title.toLowerCase().includes(q) ||
      (e.description?.toLowerCase().includes(q) ?? false) ||
      (e.location?.toLowerCase().includes(q) ?? false)

    const eventDate = e.event_date ? new Date(e.event_date) : null
    const matchesDateFrom = !dateFrom || !eventDate || eventDate >= new Date(dateFrom)
    const matchesDateTo = !dateTo || !eventDate || eventDate <= new Date(dateTo + 'T23:59:59')

    const matchesPriceMin = !minPrice || e.price >= Number(minPrice)
    const matchesPriceMax = !maxPrice || e.price <= Number(maxPrice)

    return matchesCategory && matchesSearch && matchesDateFrom && matchesDateTo && matchesPriceMin && matchesPriceMax
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
        <div className="events-filters-row">
          <input
            type="date"
            value={dateFrom}
            onChange={(e) => setDateFrom(e.target.value)}
            title="Desde fecha"
          />
          <input
            type="date"
            value={dateTo}
            onChange={(e) => setDateTo(e.target.value)}
            title="Hasta fecha"
          />
          <input
            type="number"
            placeholder="Precio min."
            value={minPrice}
            onChange={(e) => setMinPrice(e.target.value)}
            min="0"
          />
          <input
            type="number"
            placeholder="Precio máx."
            value={maxPrice}
            onChange={(e) => setMaxPrice(e.target.value)}
            min="0"
          />
        </div>
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
