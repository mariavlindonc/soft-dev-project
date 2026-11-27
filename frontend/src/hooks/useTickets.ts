import { useEffect, useState, useCallback } from 'react'
import type { Ticket } from '../types'
import * as ticketsApi from '../api/tickets'
import { getApiErrorMessage } from '../utils/apiError'

export function useTickets() {
  const [tickets, setTickets] = useState<Ticket[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    ticketsApi.getMyTickets()
      .then((data) => {
        if (!cancelled) {
          setTickets(data)
          setLoading(false)
        }
      })
      .catch(() => {
        if (!cancelled) {
          setError('Error al cargar tus entradas')
          setLoading(false)
        }
      })
    return () => { cancelled = true }
  }, [])

  const refetch = useCallback(async () => {
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

  // Devuelven null si salió bien, o el mensaje de error para mostrar.
  const cancelTicket = useCallback(async (id: number): Promise<string | null> => {
    try {
      await ticketsApi.cancelTicket(id)
      setTickets((prev) =>
        prev.map((t) =>
          t.id === id ? { ...t, status: 'cancelled' as const, cancelled_at: new Date().toISOString() } : t,
        ),
      )
      return null
    } catch (err) {
      return getApiErrorMessage(err, 'No se pudo cancelar la entrada')
    }
  }, [])

  const transferTicket = useCallback(async (id: number, toEmail: string): Promise<string | null> => {
    try {
      await ticketsApi.transferTicket(id, toEmail)
      setTickets((prev) =>
        prev.map((t) =>
          t.id === id ? { ...t, status: 'transferred' as const, transferred_out: true, transferred_at: new Date().toISOString() } : t,
        ),
      )
      return null
    } catch (err) {
      return getApiErrorMessage(err, 'No se pudo transferir la entrada')
    }
  }, [])

  return { tickets, loading, error, refetch, cancelTicket, transferTicket }
}
