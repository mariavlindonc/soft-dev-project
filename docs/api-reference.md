# API Reference

Todas las rutas usan el prefijo `/api` (excepto `/health`). Las fechas se envian en formato **RFC 3339** (`2026-06-10T15:00:00Z`).

---

## Publicas

| Metodo | Ruta | Autenticacion | Descripcion |
|--------|------|:-------------:|-------------|
| `GET` | `/health` | - | Health check del servidor |
| `POST` | `/api/auth/register` | - | Registrar nuevo usuario |
| `POST` | `/api/auth/login` | - | Iniciar sesion, devuelve JWT |
| `GET` | `/api/events` | - | Listar eventos (filtros: `?category=&date=`) |
| `GET` | `/api/events/:id` | - | Obtener detalle de un evento |
| `GET` | `/api/events/:id/sale-status` | - | Obtener fase de venta actual |

## Autenticadas (requieren `Authorization: Bearer <token>`)

| Metodo | Ruta | Rol | Descripcion |
|--------|------|:---:|-------------|
| `POST` | `/api/tickets/purchase` | client | Comprar una entrada (opcional: `presale_code`) |
| `GET` | `/api/tickets` | client | Listar entradas del usuario autenticado |
| `PATCH` | `/api/tickets/:id/cancel` | client | Cancelar una entrada propia |
| `PATCH` | `/api/tickets/:id/transfer` | client | Transferir una entrada a otro usuario por email |

## Administrador (requieren `Authorization: Bearer <token>` + rol admin)

| Metodo | Ruta | Rol | Descripcion |
|--------|------|:---:|-------------|
| `POST` | `/api/admin/events` | admin | Crear un nuevo evento |
| `PUT` | `/api/admin/events/:id` | admin | Actualizar un evento |
| `DELETE` | `/api/admin/events/:id` | admin | Cancelar un evento (baja logica) |
| `GET` | `/api/admin/reports` | admin | Reporte global de todos los eventos |
| `GET` | `/api/admin/reports/events/:id` | admin | Reporte detallado de un evento (incluye compradores) |

---

## Ejemplos de uso

```bash
# 1. Registrar usuario
curl -X POST https://localhost:8443/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan Perez","email":"juan@test.com","password":"12345678"}'

# 2. Iniciar sesion
curl -X POST https://localhost:8443/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"juan@test.com","password":"12345678"}'
# → {"token":"eyJhbGciOiJIUzI1NiIs..."}

# 3. Listar eventos
curl https://localhost:8443/api/events

# 4. Comprar entrada (con token)
curl -X POST https://localhost:8443/api/tickets/purchase \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"event_id":1}'

# 5. Ver mis entradas
curl https://localhost:8443/api/tickets \
  -H "Authorization: Bearer <token>"

# 6. Cancelar entrada
curl -X PATCH https://localhost:8443/api/tickets/1/cancel \
  -H "Authorization: Bearer <token>"

# 7. Transferir entrada
curl -X PATCH https://localhost:8443/api/tickets/1/transfer \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"to_user_email":"ana@test.com"}'

# 8. Crear evento (admin)
curl -X POST https://localhost:8443/api/admin/events \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Concierto Rock",
    "event_date":"2026-08-15T20:00:00Z",
    "capacity":500,
    "price":2500.00,
    "category":"musica"
  }'

# 9. Reporte global (admin)
curl https://localhost:8443/api/admin/reports \
  -H "Authorization: Bearer <token>"
```

---

## Request/Response DTOs

### Auth

**`POST /api/auth/register`**
```json
{
  "name": "Juan Perez",
  "email": "juan@test.com",
  "password": "12345678"
}
```
```json
// 201 Created
{
  "id": 1,
  "name": "Juan Perez",
  "email": "juan@test.com",
  "role": "client"
}
```

**`POST /api/auth/login`**
```json
{
  "email": "juan@test.com",
  "password": "12345678"
}
```
```json
// 200 OK
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Events

**`POST /api/admin/events`**
```json
{
  "title": "Concierto Rock",
  "description": "La mejor banda en vivo",
  "category": "musica",
  "location": "Estadio Unico",
  "event_date": "2026-08-15T20:00:00Z",
  "duration_minutes": 180,
  "capacity": 500,
  "price": 2500.00,
  "presale_active": true,
  "presale_code": "ROCK2026",
  "presale_start_date": "2026-07-01T00:00:00Z",
  "general_sale_date": "2026-08-01T00:00:00Z"
}
```
```json
// 201 Created
{
  "id": 1,
  "title": "Concierto Rock",
  "status": "active",
  "presale_active": true,
  "capacity": 500,
  "tickets_sold": 0,
  "price": 2500.00,
  ...
}
```

**`GET /api/events`**
```json
// 200 OK
[
  {
    "id": 1,
    "title": "Concierto Rock",
    "event_date": "2026-08-15T20:00:00Z",
    "status": "active",
    "price": 2500.00,
    "capacity": 500,
    "tickets_sold": 0,
    ...
  }
]
```

**`GET /api/events/:id/sale-status`**
```json
// 200 OK
{
  "phase": "presale",
  "presale_start_date": "2026-07-01T00:00:00Z",
  "general_sale_date": "2026-08-01T00:00:00Z",
  "message": "Pre-sale is currently active. An access code is required to purchase tickets."
}
```

### Tickets

**`POST /api/tickets/purchase`**
```json
{
  "event_id": 1,
  "presale_code": "ROCK2026"
}
```
```json
// 201 Created
{
  "id": 10,
  "user_id": 1,
  "event_id": 1,
  "status": "active",
  "purchase_price": 2500.00,
  "purchased_at": "2026-07-15T10:30:00Z"
}
```

**`PATCH /api/tickets/:id/transfer`**
```json
{
  "to_user_email": "ana@test.com"
}
```
```json
// 200 OK
{
  "message": "ticket transferred successfully"
}
```

---

## Codigos de error HTTP

| Codigo | Significado |
|--------|-------------|
| `200` | OK |
| `201` | Creado exitosamente |
| `204` | Exitoso sin contenido (DELETE, cancelaciones) |
| `400` | Error de validacion (campos faltantes, formato invalido) |
| `401` | No autenticado (token faltante o invalido) |
| `403` | No autorizado (rol incorrecto) |
| `404` | Recurso no encontrado |
| `500` | Error interno del servidor |
