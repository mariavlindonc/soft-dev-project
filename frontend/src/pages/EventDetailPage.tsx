import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import type { Event, SaleStatus } from '../types'
import { getEventById, getSaleStatus } from '../api/events'
import { purchaseTicket } from '../api/tickets'
import { useAuth } from '../context/AuthContext'
import { formatDate, formatPrice } from '../utils/format'

export default function EventDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { isAuthenticated } = useAuth()

  const [event, setEvent] = useState<Event | null>(null)
  const [saleStatus, setSaleStatus] = useState<SaleStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [presaleCode, setPresaleCode] = useState('')
  const [purchasing, setPurchasing] = useState(false)
  const [purchaseError, setPurchaseError] = useState<string | null>(null)
  const [purchased, setPurchased] = useState(false)

  useEffect(() => {
    if (!id) return
    const numId = Number(id)
    setLoading(true)
    setError(null)
    Promise.all([getEventById(numId), getSaleStatus(numId)])
      .then(([ev, sale]) => {
        setEvent(ev)
        setSaleStatus(sale)
      })
      .catch(() => setError('No se pudo cargar el evento'))
      .finally(() => setLoading(false))
  }, [id])

  const isPresaleCodeRequired =
    saleStatus?.phase === 'presale' && saleStatus.message.toLowerCase().includes('code')

  async function handlePurchase() {
    if (!event) return
    setPurchaseError(null)
    setPurchasing(true)
    try {
      await purchaseTicket({
        event_id: event.id,
        presale_code: isPresaleCodeRequired ? presaleCode : undefined,
      })
      setPurchased(true)
    } catch (err: unknown) {
      const msg =
        err && typeof err === 'object' && 'response' in err
          ? (err as { response: { data: { error: string } } }).response?.data?.error
          : 'Error al comprar la entrada'
      setPurchaseError(msg ?? 'Error al comprar la entrada')
    } finally {
      setPurchasing(false)
    }
  }

  if (loading) {
    return (
      <div className="page">
        <div className="loading-container"><div className="spinner" /></div>
      </div>
    )
  }

  if (error || !event) {
    return (
      <div className="page">
        <div className="alert alert-error">{error ?? 'Evento no encontrado'}</div>
        <Link to="/events" className="btn btn-outline">Volver a eventos</Link>
      </div>
    )
  }

  const phaseLabel: Record<string, string> = {
    presale: 'Preventa',
    public: 'Venta general',
    not_yet_open: 'No disponible',
    no_presale: 'Venta general',
  }

  const phaseClass: Record<string, string> = {
    presale: 'presale',
    public: 'public',
    not_yet_open: 'closed',
    no_presale: 'public',
  }

  return (
    <div className="page">
      <Link to="/events" className="back-link">← Volver a eventos</Link>

      {purchased ? (
        <div className="purchase-success">
          <h3>¡Compra exitosa!</h3>
          <p>Tu entrada para <strong>{event.title}</strong> fue confirmada.</p>
          <div className="purchase-success-actions">
            <Link to="/tickets" className="btn btn-primary">Ver mis entradas</Link>
            <Link to="/events" className="btn btn-outline">Seguir explorando</Link>
          </div>
        </div>
      ) : (
        <div className="event-detail">
          <div
            className="event-detail-image"
            style={{ backgroundImage: event.image_url ? `url(${event.image_url})` : undefined }}
          />

          <div className="event-detail-info">
            <h1>{event.title}</h1>
            <span className="event-category">{event.category ?? 'General'}</span>

            <div className="event-detail-meta">
              <div className="event-detail-meta-item">
                <span>📅</span>
                <span>{formatDate(event.event_date)}</span>
              </div>
              {event.location && (
                <div className="event-detail-meta-item">
                  <span>📍</span>
                  <span>{event.location}</span>
                </div>
              )}
              <div className="event-detail-meta-item">
                <span>⏱</span>
                <span>{event.duration_minutes} min</span>
              </div>
            </div>

            {event.description && (
              <p className="event-detail-description">{event.description}</p>
            )}

            <div className="event-detail-price">{formatPrice(event.price)}</div>

            <div className="event-detail-capacity">
              Capacidad: <span>{event.tickets_sold}/{event.capacity}</span>
              {event.tickets_sold >= event.capacity && (
                <span className="sold-out-label">Agotado</span>
              )}
            </div>

            {event.tickets_sold < event.capacity && (
              <div className="purchase-section">
                {saleStatus && (
                  <span className={`sale-status ${phaseClass[saleStatus.phase] ?? 'closed'}`}>
                    {phaseLabel[saleStatus.phase] ?? saleStatus.phase}
                  </span>
                )}

                {saleStatus && (
                  <p className="alert alert-info">{saleStatus.message}</p>
                )}

                {saleStatus?.phase === 'presale' && isPresaleCodeRequired && (
                  <div className="form-group">
                    <label htmlFor="presaleCode">Código de preventa</label>
                    <input
                      id="presaleCode"
                      type="text"
                      value={presaleCode}
                      onChange={(e) => setPresaleCode(e.target.value)}
                      placeholder="Ingresá tu código"
                    />
                  </div>
                )}

                {purchaseError && <div className="form-global-error">{purchaseError}</div>}

                {saleStatus && saleStatus.phase !== 'not_yet_open' && (
                  <div className="purchase-actions">
                    {isAuthenticated ? (
                      <button
                        type="button"
                        className="btn btn-primary btn-lg"
                        onClick={handlePurchase}
                        disabled={purchasing || (isPresaleCodeRequired && !presaleCode.trim())}
                      >
                        {purchasing ? 'Comprando…' : 'Comprar entrada'}
                      </button>
                    ) : (
                      <div className="alert alert-info">
                        <Link to="/login" className="auth-link">Iniciá sesión</Link> para comprar entradas.
                      </div>
                    )}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
