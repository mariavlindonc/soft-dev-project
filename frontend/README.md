# Ceibo Tickets — Frontend

SPA para la gestion de eventos y venta de entradas. Consume la API REST del backend.

## Stack

React 19 + TypeScript 6 + Vite 8 + React Router DOM v7 + Axios

## Inicio rapido

```bash
npm install
npm run dev       # http://localhost:5173
```

## Estructura

```
src/
├── api/          # Cliente Axios con interceptors (JWT y errores)
├── components/   # Componentes reutilizables (Layout, Navbar, EventCard, PaymentModal, etc.)
├── context/      # AuthContext (estado global de autenticacion + localStorage)
├── hooks/        # useEvents, useTickets (logica de negocio)
├── pages/        # Vistas completas (Home, Events, EventDetail, Auth, Tickets, AdminPanel, static pages)
├── types/        # Interfaces TypeScript (Event, Ticket, User, etc.)
├── data/         # Datos mock para desarrollo sin backend
└── utils/        # formatDate, formatPrice
```

## Paginas

| Ruta | Vista |
|------|-------|
| `/` | Home |
| `/events` | Listado de eventos con filtros |
| `/events/:id` | Detalle de evento + compra |
| `/login`, `/register` | Autenticacion |
| `/tickets` | Mis Entradas (cancelar, transferir) |
| `/admin/*` | Panel admin (CRUD eventos, reportes) |
| `/faq`, `/terms`, `/privacy` | Paginas estaticas |

## Comandos utiles

```bash
npm run dev       # desarrollo
npm run build     # produccion
npm run preview   # previsualizar build
npm run lint      # ESLint
```
