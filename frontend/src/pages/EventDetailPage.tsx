import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import type { Event, SaleStatus } from '../types'
import { getEventById, getSaleStatus } from '../api/events'
import { purchaseTicket } from '../api/tickets'
import { mockEvents } from '../data/mockEvents'
import { useAuth } from '../context/AuthContext'
import { formatDateTime, formatPrice } from '../utils/format'
import PaymentModal from '../components/payment/PaymentModal'

function computeSaleStatus(event: Event): SaleStatus {
  if (event.status === 'cancelled') {
    return { phase: 'not_yet_open', presale_start_date: null, general_sale_date: null, message: 'Evento cancelado' }
  }
  if (event.tickets_sold >= event.capacity) {
    return { phase: 'not_yet_open', presale_start_date: null, general_sale_date: null, message: 'Entradas agotadas' }
  }
  if (event.presale_active && event.presale_start_date && new Date(event.presale_start_date) <= new Date()) {
    if (event.general_sale_date && new Date(event.general_sale_date) <= new Date()) {
      return { phase: 'public', presale_start_date: event.presale_start_date, general_sale_date: event.general_sale_date, message: 'Venta general disponible' }
    }
    return { phase: 'presale', presale_start_date: event.presale_start_date, general_sale_date: event.general_sale_date, message: 'Preventa disponible — ingresá tu código' }
  }
  if (event.presale_active && event.presale_start_date && new Date(event.presale_start_date) > new Date()) {
    return { phase: 'not_yet_open', presale_start_date: event.presale_start_date, general_sale_date: event.general_sale_date, message: 'Preventa próximamente' }
  }
  return { phase: 'public', presale_start_date: null, general_sale_date: null, message: 'Entradas disponibles' }
}

export default function EventDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { isAuthenticated, login, register } = useAuth()

  const [event, setEvent] = useState<Event | null>(null)
  const [saleStatus, setSaleStatus] = useState<SaleStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [quantity, setQuantity] = useState(1)
  const [presaleCode, setPresaleCode] = useState('')
  const [purchasing, setPurchasing] = useState(false)
  const [purchaseError, setPurchaseError] = useState<string | null>(null)
  const [purchased, setPurchased] = useState(false)
  const [showPaymentModal, setShowPaymentModal] = useState(false)

  const [authName, setAuthName] = useState('')
  const [authEmail, setAuthEmail] = useState('')
  const [authPassword, setAuthPassword] = useState('')
  const [isRegistering, setIsRegistering] = useState(false)

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
      .catch(() => {
        const mock = mockEvents.find((e) => e.id === numId)
        if (mock) {
          setEvent(mock)
          setSaleStatus(computeSaleStatus(mock))
        } else {
          setError('Evento no encontrado')
        }
      })
      .finally(() => setLoading(false))
  }, [id])

  const isPresaleCodeRequired =
    saleStatus?.phase === 'presale' && saleStatus.message.toLowerCase().includes('code')

  const available = event ? event.capacity - event.tickets_sold : 0

  async function handlePurchase() {
    if (!event) return
    setPurchaseError(null)

    if (!isAuthenticated) {
      setPurchasing(true)
      try {
        if (isRegistering) {
          await register({ name: authName, email: authEmail, password: authPassword })
        } else {
          await login({ email: authEmail, password: authPassword })
        }
      } catch {
        setPurchaseError('Error al iniciar sesión. Verificá tus datos.')
        setPurchasing(false)
        return
      }
      setPurchasing(false)
    }

    setShowPaymentModal(true)
  }

  async function handlePaymentConfirm() {
    if (!event) return
    await purchaseTicket({
      event_id: event.id,
      presale_code: isPresaleCodeRequired ? presaleCode : undefined,
    })
    setPurchased(true)
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
                <span>{formatDateTime(event.event_date)}</span>
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

                <div className="purchase-quantity">
                  <span className="purchase-quantity-label">Cantidad</span>
                  <div className="purchase-quantity-controls">
                    <button
                      type="button"
                      className="purchase-quantity-btn"
                      onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                      disabled={quantity <= 1}
                    >
                      −
                    </button>
                    <span className="purchase-quantity-value">{quantity}</span>
                    <button
                      type="button"
                      className="purchase-quantity-btn"
                      onClick={() => setQuantity((q) => Math.min(available, q + 1))}
                      disabled={quantity >= available}
                    >
                      +
                    </button>
                  </div>
                  <span className="purchase-total">{formatPrice(event.price * quantity)}</span>
                </div>

                {!isAuthenticated && (
                  <div className="purchase-auth-form">
                    <p className="purchase-auth-title">{isRegistering ? 'Crear cuenta' : 'Iniciar sesión'}</p>
                    {isRegistering && (
                      <div className="form-group">
                        <label htmlFor="authName">Nombre</label>
                        <input id="authName" type="text" value={authName} onChange={(e) => setAuthName(e.target.value)} placeholder="Tu nombre" />
                      </div>
                    )}
                    <div className="form-group">
                      <label htmlFor="authEmail">Email</label>
                      <input id="authEmail" type="email" value={authEmail} onChange={(e) => setAuthEmail(e.target.value)} placeholder="correo@ejemplo.com" />
                    </div>
                    <div className="form-group">
                      <label htmlFor="authPassword">Contraseña</label>
                      <input id="authPassword" type="password" value={authPassword} onChange={(e) => setAuthPassword(e.target.value)} placeholder="••••••••" />
                    </div>
                    <button type="button" className="btn btn-outline btn-sm" onClick={() => { setIsRegistering(!isRegistering); setPurchaseError(null) }}>
                      {isRegistering ? 'Ya tengo cuenta' : 'Crear cuenta nueva'}
                    </button>
                  </div>
                )}

                {purchaseError && <div className="form-global-error">{purchaseError}</div>}

                <div className="purchase-actions">
                  <button
                    type="button"
                    className="btn btn-primary btn-lg"
                    onClick={handlePurchase}
                    disabled={purchasing || (isPresaleCodeRequired && !presaleCode.trim()) || (!isAuthenticated && !authEmail.trim() || !authPassword.trim())}
                  >
                    {purchasing ? 'Comprando…' : `Comprar ${quantity > 1 ? `${quantity} entradas` : 'entrada'}`}
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {showPaymentModal && event && (
        <PaymentModal
          eventTitle={event.title}
          quantity={quantity}
          total={event.price * quantity}
          onConfirm={handlePaymentConfirm}
          onClose={() => setShowPaymentModal(false)}
        />
      )}
    </div>
  )
}
