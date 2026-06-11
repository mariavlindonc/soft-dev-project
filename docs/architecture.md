# Decisiones de Diseno

## 1. DAOs como interfaces para testabilidad

Los repositorios (UserDAO, EventDAO, TicketDAO) se definen como interfaces en `dao/interfaces.go`. Esto permite inyectar mocks en los servicios y testear la logica de negocio sin base de datos real, siguiendo el principio de inversion de dependencias.

```go
type UserDAO interface {
    Create(user *domain.User) error
    FindByEmail(email string) (*domain.User, error)
    FindByID(id uint) (*domain.User, error)
}
```

Cada implementacion concreta verifica en compilacion que satisface la interfaz:

```go
var _ UserDAO = (*UserDAOImpl)(nil)
```

---

## 2. Fases de venta con logica de dominio pura

La preventa se modela con fechas (`presale_start_date`, `general_sale_date`) y un metodo `CurrentSalePhase()` en la entidad `Event`. Esto centraliza la logica en el dominio y la hace testeable sin depender del servicio o la base de datos.

```go
func (e *Event) CurrentSalePhase(now time.Time) SalePhase {
    if !e.PresaleActive || e.PresaleStartDate == nil || e.GeneralSaleDate == nil {
        return PhaseNoPresale
    }
    if now.Before(*e.PresaleStartDate) {
        return PhaseNotYetOpen
    }
    if now.Before(*e.GeneralSaleDate) {
        return PhasePresale
    }
    return PhasePublic
}
```

Las fases posibles:

| Fase | Descripcion |
|------|-------------|
| `not_yet_open` | La preventa no ha comenzado aun |
| `presale` | Preventa activa (requiere codigo de acceso) |
| `public` | Venta general abierta (sin codigo) |
| `no_presale` | El evento no tiene preventa configurada |

---

## 3. Email client con strategy pattern

El cliente de email (`clients/email_client.go`) tiene dos implementaciones intercambiables via variable de entorno:

| Implementacion | Variable de entorno | Comportamiento |
|----------------|--------------------|----------------|
| `logEmailClient` | `EMAIL_PROVIDER=log` (default) | Solo registra en consola el contenido del email |
| `smtpEmailClient` | `EMAIL_PROVIDER=smtp` | Envia emails reales via SMTP con timeout de 30s |

```go
type EmailClient interface {
    SendPurchaseConfirmation(to string, info TicketInfo) error
    SendCancellationNotice(to string, info TicketInfo) error
    SendTransferNotice(from, to string, info TicketInfo) error
}
```

La interfaz permite mockear en tests y cambiar la implementacion sin modificar codigo de negocio.

---

## 4. Transacciones atomicas

Las operaciones criticas (compra, cancelacion, transferencia de entradas) se ejecutan dentro de transacciones SQL, garantizando consistencia entre la creacion/modificacion del ticket y el ajuste del contador `tickets_sold` del evento.

Flujo de compra:
1. Iniciar transaccion
2. Verificar capacidad disponible (`tickets_sold < capacity`)
3. Crear ticket con status `active`
4. Incrementar `tickets_sold` del evento
5. Commitear transaccion (o rollback si algo falla)
6. Enviar email de confirmacion (fire-and-forget, errores solo se loguean)

---

## 5. Seguridad por diseno

- **HTTPS obligatorio** - El servidor no inicia sin TLS_CERT_FILE y TLS_KEY_FILE configurados
- **Headers de seguridad** - HSTS, X-Content-Type-Options, X-Frame-Options (DENY), CSP (`default-src 'self'`), Referrer-Policy
- **CORS configurable** - Origenes permitidos via `CORS_ALLOWED_ORIGINS` (default: `*`)
- **Bcrypt costo 12** - Hash de contrasenas con factor de trabajo elevado
- **JWT 24h** - Tokens con expiracion fija, secreto configurable via `JWT_SECRET`
- **Campos sensibles excluidos del JSON** - `password_hash` y `presale_code` usan `json:"-"` en las entidades
- **No informacion de existencia** - Violaciones de propiedad de tickets devuelven `404 Not Found` en lugar de `403`, para no filtrar la existencia del recurso

---

## 6. Capas del backend

```
┌─────────────────────────────────────────────────────────────┐
│                        controllers                          │
│              Handlers HTTP + middleware + DTOs               │
├─────────────────────────────────────────────────────────────┤
│                         services                            │
│          Logica de negocio + validaciones + errores          │
├─────────────────────────────────────────────────────────────┤
│                     dao (interfaces)                         │
│          Abstraccion de repositorio (interfaces)             │
├─────────────────────────────────────────────────────────────┤
│                     dao (implementaciones)                   │
│                    GORM + MySQL                              │
└─────────────────────────────────────────────────────────────┘
```

- Los controllers solo manejan HTTP (binding, response codes) y delegan a services
- Los services contienen toda la logica de negocio y validaciones
- Los DAOs son intercambiables (mocks en tests, GORM en produccion)
- Las entidades de dominio (`domain/`) son estructuras planas sin dependencias externas
