import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useTickets } from '../hooks/useTickets'
import { formatDate, formatPrice } from '../utils/format'

export default function TicketsPage() {
  const { isAuthenticated } = useAuth()
  const { tickets, loading, error, cancelTicket, transferTicket, refetch } = useTickets()

  const [transferModal, setTransferModal] = useState<{ ticketId: number; eventTitle: string } | null>(null)
  const [transferEmail, setTransferEmail] = useState('')
  const [transferring, setTransferring] = useState(false)
  const [transferError, setTransferError] = useState<string | null>(null)
  const [cancellingId, setCancellingId] = useState<number | null>(null)

  if (!isAuthenticated) {
    return (
      <div className="page">
        <div className="tickets-header"><h1>Mis Entradas</h1></div>
        <div className="alert alert-info">
          <Link to="/login" className="auth-link">Iniciá sesión</Link> para ver tus entradas.
        </div>
      </div>
    )
  }

  async function handleCancel(id: number) {
    setCancellingId(id)
    const ok = await cancelTicket(id)
    setCancellingId(null)
    if (!ok) alert('No se pudo cancelar la entrada')
  }

  function openTransfer(ticketId: number, eventTitle: string) {
    setTransferModal({ ticketId, eventTitle })
    setTransferEmail('')
    setTransferError(null)
  }

  async function handleTransfer() {
    if (!transferModal) return
    if (!transferEmail.trim() || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(transferEmail)) {
      setTransferError('Ingresá un correo válido')
      return
    }
    setTransferring(true)
    setTransferError(null)
    try {
      const ok = await transferTicket(transferModal.ticketId, transferEmail)
      if (ok) {
        setTransferModal(null)
        refetch()
      } else {
        setTransferError('No se pudo transferir la entrada')
      }
    } catch {
      setTransferError('Error al transferir la entrada')
    } finally {
      setTransferring(false)
    }
  }

  const statusLabel: Record<string, string> = {
    active: 'Activa',
    cancelled: 'Cancelada',
    transferred: 'Transferida',
  }

  const statusClass: Record<string, string> = {
    active: 'ticket-status-active',
    cancelled: 'ticket-status-cancelled',
    transferred: 'ticket-status-transferred',
  }

  return (
    <div className="page">
      <div className="tickets-header">
        <h1>Mis Entradas</h1>
      </div>

      {loading && <div className="loading-text">Cargando tus entradas…</div>}

      {error && <div className="alert alert-error">{error}</div>}

      {!loading && !error && tickets.length === 0 && (
        <div className="tickets-empty">
          <p>No tenés entradas</p>
          <p className="form-footer">
            <Link to="/events">Explorá eventos</Link> y comprá tus primeras entradas.
          </p>
        </div>
      )}

      {!loading && !error && tickets.length > 0 && (
        <div className="ticket-list">
          {tickets.map((ticket) => (
            <div key={ticket.id} className="ticket-card">
              <div className="ticket-card-info">
                <h3>{ticket.event_title || `Evento #${ticket.event_id}`}</h3>
                {ticket.event_date && (
                  <p>{formatDate(ticket.event_date)}</p>
                )}
                <p>Comprado el {formatDate(ticket.purchased_at)} — {formatPrice(ticket.purchase_price)}</p>
                <span className={`ticket-card-status ${statusClass[ticket.status]}`}>
                  {statusLabel[ticket.status]}
                </span>
              </div>
              <div className="ticket-card-actions">
                {ticket.status === 'active' && (
                  <>
                    <button
                      type="button"
                      className="btn btn-outline btn-sm"
                      onClick={() => openTransfer(ticket.id, ticket.event_title || `Evento #${ticket.event_id}`)}
                    >
                      Transferir
                    </button>
                    <button
                      type="button"
                      className="btn btn-danger btn-sm"
                      onClick={() => handleCancel(ticket.id)}
                      disabled={cancellingId === ticket.id}
                    >
                      {cancellingId === ticket.id ? 'Cancelando…' : 'Cancelar'}
                    </button>
                  </>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {transferModal && (
        <div className="modal-overlay" onClick={() => setTransferModal(null)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>Transferir entrada</h3>
            <p>Transferí tu entrada para <strong>{transferModal.eventTitle}</strong> a otro usuario.</p>

            {transferError && <div className="form-global-error">{transferError}</div>}

            <div className="form-group">
              <label htmlFor="transferEmail">Correo del destinatario</label>
              <input
                id="transferEmail"
                type="email"
                value={transferEmail}
                onChange={(e) => setTransferEmail(e.target.value)}
                placeholder="correo@ejemplo.com"
              />
            </div>

            <div className="modal-actions">
              <button type="button" className="btn btn-outline" onClick={() => setTransferModal(null)}>
                Cancelar
              </button>
              <button type="button" className="btn btn-primary" onClick={handleTransfer} disabled={transferring}>
                {transferring ? 'Transfiriendo…' : 'Transferir'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
