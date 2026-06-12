import { useEffect, useState } from 'react'
import { Routes, Route, Link, useNavigate, useLocation, useParams } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useEvents } from '../hooks/useEvents'
import * as adminApi from '../api/admin'
import { getEventById } from '../api/events'
import { formatDate, formatPrice } from '../utils/format'
import type { GlobalReport, EventReport, CreateEventData, UpdateEventData } from '../types/admin'
import type { Event } from '../types'

function AdminNav() {
  const { pathname } = useLocation()
  const isActive = (path: string) => pathname === path || pathname.startsWith(path + '/') ? 'active' : ''

  return (
    <aside className="admin-sidebar">
      <h3>Panel Admin</h3>
      <nav className="admin-nav">
        <Link to="/admin" className={isActive('/admin') === 'active' && pathname === '/admin' ? 'active' : ''}>
          Dashboard
        </Link>
        <Link to="/admin/events" className={isActive('/admin/events')}>
          Gestionar Eventos
        </Link>
      </nav>
    </aside>
  )
}

function Dashboard() {
  const [report, setReport] = useState<GlobalReport | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    adminApi.getGlobalReport()
      .then(setReport)
      .catch(() => setError('Error al cargar el reporte'))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="loading-container"><div className="spinner" /></div>
  if (error) return <div className="alert alert-error">{error}</div>
  if (!report) return null

  return (
    <div>
      <div className="admin-header">
        <h1>Dashboard</h1>
      </div>

      <div className="admin-stats-grid">
        <div className="admin-stat-card">
          <h4>Total Eventos</h4>
          <p>{report.total_events}</p>
        </div>
        <div className="admin-stat-card">
          <h4>Entradas Vendidas</h4>
          <p>{report.total_tickets_sold}</p>
        </div>
      </div>

      <h2 className="admin-section-title">Eventos</h2>
      <table className="admin-table">
        <thead>
          <tr>
            <th>Evento</th>
            <th>Capacidad</th>
            <th>Vendidas</th>
            <th>Ocupación</th>
          </tr>
        </thead>
        <tbody>
          {report.events.map((ev) => (
            <tr key={ev.event_id}>
              <td>
                <Link to={`/admin/reports/${ev.event_id}`} style={{ color: 'var(--color-primary)', fontWeight: 500 }}>
                  {ev.title}
                </Link>
              </td>
              <td>{ev.capacity}</td>
              <td>{ev.tickets_sold}</td>
              <td>{ev.occupancy.toFixed(1)}%</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function AdminEventList() {
  const navigate = useNavigate()
  const { events, loading, error, refetch } = useEvents()

  async function handleDelete(id: number) {
    if (!window.confirm('¿Cancelar este evento? Se cancelarán todas las entradas activas.')) return
    try {
      await adminApi.deleteEvent(id)
      refetch()
    } catch {
      alert('Error al cancelar el evento')
    }
  }

  if (loading) return <div className="loading-container"><div className="spinner" /></div>

  return (
    <div>
      <div className="admin-header">
        <h1>Gestionar Eventos</h1>
        <Link to="/admin/events/new" className="btn btn-primary">Crear Evento</Link>
      </div>

      {error && <div className="alert alert-error">{error}</div>}

      {events.length === 0 ? (
        <p style={{ color: 'var(--color-text-secondary)', textAlign: 'center', padding: '2rem' }}>
          No hay eventos creados todavía.
        </p>
      ) : (
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Título</th>
            <th>Categoría</th>
            <th>Fecha</th>
            <th>Precio</th>
            <th>Estado</th>
            <th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {events.map((event) => (
            <tr key={event.id}>
              <td>{event.id}</td>
              <td>{event.title}</td>
              <td>{event.category ?? '-'}</td>
              <td>{formatDate(event.event_date)}</td>
              <td>{formatPrice(event.price)}</td>
              <td>
                <span className={`sale-status ${event.status}`}>
                  {event.status === 'active' ? 'Activo' : event.status === 'presale' ? 'Preventa' : event.status === 'sold_out' ? 'Agotado' : event.status === 'cancelled' ? 'Cancelado' : event.status}
                </span>
              </td>
              <td>
                <div className="admin-table-actions">
                  <button type="button" className="btn btn-outline btn-sm" onClick={() => navigate(`/admin/events/${event.id}/edit`)}>
                    Editar
                  </button>
                  <Link to={`/admin/reports/${event.id}`} className="btn btn-outline btn-sm">Reporte</Link>
                  {event.status !== 'cancelled' && (
                    <button type="button" className="btn btn-danger btn-sm" onClick={() => handleDelete(event.id)}>
                      Cancelar
                    </button>
                  )}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      )}
    </div>
  )
}

function AdminEventForm() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isEdit = !!id

  const [form, setForm] = useState<CreateEventData>({
    title: '',
    description: '',
    image_url: '',
    category: '',
    location: '',
    event_date: '',
    duration_minutes: 60,
    capacity: 100,
    price: 0,
    presale_active: false,
    presale_code: '',
    presale_start_date: '',
    general_sale_date: '',
  })
  const [loading, setLoading] = useState(isEdit)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    getEventById(Number(id))
      .then((event: Event) => {
        setForm({
          title: event.title,
          description: event.description ?? '',
          image_url: event.image_url ?? '',
          category: event.category ?? '',
          location: event.location ?? '',
          event_date: event.event_date.slice(0, 16),
          duration_minutes: event.duration_minutes,
          capacity: event.capacity,
          price: Number(event.price),
          presale_active: event.presale_active,
          presale_code: '',
          presale_start_date: event.presale_start_date?.slice(0, 16) ?? '',
          general_sale_date: event.general_sale_date?.slice(0, 16) ?? '',
        })
        setLoading(false)
      })
      .catch(() => {
        setError('Error al cargar el evento')
        setLoading(false)
      })
  }, [id])

  function handleChange(field: string, value: string | number | boolean) {
    setForm((prev) => ({ ...prev, [field]: value }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSaving(true)
    setError(null)
    try {
      const payload = {
        ...form,
        event_date: form.event_date ? new Date(form.event_date).toISOString() : '',
        presale_start_date: form.presale_start_date ? new Date(form.presale_start_date).toISOString() : undefined,
        general_sale_date: form.general_sale_date ? new Date(form.general_sale_date).toISOString() : undefined,
        presale_code: form.presale_code || undefined,
      }
      if (isEdit) {
        const updateData: UpdateEventData = { ...payload }
        await adminApi.updateEvent(Number(id), updateData)
      } else {
        await adminApi.createEvent(payload as CreateEventData)
      }
      navigate('/admin/events')
    } catch (err: unknown) {
      const msg =
        err && typeof err === 'object' && 'response' in err
          ? (err as { response: { data: { error: string } } }).response?.data?.error
          : 'Error al guardar el evento'
      setError(msg ?? 'Error al guardar el evento')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="loading-container"><div className="spinner" /></div>

  return (
    <div>
      <div className="admin-header">
        <h1>{isEdit ? 'Editar Evento' : 'Crear Evento'}</h1>
      </div>

      {error && <div className="form-global-error">{error}</div>}

      <form className="admin-form" onSubmit={handleSubmit} noValidate>
        <div className="form-group">
          <label htmlFor="title">Título *</label>
          <input id="title" required value={form.title} onChange={(e) => handleChange('title', e.target.value)} />
        </div>

        <div className="form-group">
          <label htmlFor="description">Descripción</label>
          <textarea id="description" value={form.description ?? ''} onChange={(e) => handleChange('description', e.target.value)} />
        </div>

        <div className="admin-form-row">
          <div className="form-group">
            <label htmlFor="category">Categoría</label>
            <input id="category" value={form.category ?? ''} onChange={(e) => handleChange('category', e.target.value)} />
          </div>
          <div className="form-group">
            <label htmlFor="location">Ubicación</label>
            <input id="location" value={form.location ?? ''} onChange={(e) => handleChange('location', e.target.value)} />
          </div>
        </div>

        <div className="admin-form-row">
          <div className="form-group">
            <label htmlFor="event_date">Fecha del evento *</label>
            <input id="event_date" type="datetime-local" required value={form.event_date} onChange={(e) => handleChange('event_date', e.target.value)} />
          </div>
          <div className="form-group">
            <label htmlFor="duration_minutes">Duración (min)</label>
            <input id="duration_minutes" type="number" min={1} value={form.duration_minutes ?? ''} onChange={(e) => handleChange('duration_minutes', Number(e.target.value))} />
          </div>
        </div>

        <div className="admin-form-row">
          <div className="form-group">
            <label htmlFor="capacity">Capacidad *</label>
            <input id="capacity" type="number" required min={1} value={form.capacity} onChange={(e) => handleChange('capacity', Number(e.target.value))} />
          </div>
          <div className="form-group">
            <label htmlFor="price">Precio *</label>
            <input id="price" type="number" required min={0} step="0.01" value={form.price} onChange={(e) => handleChange('price', Number(e.target.value))} />
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="image_url">URL de imagen</label>
          <input id="image_url" value={form.image_url ?? ''} onChange={(e) => handleChange('image_url', e.target.value)} />
        </div>

        <h2 className="admin-section-title">Preventa</h2>

        <div className="form-group" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
          <input id="presale_active" type="checkbox" checked={form.presale_active} onChange={(e) => handleChange('presale_active', e.target.checked)} style={{ width: 'auto' }} />
          <label htmlFor="presale_active" style={{ margin: 0 }}>Activar preventa</label>
        </div>

        {form.presale_active && (
          <>
            <div className="admin-form-row">
              <div className="form-group">
                <label htmlFor="presale_code">Código de preventa</label>
                <input id="presale_code" value={form.presale_code ?? ''} onChange={(e) => handleChange('presale_code', e.target.value)} />
              </div>
              <div className="form-group">
                <label htmlFor="presale_start_date">Inicio preventa</label>
                <input id="presale_start_date" type="datetime-local" value={form.presale_start_date ?? ''} onChange={(e) => handleChange('presale_start_date', e.target.value)} />
              </div>
            </div>
            <div className="form-group">
              <label htmlFor="general_sale_date">Inicio venta general</label>
              <input id="general_sale_date" type="datetime-local" value={form.general_sale_date ?? ''} onChange={(e) => handleChange('general_sale_date', e.target.value)} />
            </div>
          </>
        )}

        <div style={{ display: 'flex', gap: '0.75rem', marginTop: '1.5rem' }}>
          <button type="submit" className="btn btn-primary btn-lg" disabled={saving}>
            {saving ? 'Guardando…' : isEdit ? 'Actualizar Evento' : 'Crear Evento'}
          </button>
          <button type="button" className="btn btn-outline btn-lg" onClick={() => navigate('/admin/events')}>
            Cancelar
          </button>
        </div>
      </form>
    </div>
  )
}

function AdminEventReport() {
  const { id } = useParams<{ id: string }>()
  const [report, setReport] = useState<EventReport | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    adminApi.getEventReport(Number(id))
      .then(setReport)
      .catch(() => setError('Error al cargar el reporte'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) return <div className="loading-container"><div className="spinner" /></div>
  if (error) return <div className="alert alert-error">{error}</div>
  if (!report) return null

  return (
    <div>
      <Link to="/admin" className="back-link">← Volver al dashboard</Link>
      <div className="admin-header">
        <h1>{report.title}</h1>
      </div>

      <div className="admin-stats-grid">
        <div className="admin-stat-card">
          <h4>Capacidad</h4>
          <p>{report.capacity}</p>
        </div>
        <div className="admin-stat-card">
          <h4>Entradas Vendidas</h4>
          <p>{report.tickets_sold}</p>
        </div>
        <div className="admin-stat-card">
          <h4>Ocupación</h4>
          <p>{report.occupancy.toFixed(1)}%</p>
        </div>
      </div>

      <h2 className="admin-section-title">Compradores ({report.buyers.length})</h2>
      <table className="admin-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Nombre</th>
            <th>Email</th>
          </tr>
        </thead>
        <tbody>
          {report.buyers.map((buyer) => (
            <tr key={buyer.user_id}>
              <td>{buyer.user_id}</td>
              <td>{buyer.name}</td>
              <td>{buyer.email}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function AdminPanel() {
  const { isAuthenticated, isAdmin } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login')
    } else if (!isAdmin) {
      navigate('/')
    }
  }, [isAuthenticated, isAdmin, navigate])

  if (!isAuthenticated || !isAdmin) return null

  return (
    <div className="page">
      <div className="admin-layout">
        <AdminNav />
        <main>
          <Routes>
            <Route index element={<Dashboard />} />
            <Route path="events" element={<AdminEventList />} />
            <Route path="events/new" element={<AdminEventForm />} />
            <Route path="events/:id/edit" element={<AdminEventForm />} />
            <Route path="reports/:id" element={<AdminEventReport />} />
            <Route path="*" element={<Dashboard />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}
