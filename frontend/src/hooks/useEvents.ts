import { useEffect, useState, useCallback } from 'react'
import type { Event, EventFilters } from '../types'
import { getEvents } from '../api/events'
import { mockEvents } from '../data/mockEvents'

export function useEvents(initialFilters?: EventFilters) {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [filters, setFilters] = useState<EventFilters | undefined>(initialFilters)

  const fetchEvents = useCallback(async (f?: EventFilters) => {
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

  useEffect(() => {
    fetchEvents(filters)
  }, [filters, fetchEvents])

  const updateFilters = useCallback((newFilters: EventFilters) => {
    setFilters(newFilters)
  }, [])

  const clearFilters = useCallback(() => {
    setFilters(undefined)
  }, [])

  return { events, loading, error, filters, updateFilters, clearFilters, refetch: fetchEvents }
}
