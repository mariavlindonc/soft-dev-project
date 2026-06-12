import { useEffect, useState, useCallback } from 'react'
import type { Event, EventFilters } from '../types'
import { getEvents } from '../api/events'
import { mockEvents } from '../data/mockEvents'

export function useEvents(initialFilters?: EventFilters) {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [filters, setFilters] = useState<EventFilters | undefined>(initialFilters)

  useEffect(() => {
    let cancelled = false
    getEvents(filters)
      .then((data) => {
        if (!cancelled) {
          setEvents(data)
          setLoading(false)
        }
      })
      .catch(() => {
        if (!cancelled) {
          setEvents(mockEvents)
          setLoading(false)
        }
      })
    return () => { cancelled = true }
  }, [filters])

  const refetch = useCallback(async (f?: EventFilters) => {
    setLoading(true)
    setError(null)
    try {
      const data = await getEvents(f)
      setEvents(data)
    } catch {
      setEvents(mockEvents)
    } finally {
      setLoading(false)
    }
  }, [])

  const updateFilters = useCallback((newFilters: EventFilters) => {
    setFilters(newFilters)
  }, [])

  const clearFilters = useCallback(() => {
    setFilters(undefined)
  }, [])

  return { events, loading, error, filters, updateFilters, clearFilters, refetch }
}
