import { useState } from 'react'

interface PaymentModalProps {
  eventTitle: string
  quantity: number
  total: number
  onConfirm: () => Promise<void>
  onClose: () => void
}

export default function PaymentModal({ eventTitle, quantity, total, onConfirm, onClose }: PaymentModalProps) {
  const [confirming, setConfirming] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleConfirm() {
    setConfirming(true)
    setError(null)
    try {
      await onConfirm()
      onClose()
    } catch {
      setError('Error al procesar la compra. Intentá de nuevo.')
      setConfirming(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content payment-modal" onClick={(e) => e.stopPropagation()}>
        {confirming ? (
          <div className="payment-processing">
            <div className="spinner" />
            <h3>Procesando compra</h3>
            <p className="payment-processing-text">No cierres esta ventana</p>
          </div>
        ) : (
          <>
            <div className="payment-header">
              <h3>Confirmar compra</h3>
              <button type="button" className="payment-close" onClick={onClose}>✕</button>
            </div>

            <div className="payment-summary">
              <div className="payment-summary-row">
                <span className="payment-summary-label">Evento</span>
                <span className="payment-summary-value">{eventTitle}</span>
              </div>
              <div className="payment-summary-row">
                <span className="payment-summary-label">Entradas</span>
                <span className="payment-summary-value">{quantity} {quantity === 1 ? 'entrada' : 'entradas'}</span>
              </div>
              <div className="payment-summary-row payment-summary-total">
                <span className="payment-summary-label">Total</span>
                <span className="payment-summary-value">${total.toFixed(2)}</span>
              </div>
            </div>

            {error && <div className="form-global-error">{error}</div>}

            <div className="payment-actions">
              <button type="button" className="btn btn-outline" onClick={onClose}>
                Cancelar
              </button>
              <button type="button" className="btn btn-primary btn-lg" onClick={handleConfirm}>
                Confirmar compra
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
