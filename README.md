# Ceibo Tickets - Sistema de Gestion de Eventos y Entradas

Sistema de venta de entradas para eventos con roles de cliente y administrador, desarrollado como practico integrador 2026.

---

## Tecnologias utilizadas

| Capa              | Tecnologia                     |
| ----------------- | ------------------------------ |
| **Backend**       | Go 1.26, Gin, GORM, MySQL      |
| **Frontend**      | React 19, TypeScript 6, Vite 8 |
| **Autenticacion** | JWT (golang-jwt), bcrypt       |
| **Testing**       | Go testing + testify           |

> [Arquitectura detallada](docs/architecture.md) — [API Reference](docs/api-reference.md)

---

## Instalacion y uso

```bash
# Clonar
git clone <repo-url>
cd soft-dev-project

# Backend
cd backend
cp .env.example .env    # configurar credenciales
go mod tidy
go run main.go

# Frontend
cd frontend
npm install
npm run dev
```

**Requisitos previos:**
Go 1.26+, Node.js 26+, MySQL 8+

> Proba la API con los [ejemplos de curl](docs/api-reference.md#ejemplos-de-uso) luego de iniciar el servidor.

---

## Variables de entorno

Archivo `.env` en `backend/`. Ver [`.env.example`](backend/.env.example) para valores completos.

| Variable                                                  | Default                    | Descripcion                                 |
| --------------------------------------------------------- | -------------------------- | ------------------------------------------- |
| `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` | `localhost:3306/ceibo_db`  | Conexion MySQL ([esquema](docs/database.md)) |
| `JWT_SECRET`                                              | `change-me-...`            | Clave para firmar JWT                       |
| `APP_PORT`                                                | `8443`                     | Puerto HTTPS                                |
| `TLS_CERT_FILE`, `TLS_KEY_FILE`                           | `./certs/server.{crt,key}` | TLS obligatorio                             |
| `CORS_ALLOWED_ORIGINS`                                    | `*`                        | Origenes CORS permitidos                    |
| `EMAIL_PROVIDER`                                          | `log`                      | `log` (desarrollo) o `smtp` (produccion)    |
| `SMTP_*`                                                  | -                          | Configuracion SMTP (si EMAIL_PROVIDER=smtp) |

---

## Estructura del proyecto

```
soft-dev-project/
├── backend/
│   ├── main.go              # Entry point, DI, rutas, middlewares
│   ├── clients/             # Email client — [strategy pattern](docs/architecture.md#3-email-client-con-strategy-pattern)
│   ├── controllers/         # Handlers HTTP — [API reference](docs/api-reference.md)
│   ├── dao/                 # Interfaces + [implementaciones GORM](docs/architecture.md#1-daos-como-interfaces-para-testabilidad)
│   ├── domain/              # Entidades User, Event, Ticket — [modelo de datos](docs/database.md)
│   ├── logger/              # Logger estructurado JSON
│   └── services/            # Logica de negocio + [tests](docs/testing.md#tests-de-servicios)
├── frontend/                # React + TypeScript + Vite
├── database/                # Schema SQL + datos de ejemplo
├── docs/                    # Documentacion detallada
└── README.md
```

---

## Funcionalidades

### Cliente

- Exploracion y filtro de eventos por categoria y fecha
- Detalle de evento con informacion completa
- Compra de entradas (con soporte de preventa y codigo de acceso)
- Historial "Mis Entradas" con estado de cada una
- Cancelacion de compra
- Traspaso de entrada a otro usuario por email
- Notificaciones por email (confirmacion, cancelacion, transferencia)

### Administrador

- Creacion, actualizacion y cancelacion de eventos
- Reportes y metricas globales (ocupacion, ventas totales)
- Reporte detallado por evento (incluye lista de compradores)

### Funcionalidad Extra (Bonus Track)

- Preventa con codigo de acceso y [fechas diferenciadas](docs/architecture.md#2-fases-de-venta-con-logica-de-dominio-pura)
- Fases de venta: no abierta → preventa → venta general
- [Logging estructurado JSON](docs/architecture.md#3-email-client-con-strategy-pattern), email intercambiable, [transacciones atomicas](docs/architecture.md#4-transacciones-atomicas)

---

## Documentacion detallada

| Documento                              | Contenido                                                     |
| -------------------------------------- | ------------------------------------------------------------- |
| [API Reference](docs/api-reference.md) | Endpoints, ejemplos curl, request/response DTOs, codigos HTTP |
| [Base de Datos](docs/database.md)      | Entidades, campos, relaciones, diagrama ER                    |
| [Arquitectura](docs/architecture.md)   | Patrones, diseno, capas, seguridad                            |
| [Testing](docs/testing.md)             | Tests por paquete, cobertura, casos                           |

---

## Estado del proyecto

| Componente    |    Estado     | Detalle                                                       |
| ------------- | :-----------: | ------------------------------------------------------------- |
| Backend API   |   Completo    | 16 endpoints, 72+ tests, arquitectura limpia                  |
| Frontend      | En desarrollo | Template inicial. Pendiente: routing, pages, API, auth, forms |
| Base de datos |   Completo    | Schema SQL + migraciones GORM                                 |
| Documentacion |   Completo    | README + docs detallados                                      |

---

## Autores

BRUA, Jonathan
HERNANDEZ, Juan
LINDON, Maria Victoria
.., Athina

Proyecto integrador 2026 - Universidad Catolica de Cordoba (UCC).
