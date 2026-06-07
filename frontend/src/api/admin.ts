import client from './client'
import type { Event } from '../types/event'
import type { CreateEventData, UpdateEventData, GlobalReport, EventReport } from '../types/admin'

export async function createEvent(data: CreateEventData): Promise<Event> {
  const res = await client.post('/admin/events', data)
  return res.data
}

export async function updateEvent(id: number, data: UpdateEventData): Promise<Event> {
  const res = await client.put(`/admin/events/${id}`, data)
  return res.data
}

export async function deleteEvent(id: number): Promise<void> {
  await client.delete(`/admin/events/${id}`)
}

export async function getGlobalReport(): Promise<GlobalReport> {
  const res = await client.get('/admin/reports')
  return res.data
}

export async function getEventReport(id: number): Promise<EventReport> {
  const res = await client.get(`/admin/reports/events/${id}`)
  return res.data
}
