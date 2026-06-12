# Testing

```bash
cd backend
go test ./... -v -cover
```

---

## Cobertura por paquete

| Paquete | Tipo | Tests |
|---------|------|:-----:|
| `domain/` | Unitario (puro, sin dependencias) | 9/9 |
| `utils/` | Unitario (puro) | 9/9 |
| `services/` | Unitario (con testify/mock) | 27/27 |
| `controllers/` | Integracion (httptest) | 27/27 |

---

## Tests de dominio

Prueban logica pura sin dependencias externas.

| Archivo | Casos cubiertos |
|---------|-----------------|
| `event_test.go` | `CurrentSalePhase` - 9 subtests que cubren las 4 fases de venta mas casos borde con fechas nil y valores exactos en los limites |

---

## Tests de utils

Prueban funciones de utilidad aisladas.

| Archivo | Casos cubiertos |
|---------|-----------------|
| `password_test.go` | HashPassword, CheckPassword (correcta, incorrecta, vacia), HashPasswordUnique |
| `jwt_test.go` | GenerateJWT, ValidateToken (valido, firma invalida, malformado, vacio) |

---

## Tests de servicios

Usan **testify/mock** para simular los DAOs y el cliente de email. Cada servicio se testea de forma aislada sin base de datos.

| Archivo | Casos cubiertos |
|---------|-----------------|
| `auth_service_test.go` | Register (exito, password corta, email duplicado), Login (exito, password incorrecta, email desconocido) |
| `event_service_test.go` | GetAll, GetByID (exito, no encontrado), Create (valido, titulo vacio, capacidad cero, fecha pasada), Cancel (activo, ya cancelado, no encontrado), Update (exito, cancelado, no encontrado) |
| `ticket_service_test.go` | Purchase (exito, evento cancelado, sin capacidad), PurchasePresale (codigo correcto, sin codigo, codigo incorrecto), CancelTicket (propio, ajeno, ya cancelado), Transfer (a otro usuario, a si mismo), GetByUser, PurchaseNotFound |
| `report_service_test.go` | GetEventReport (exito, no encontrado), GetGlobalReport |

---

## Tests de controladores

Usan **net/http/httptest** para enviar requests HTTP directamente contra los handlers con servicios mockeados.

| Archivo | Casos cubiertos |
|---------|-----------------|
| `middleware_test.go` | AuthRequired (sin header, formato invalido, token malformado), AdminRequired (admin pasa, client 403, sin rol 403) |
| `auth_controller_test.go` | Register (201, 400 campos faltantes, 400 email duplicado), Login (200 con token, 401 credenciales invalidas, 400 campo faltante) |
| `event_controller_test.go` | GetAll (200), GetByID (200, 404, 400), Create (201, 400), Update (200, 404), Delete (204, 404), GetSaleStatus (200) |
| `ticket_controller_test.go` | Purchase (201, 400 evento cancelado, 400 campo faltante), GetMyTickets (200), Cancel (204, 404), Transfer (200, 400 campo faltante) |
| `admin_controller_test.go` | GetReports (200), GetEventReport (200, 404) |
