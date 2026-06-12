import client from './client'
import type { Ticket, PurchaseRequest } from '../types/ticket'

export async function purchaseTicket(data: PurchaseRequest): Promise<Ticket> {
  const res = await client.post('/tickets/purchase', data)
  return res.data
}

export async function getMyTickets(): Promise<Ticket[]> {
  const res = await client.get('/tickets')
  return res.data
}

export async function cancelTicket(id: number): Promise<void> {
  await client.patch(`/tickets/${id}/cancel`)
}

export async function transferTicket(id: number, toEmail: string): Promise<void> {
  await client.patch(`/tickets/${id}/transfer`, { to_user_email: toEmail })
}
