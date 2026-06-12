# Modelo de Base de Datos

![Diagrama ER](Diagrama.jpg)

El esquema completo se encuentra en [database/schema.sql](../database/schema.sql).

---

## Entidades

### users - Clientes y administradores

| Campo | Tipo | Restricciones | Detalle |
|-------|------|:-------------|---------|
| `id` | INT UNSIGNED | PK, AUTO_INCREMENT | Identificador unico |
| `name` | VARCHAR(150) | NOT NULL | Nombre del usuario |
| `email` | VARCHAR(255) | NOT NULL, UNIQUE | Email de inicio de sesion |
| `password_hash` | VARCHAR(255) | NOT NULL | Hash bcrypt (nunca texto plano) |
| `role` | ENUM('client','admin') | NOT NULL, DEFAULT 'client' | Rol en el sistema |
| `created_at` | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Fecha de registro |
| `updated_at` | DATETIME | NOT NULL, ON UPDATE CURRENT_TIMESTAMP | Ultima modificacion |
| `deleted_at` | DATETIME | NULL | Soft-delete (GORM) |

### events - Eventos con soporte de preventa

| Campo | Tipo | Restricciones | Detalle |
|-------|------|:-------------|---------|
| `id` | INT UNSIGNED | PK, AUTO_INCREMENT | Identificador unico |
| `title` | VARCHAR(200) | NOT NULL | Titulo del evento |
| `description` | TEXT | NULL | Descripcion larga |
| `image_url` | VARCHAR(500) | NULL | URL de imagen ilustrativa |
| `category` | VARCHAR(100) | NULL, INDEX | Categoria para filtros y busqueda |
| `location` | VARCHAR(300) | NULL | Lugar del evento |
| `event_date` | DATETIME | NOT NULL, INDEX | Fecha y hora del evento |
| `duration_minutes` | INT | NOT NULL, DEFAULT 0 | Duracion estimada |
| `capacity` | INT | NOT NULL, DEFAULT 0 | Capacidad total de asistentes |
| `tickets_sold` | INT | NOT NULL, DEFAULT 0 | Entradas vendidas (CHECK: `<= capacity`) |
| `price` | DECIMAL(10,2) | NOT NULL, DEFAULT 0.00 | Precio por entrada |
| `status` | ENUM('active','presale','sold_out','cancelled') | NOT NULL, DEFAULT 'active', INDEX | Estado actual del evento |
| `presale_active` | TINYINT(1) | NOT NULL, DEFAULT 0 | Habilita/deshabilita preventa |
| `presale_code` | VARCHAR(100) | NULL | Codigo de acceso para preventa (no expuesto en JSON) |
| `presale_start_date` | DATETIME | NULL | Inicio de la venta en preventa |
| `general_sale_date` | DATETIME | NULL | Inicio de la venta general |
| `created_by_id` | INT UNSIGNED | NOT NULL, FK → users(id) | Administrador que creo el evento |
| `created_at` | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Fecha de creacion |
| `updated_at` | DATETIME | NOT NULL, ON UPDATE CURRENT_TIMESTAMP | Ultima modificacion |
| `deleted_at` | DATETIME | NULL, INDEX | Soft-delete (GORM) |

**Indices compuestos:**
- `idx_events_status_date` (status, event_date) - Optimiza busquedas por estado + fecha

**Check constraints:**
- `chk_events_capacity`: `tickets_sold >= 0 AND tickets_sold <= capacity`

### tickets - Entradas adquiridas por los usuarios

| Campo | Tipo | Restricciones | Detalle |
|-------|------|:-------------|---------|
| `id` | INT UNSIGNED | PK, AUTO_INCREMENT | Identificador unico |
| `user_id` | INT UNSIGNED | NOT NULL, FK → users(id), INDEX | Propietario actual de la entrada |
| `event_id` | INT UNSIGNED | NOT NULL, FK → events(id), INDEX | Evento al que pertenece |
| `status` | ENUM('active','cancelled','transferred') | NOT NULL, DEFAULT 'active' | Estado actual |
| `purchase_price` | DECIMAL(10,2) | NOT NULL, DEFAULT 0.00 | Precio pagado al momento de la compra |
| `purchased_at` | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Momento de la compra |
| `cancelled_at` | DATETIME | NULL | Momento de cancelacion |
| `transferred_at` | DATETIME | NULL | Momento de transferencia |
| `transferred_to_id` | INT UNSIGNED | NULL, FK → users(id), INDEX | Usuario destino de la transferencia |
| `created_at` | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Fecha de creacion del registro |

**Indices compuestos:**
- `idx_tickets_event_status` (event_id, status) - Busquedas de entradas activas por evento
- `idx_tickets_user_id_status` (user_id, status) - Historial de entradas por usuario

---

## Relaciones

```
users 1 ── * tickets     (un usuario puede tener muchas entradas)
users 1 ── * events      (un admin puede crear muchos eventos)
events 1 ── * tickets    (un evento tiene muchas entradas)
tickets * ── 1 tickets  (una entrada transferida apunta al nuevo dueño via transferred_to_id)
```

## Diagrama entidad-relacion

![Diagrama ER](Diagrama.jpg)
