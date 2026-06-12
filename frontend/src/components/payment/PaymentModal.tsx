import { useState, useRef, useEffect } from 'react'

interface PaymentModalProps {
  eventTitle: string
  quantity: number
  total: number
  onConfirm: () => Promise<void>
  onClose: () => void
}

type CardBrand = 'visa' | 'mastercard' | 'amex' | 'unknown'

function detectCardBrand(number: string): CardBrand {
  const clean = number.replace(/\s/g, '')
  if (/^4/.test(clean)) return 'visa'
  if (/^5[1-5]/.test(clean)) return 'mastercard'
  if (/^3[47]/.test(clean)) return 'amex'
  return 'unknown'
}

function formatCardNumber(value: string): string {
  const digits = value.replace(/\D/g, '')
  const brand = detectCardBrand(digits)
  const groups: number[] = brand === 'amex' ? [4, 6, 5] : [4, 4, 4, 4]
  let result = ''
  let i = 0
  for (const g of groups) {
    if (i >= digits.length) break
    result += digits.slice(i, i + g) + ' '
    i += g
  }
  return result.trim()
}

function formatExpiry(value: string): string {
  const digits = value.replace(/\D/g, '')
  if (digits.length >= 2) {
    return digits.slice(0, 2) + '/' + digits.slice(2, 4)
  }
  return digits
}

const brandLogos: Record<CardBrand, string> = {
  visa: 'VISA',
  mastercard: 'MC',
  amex: 'AMEX',
  unknown: '💳',
}

export default function PaymentModal({ eventTitle, quantity, total, onConfirm, onClose }: PaymentModalProps) {
  const [step, setStep] = useState<'form' | 'processing' | 'error'>('form')
  const [cardNumber, setCardNumber] = useState('')
  const [cardName, setCardName] = useState('')
  const [expiry, setExpiry] = useState('')
  const [cvv, setCvv] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [saveCard, setSaveCard] = useState(false)
  const [acceptedTerms, setAcceptedTerms] = useState(false)

  const numberRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    numberRef.current?.focus()
  }, [])

  const brand = detectCardBrand(cardNumber)

  function validate(): string | null {
    const cleanNumber = cardNumber.replace(/\s/g, '')
    if (cleanNumber.length < 13 || cleanNumber.length > 19) {
      return 'Número de tarjeta inválido'
    }
    if (cardName.trim().length < 3) {
      return 'Ingresá el nombre del titular'
    }
    const expiryDigits = expiry.replace(/\D/g, '')
    if (expiryDigits.length !== 4) {
      return 'Fecha de vencimiento inválida'
    }
    const month = parseInt(expiryDigits.slice(0, 2), 10)
    const year = parseInt(expiryDigits.slice(2, 4), 10) + 2000
    if (month < 1 || month > 12) {
      return 'Mes de vencimiento inválido'
    }
    const now = new Date()
    const expiryDate = new Date(year, month, 0)
    if (expiryDate < now) {
      return 'La tarjeta está vencida'
    }
    if (cvv.length < 3 || cvv.length > 4) {
      return 'Código de seguridad inválido'
    }
    if (!acceptedTerms) {
      return 'Debés aceptar los términos y condiciones'
    }
    return null
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const err = validate()
    if (err) {
      setError(err)
      return
    }

    setStep('processing')
    setError(null)

    // Simulate payment processing delay
    await new Promise((r) => setTimeout(r, 2000))

    try {
      await onConfirm()
      onClose()
    } catch {
      setStep('error')
      setError('Error al procesar el pago. Intentá de nuevo.')
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content payment-modal" onClick={(e) => e.stopPropagation()}>
        {step === 'processing' ? (
          <div className="payment-processing">
            <div className="spinner" />
            <h3>Procesando pago</h3>
            <p className="payment-processing-text">No cierres esta ventana</p>
          </div>
        ) : (
          <>
            <div className="payment-header">
              <h3>Pagar con tarjeta</h3>
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

            <div className="payment-brand-indicator">
              <span className={`payment-brand-logo payment-brand-${brand}`}>
                {brandLogos[brand]}
              </span>
            </div>

            <form className="payment-form" onSubmit={handleSubmit} noValidate>
              {error && <div className="form-global-error">{error}</div>}

              <div className="form-group">
                <label htmlFor="cardNumber">Número de tarjeta</label>
                <input
                  ref={numberRef}
                  id="cardNumber"
                  type="text"
                  inputMode="numeric"
                  maxLength={23}
                  placeholder="1234 5678 9012 3456"
                  value={cardNumber}
                  onChange={(e) => setCardNumber(formatCardNumber(e.target.value))}
                />
              </div>

              <div className="form-group">
                <label htmlFor="cardName">Titular</label>
                <input
                  id="cardName"
                  type="text"
                  placeholder="Como figura en la tarjeta"
                  value={cardName}
                  onChange={(e) => setCardName(e.target.value)}
                  autoComplete="cc-name"
                />
              </div>

              <div className="payment-form-row">
                <div className="form-group">
                  <label htmlFor="expiry">Vencimiento</label>
                  <input
                    id="expiry"
                    type="text"
                    inputMode="numeric"
                    maxLength={5}
                    placeholder="MM/AA"
                    value={expiry}
                    onChange={(e) => setExpiry(formatExpiry(e.target.value))}
                    autoComplete="cc-exp"
                  />
                </div>
                <div className="form-group">
                  <label htmlFor="cvv">CVV</label>
                  <input
                    id="cvv"
                    type="text"
                    inputMode="numeric"
                    maxLength={4}
                    placeholder="123"
                    value={cvv}
                    onChange={(e) => setCvv(e.target.value.replace(/\D/g, ''))}
                    autoComplete="cc-csc"
                  />
                </div>
              </div>

              <div className="payment-checkbox-row">
                <input
                  id="saveCard"
                  type="checkbox"
                  checked={saveCard}
                  onChange={(e) => setSaveCard(e.target.checked)}
                />
                <label htmlFor="saveCard">Guardar tarjeta para próximas compras</label>
              </div>

              <div className="payment-checkbox-row">
                <input
                  id="acceptedTerms"
                  type="checkbox"
                  checked={acceptedTerms}
                  onChange={(e) => setAcceptedTerms(e.target.checked)}
                />
                <label htmlFor="acceptedTerms">
                  Acepto los <a href="/terms" target="_blank">términos y condiciones</a>
                </label>
              </div>

              <div className="payment-actions">
                <button type="button" className="btn btn-outline" onClick={onClose}>
                  Cancelar
                </button>
                <button type="submit" className="btn btn-primary btn-lg">
                  Pagar ${total.toFixed(2)}
                </button>
              </div>
            </form>

            <p className="payment-footer-text">
              Pago seguro cifrado con SSL
            </p>
          </>
        )}
      </div>
    </div>
  )
}
