import { useEffect, useState, useCallback } from 'react'
import type { Ticket } from '../types'
import * as ticketsApi from '../api/tickets'

export function useTickets() {
  const [tickets, setTickets] = useState<Ticket[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTickets = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await ticketsApi.getMyTickets()
      setTickets(data)
    } catch {
      setError('Error al cargar tus entradas')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchTickets()
  }, [fetchTickets])

  const cancelTicket = useCallback(async (id: number) => {
    try {
      await ticketsApi.cancelTicket(id)
      setTickets((prev) =>
        prev.map((t) =>
          t.id === id ? { ...t, status: 'cancelled' as const, cancelled_at: new Date().toISOString() } : t,
        ),
      )
      return true
    } catch {
      return false
    }
  }, [])

  const transferTicket = useCallback(async (id: number, toEmail: string) => {
    try {
      await ticketsApi.transferTicket(id, toEmail)
      setTickets((prev) =>
        prev.map((t) =>
          t.id === id ? { ...t, status: 'transferred' as const, transferred_at: new Date().toISOString() } : t,
        ),
      )
      return true
    } catch {
      return false
    }
  }, [])

  return { tickets, loading, error, refetch: fetchTickets, cancelTicket, transferTicket }
}
