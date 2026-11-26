import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { AxiosError } from 'axios'
import PaymentModal from './PaymentModal'

function renderModal(onConfirm: () => Promise<void>) {
  const onClose = vi.fn()
  render(
    <PaymentModal eventTitle="Recital de prueba" quantity={2} total={10000} onConfirm={onConfirm} onClose={onClose} />,
  )
  return { onClose }
}

describe('PaymentModal', () => {
  it('muestra el resumen de la compra', () => {
    renderModal(vi.fn().mockResolvedValue(undefined))

    expect(screen.getByText('Recital de prueba')).toBeInTheDocument()
    expect(screen.getByText('2 entradas')).toBeInTheDocument()
    expect(screen.getByText('$10000.00')).toBeInTheDocument()
  })

  it('confirma la compra y cierra el modal', async () => {
    const user = userEvent.setup()
    const onConfirm = vi.fn().mockResolvedValue(undefined)
    const { onClose } = renderModal(onConfirm)

    await user.click(screen.getByRole('button', { name: 'Confirmar compra' }))

    expect(onConfirm).toHaveBeenCalledTimes(1)
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('si la compra falla, muestra el motivo real y no cierra el modal', async () => {
    const user = userEvent.setup()
    const backendError = Object.assign(new AxiosError('Request failed'), {
      response: { data: { error: 'event 3: invalid presale code' } },
    })
    const { onClose } = renderModal(vi.fn().mockRejectedValue(backendError))

    await user.click(screen.getByRole('button', { name: 'Confirmar compra' }))

    expect(await screen.findByText('El código de preventa no es válido.')).toBeInTheDocument()
    expect(onClose).not.toHaveBeenCalled()
  })

  it('si el backend no responde, muestra el mensaje genérico', async () => {
    const user = userEvent.setup()
    renderModal(vi.fn().mockRejectedValue(new AxiosError('Network Error')))

    await user.click(screen.getByRole('button', { name: 'Confirmar compra' }))

    expect(await screen.findByText('Error al procesar la compra. Intentá de nuevo.')).toBeInTheDocument()
  })

  it('el botón Cancelar cierra el modal sin comprar', async () => {
    const user = userEvent.setup()
    const onConfirm = vi.fn()
    const { onClose } = renderModal(onConfirm)

    await user.click(screen.getByRole('button', { name: 'Cancelar' }))

    expect(onClose).toHaveBeenCalledTimes(1)
    expect(onConfirm).not.toHaveBeenCalled()
  })
})
