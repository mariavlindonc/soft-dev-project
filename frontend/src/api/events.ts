import client from './client'
import type { Event, EventFilters, SaleStatus } from '../types/event'

export async function getEvents(filters?: EventFilters): Promise<Event[]> {
  const res = await client.get('/events', { params: filters })
  return res.data
}

export async function getEventById(id: number): Promise<Event> {
  const res = await client.get(`/events/${id}`)
  return res.data
}

export async function getSaleStatus(id: number): Promise<SaleStatus> {
  const res = await client.get(`/events/${id}/sale-status`)
  return res.data
}
