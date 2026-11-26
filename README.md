# Ceibo Tickets - Sistema de Gestion de Eventos y Entradas

Sistema de venta de entradas para eventos con roles de **cliente** y **administrador**, desarrollado como practico integrador 2026 de la Universidad Catolica de Cordoba. El backend expone una API REST en Go con arquitectura MVC, el frontend es una SPA en React + TypeScript, y los datos se persisten en MySQL con GORM.

El flujo del sistema es el siguiente:
El usuario ingresa al sitio, explora el catalogo de eventos con filtros por categoria y fecha, y selecciona un evento para ver su detalle (descripcion, fecha, ubicacion, capacidad, precio). Si desea comprar, debe registrarse o iniciar sesion; luego puede adquirir una o mas entradas, ingresando un codigo de preventa si corresponde. La compra se procesa con validaciones de capacidad y fase de venta, y se confirma via email. El usuario puede ver sus entradas en "Mis Entradas", cancelarlas o transferirlas a otro usuario por email. Del lado del administrador, este puede crear, editar y cancelar eventos, configurar preventas con codigo de acceso y fechas diferenciadas, y consultar reportes globales o por evento con detalle de compradores.

---

## Tabla de contenidos

1. [Tecnologias utilizadas](#tecnologias-utilizadas)
2. [Requisitos previos](#requisitos-previos)
3. [Instalacion y uso](#instalacion-y-uso)
4. [Estructura del proyecto](#estructura-del-proyecto)
5. [Capturas de pantalla](#capturas-de-pantalla)
6. [Diagrama de base de datos](#diagrama-de-base-de-datos)
7. [Endpoints de la API](#endpoints-de-la-api)
8. [Decisiones de diseno](#decisiones-de-diseno)
9. [Testing](#testing)
10. [Bonus Track](#bonus-track)
11. [Autores](#autores)

---

## Tecnologias utilizadas

| Capa              | Tecnologia                                    |
| ----------------- | --------------------------------------------- |
| **Backend**       | Go 1.26, Gin, GORM, driver MySQL (go-sql-driver) |
| **Base de datos** | MySQL 8                                       |
| **Frontend**      | React 19, TypeScript, Vite, React Router, Axios |
| **Autenticacion** | JWT (golang-jwt), bcrypt                      |
| **Testing**       | Go testing + testify + httptest + SQLite en memoria (DAOs) |
| **DevOps**        | Docker, Docker Compose, Nginx                 |
| **Seguridad**     | TLS, headers de seguridad, CORS configurable  |

---

## Requisitos previos

- Go 1.26+
- Node.js 26+
- MySQL 8+ (solo si corre local sin Docker)
- Docker + Docker Compose (para levantar con contenedores)

---

## Instalacion y uso

### Opcion 1: Con Docker (recomendado)

Levantá todo el stack (base de datos, backend y frontend) con un solo comando:

```bash
git clone <repo-url>
cd soft-dev-project
docker-compose up --build
```

Los servicios quedan disponibles en:

| Servicio  | URL                          |
| --------- | ---------------------------- |
| Frontend  | http://localhost              |
| Backend   | https://localhost:8443        |
| MySQL     | localhost:3307                |

Para detener y limpiar:

```bash
docker-compose down           # detiene los contenedores
docker-compose down -v        # detiene y elimina los volúmenes (borra la base)
```

Los datos seed (13 usuarios, 23 eventos, 35 tickets) se cargan automaticamente al iniciar la base por primera vez. Algunos usuarios de prueba:

| Email                       | Password     | Rol      |
| --------------------------- | ------------ | -------- |
| carlos@ceibo.com            | carlos123    | admin    |
| ana@ceibo.com               | analucia123  | admin    |
| pedro@ceibo.com             | pedro123     | admin    |
| sofia@gmail.com             | sofia123     | cliente  |
| lautaro@hotmail.com         | lautaro123   | cliente  |

El resto de los clientes sigue el patron `[nombre]123` (ver `database/data.sql`).

> **Si el backend no arranca la primera vez:** al crear la base desde cero, MySQL tarda unos segundos en cargar el esquema y los datos, y el backend puede intentar conectarse antes y cerrarse. Se ve como `Exited` en `docker-compose ps`. Basta con esperar unos segundos y ejecutar `docker-compose start backend`.

### Opcion 2: Local (sin Docker)

```bash
# Clonar
git clone <repo-url>
cd soft-dev-project

# Backend
cd backend
cp .env.example .env   # configurar credenciales de MySQL
go mod tidy
go run ./cmd/gencerts  # genera certs/server.crt y certs/server.key (solo la primera vez)
go run main.go

# En otra terminal - Frontend
cd frontend
npm install
npm run dev
```

Las variables de entorno se configuran en `backend/.env` (ver `.env.example`). El servidor inicia en `https://localhost:8443` (TLS obligatorio) y el frontend en `http://localhost:5173`.

### Ejecutar los tests

```bash
# Backend: tests unitarios y de integracion con cobertura
cd backend
go test ./... -cover

# Frontend: chequeo de tipos, build y lint
cd frontend
npm run build
npm run lint
```

El detalle de los tests y la cobertura esta en la seccion [Testing](#testing).

---

## Estructura del proyecto

```
soft-dev-project/
├── backend/
│   ├── main.go                 # Punto de entrada, configura servidor HTTP y rutas
│   ├── Dockerfile              # Build multi-etapa (golang -> alpine)
│   ├── .env.example            # Variables de entorno de ejemplo
│   ├── cmd/
│   │   └── gencerts/main.go    # Generador de certificados auto-firmados
│   ├── domain/                 # Entidades del negocio (User, Event, Ticket)
│   ├── dao/                    # Acceso a datos via GORM (interfaces + implementaciones)
│   ├── services/               # Logica de negocio
│   ├── controllers/            # Handlers HTTP + middleware (auth, admin, CORS, seguridad)
│   ├── clients/                # Clientes externos (email via SMTP o log)
│   ├── utils/                  # Utilidades (JWT, password hashing)
│   ├── logger/                 # Logger estructurado JSON
│   └── certs/                  # Certificados TLS
├── frontend/
│   ├── Dockerfile              # Build multi-etapa (node -> nginx)
│   ├── nginx.conf              # SPA fallback + reverse proxy a backend
│   ├── src/
│   │   ├── api/                # Capa de cliente HTTP (axios)
│   │   ├── types/              # Definiciones TypeScript
│   │   ├── context/            # AuthContext (estado de sesion)
│   │   ├── hooks/              # Custom hooks (useEvents, useTickets)
│   │   ├── pages/              # Vistas: Home, Events, Detail, Auth, Tickets, Admin
│   │   ├── components/         # Componentes reutilizables (cards, modals, layout)
│   │   └── data/               # Datos mock como fallback
│   └── public/                 # Assets estaticos
├── database/
│   ├── schema.sql              # DDL: CREATE TABLE con FK, indices, constraints
│   └── data.sql                # Datos seed (13 usuarios, 23 eventos, 35 tickets)
├── docs/
│   ├── Diagrama.png            # Diagrama ER (imagen)
│   ├── Diagrama.svg            # Diagrama ER (fuente editable)
│   └── screenshots/            # Capturas de pantalla del sistema
├── docker-compose.yml          # Orquestacion de 3 servicios
└── README.md
```

---

## Capturas de pantalla

| Vista                      | Imagen                                               |
| -------------------------- | ---------------------------------------------------- |
| Pagina principal           | ![Home](docs/screenshots/home.jpg)                   |
| Catalogo de eventos        | ![Eventos](docs/screenshots/home-events.jpg)         |
| Detalle de evento          | ![Detalle](docs/screenshots/event-detail.jpg)        |
| Inicio de sesion           | ![Login](docs/screenshots/login.jpg)                 |
| Registro de usuario        | ![Register](docs/screenshots/register.jpg)           |
| Mis Entradas               | ![Mis Entradas](docs/screenshots/my-tickets.jpg)     |
| Flujo de compra            | ![Compra](docs/screenshots/purchase.jpg)             |
| Panel administrador        | ![Admin](docs/screenshots/admin-panel.jpg)           |
| Gestion de eventos (Admin) | ![Admin Eventos](docs/screenshots/admin-events.jpg)  |
| Reporte de evento (Admin)  | ![Reporte](docs/screenshots/admin-report.jpg)        |

---

## Diagrama de base de datos

![Diagrama ER](docs/Diagrama.png)

La imagen esta en `docs/Diagrama.png` y su fuente editable en `docs/Diagrama.svg`.

El modelo consta de tres tablas principales:

- **users**: almacenamiento de usuarios con roles (cliente/admin), email unico, password hasheada con bcrypt, y soft-delete.
- **events**: eventos con soporte para preventa (codigo + fechas), estados (active, presale, sold_out, cancelled), capacidad maxima y soft-delete.
- **tickets**: entradas vendidas con estados (active, cancelled, transferred), precio de compra, y referencia al usuario dueño, al evento y al destinatario de una transferencia.

Relaciones implementadas mediante claves foraneas (todas 1 a N):

| Clave foranea               | Referencia  | Significado                                       |
| --------------------------- | ----------- | ------------------------------------------------- |
| `events.created_by_id`      | `users.id`  | Administrador que creo el evento                  |
| `tickets.user_id`           | `users.id`  | Dueño actual de la entrada                        |
| `tickets.event_id`          | `events.id` | Evento al que pertenece la entrada                |
| `tickets.transferred_to_id` | `users.id`  | Destinatario de una transferencia (opcional, NULL) |

El esquema completo esta en `database/schema.sql` con datos de ejemplo en `database/data.sql` (13 usuarios, 23 eventos, 35 tickets).

---

## Endpoints de la API

### Autenticacion (publicos)

| Metodo | Endpoint                | Descripcion            |
| ------ | ----------------------- | ---------------------- |
| POST   | `/api/auth/register`    | Registro de usuario    |
| POST   | `/api/auth/login`       | Login y obtencion JWT  |

### Eventos (publicos)

| Metodo | Endpoint                        | Descripcion                             |
| ------ | ------------------------------- | --------------------------------------- |
| GET    | `/api/events`                   | Listado de eventos (filtros opcionales) |
| GET    | `/api/events/:id`               | Detalle de un evento                    |
| GET    | `/api/events/:id/sale-status`   | Estado de venta de un evento            |

### Cliente (requieren JWT)

| Metodo | Endpoint                        | Descripcion                           |
| ------ | ------------------------------- | ------------------------------------- |
| POST   | `/api/tickets/purchase`         | Compra de entrada                     |
| GET    | `/api/tickets`                  | Mis entradas                          |
| PATCH  | `/api/tickets/:id/cancel`       | Cancelar entrada                      |
| PATCH  | `/api/tickets/:id/transfer`     | Transferir entrada a otro usuario     |

### Administrador (requieren JWT + rol admin)

| Metodo | Endpoint                          | Descripcion                    |
| ------ | --------------------------------- | ------------------------------ |
| POST   | `/api/admin/events`               | Crear evento                   |
| PUT    | `/api/admin/events/:id`           | Actualizar evento              |
| DELETE | `/api/admin/events/:id`           | Cancelar evento                |
| GET    | `/api/admin/reports`              | Reporte global de ocupacion    |
| GET    | `/api/admin/reports/events/:id`   | Reporte de un evento specifico |

### Seguridad

- **JWT**: tokens con 24h de expiracion, firmados con HS256. El payload contiene `user_id` y `role`.
- **Passwords**: hasheadas con bcrypt (cost 12), nunca se almacenan en texto plano.
- **Middleware**: `AuthRequired` valida el token Bearer, `AdminRequired` verifica el rol.
- **Headers**: HSTS, X-Content-Type-Options, X-Frame-Options, CSP, Referrer-Policy.
- **CORS**: configurable via variable de entorno `CORS_ALLOWED_ORIGINS`.

---

## Decisiones de diseno

**1. DAOs como interfaces para testabilidad.** Los repositorios (UserDAO, EventDAO, TicketDAO) se definen como interfaces en `dao/interfaces.go`. Esto permite inyectar mocks en los servicios y testear la logica de negocio sin base de datos real, siguiendo el principio de inversion de dependencias.

```go
type UserDAO interface {
    Create(user *User) error
    FindByEmail(email string) (*User, error)
    FindByID(id uint) (*User, error)
}
```

**2. Fases de venta con logica de dominio pura.** La preventa se modela con fechas en la entidad `Event` y un metodo `CurrentSalePhase()` que determina la fase actual (no abierta, preventa, venta general, sin preventa). Esto centraliza la logica en el dominio y la hace testeable sin depender del servicio ni de la base de datos.

**3. Email client con strategy pattern.** El cliente de email tiene dos implementaciones intercambiables via variable de entorno (`EMAIL_PROVIDER=log` en desarrollo, `EMAIL_PROVIDER=smtp` en produccion). La interfaz permite mockear en tests y cambiar la implementacion sin modificar el codigo de negocio.

**4. Transacciones atomicas.** Las operaciones criticas (compra, cancelacion, transferencia) se ejecutan dentro de transacciones SQL via `WithTransaction()`, garantizando consistencia entre la creacion/modificacion del ticket y el ajuste del contador `tickets_sold` del evento. Si algo falla, se revierte todo (rollback).

---

## Testing

```bash
# Todos los tests
cd backend
go test ./... -v -cover

# Tests de un paquete especifico
go test ./services/... -v -cover
go test ./controllers/... -v -cover

# Generar reporte de cobertura HTML
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Sin tener Go instalado, usando Docker (la imagen del backend solo contiene el binario compilado, por eso se usa la imagen oficial de Go):

```bash
docker run --rm -v "${PWD}/backend:/app" -w /app golang:1.26-alpine go test ./... -cover
```

### Cobertura por paquete

**360 tests, todos pasan. Cobertura total del backend: 81.7%** (objetivo de la consigna: 80%, principalmente en servicios y controladores).

| Paquete         | Tipo                                    | Cobertura |
| --------------- | --------------------------------------- | :-------: |
| `services/`     | Unitario (con testify/mock)             | **99.3%** |
| `controllers/`  | Integracion HTTP (httptest)             | **95.2%** |
| `domain/`       | Unitario (puro, sin dependencias)       |  100.0%   |
| `logger/`       | Unitario (con testify)                  |  100.0%   |
| `clients/`      | Unitario (con testify)                  |   93.5%   |
| `utils/`        | Unitario (puro)                         |   90.9%   |
| `cmd/gencerts/` | Unitario (archivos temporales)          |   53.3%   |
| `dao/`          | Integracion (SQLite en memoria)         |   48.5%   |
| `backend` (`main.go`) | Unitario (middlewares CORS y seguridad) | 38.1%  |
| **Total**       |                                         | **81.7%** |

Los paquetes con menor cobertura lo son por diseño: en `dao/`, el archivo `connection.go` (conexion a MySQL, migraciones y TLS) necesita un servidor MySQL real; en `main.go` y `cmd/gencerts/`, la funcion `main()` solo arma y arranca la aplicacion. La logica de negocio (servicios) y la capa HTTP (controladores) superan ampliamente el 80%.

Para ver la cobertura linea por linea:

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out      # resumen por funcion (la ultima linea es el total)
go tool cover -html=coverage.out      # abre el reporte en el navegador
```

### Estrategia de testing

- **Dominio**: pruebas puras de logica sin dependencias externas. `CurrentSalePhase()` se testea con 9 subtests cubriendo las 4 fases de venta mas casos borde.
- **Servicios**: usan **testify/mock** para simular los DAOs y el cliente de email. Cada servicio se testea de forma aislada sin base de datos.
- **Controladores**: usan **net/http/httptest** para enviar requests HTTP directamente contra los handlers con servicios mockeados, validando status codes en casos exitosos y de error.
- **DAOs**: se prueban contra una base **SQLite en memoria** con GORM, verificando que las consultas reales funcionen sin depender de MySQL.
- **Mocks**: definidos en `services/mocks_test.go` y `controllers/mocks_test.go`. Interfaces de DAO como costuras (seams) para inyeccion de dependencias.

### Tests de dominio

| Archivo         | Casos cubiertos                                                                                                                 |
| --------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `event_test.go` | `CurrentSalePhase` - 9 subtests que cubren las 4 fases de venta mas casos borde con fechas nil y valores exactos en los limites |

### Tests de utils

| Archivo            | Casos cubiertos                                                               |
| ------------------ | ----------------------------------------------------------------------------- |
| `password_test.go` | HashPassword, CheckPassword (correcta, incorrecta, vacia), HashPasswordUnique |
| `jwt_test.go`      | GenerateJWT, ValidateToken (valido, firma invalida, malformado, vacio)        |

### Tests de acceso a datos (DAO)

| Archivo              | Casos cubiertos                                                                                                   |
| -------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `user_dao_test.go`   | Create (exito, email duplicado), FindByID y FindByEmail (encontrado y no encontrado)                              |
| `event_dao_test.go`  | Create, FindByID, FindAll (filtros por categoria, fechas y precio; oculta eventos pasados), Update, Delete, IncrementTicketsSold, DecrementTicketsSold (sin bajar de 0) |
| `ticket_dao_test.go` | Create, FindByID, FindByUserID, FindActiveByEvent, CountActiveByEvent, CancelByEvent, Save, WithTransaction       |
| `connection_test.go` | getEnvOrDefault (valor presente, ausente y vacio)                                                                 |

### Tests de clients, logger y configuracion

| Archivo                      | Casos cubiertos                                                                                              |
| ---------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `email_client_test.go`       | Seleccion de proveedor (default, log, smtp, desconocido), envio de confirmacion de compra, cancelacion y transferencia (log y SMTP con host invalido) |
| `logger_test.go`             | Niveles Info, Warn, Error, Debug y Fatal, salida a stdout/stderr, formato JSON                                |
| `main_test.go`               | parseOrigins, isOriginAllowed, headers de seguridad, CORS (preflight 204, origen permitido y no permitido)    |
| `gencerts_test.go`           | Generacion de clave y certificado, escritura de archivos PEM validos                                          |

### Tests de servicios

| Archivo                  | Casos cubiertos                                                                                                                                                                                                                                                                      |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `auth_service_test.go`   | Register (exito, password corta, email duplicado, error de conexion en busqueda, error de creacion), Login (exito, password incorrecta, email desconocido), GenerateToken, NotFound error                                                                                            |
| `event_service_test.go`  | GetAll, GetByID (exito, no encontrado), Create (valido, titulo vacio, capacidad cero, fecha pasada, presale activo, presale invalido, DAO error, campos opcionales), Cancel (activo, ya cancelado, no encontrado, error en Update, error en CancelByEvent), Update (exito, cancelado, no encontrado, multiples campos, Date/Duration, presale activo/inactivo, conserva el codigo de preventa si no se envia uno nuevo, fechas invalidas, DAO error), validatePresaleConfig (todas las ramas) |
| `ticket_service_test.go` | Purchase (exito, quantity 0 default, evento cancelado, sin capacidad, count error, create tx error, fase no abierta, no encontrado), PurchasePresale (codigo correcto, sin codigo, codigo incorrecto), CancelTicket (propio, ajeno, ya cancelado, save tx error, no encontrado), Transfer (a otro usuario, a si mismo, save tx error, ticket no encontrado, no propio, ya cancelado, ya transferido, target no encontrado), GetByUser, emails (user/event not found, send fails), toTicketInfo |
| `report_service_test.go` | GetEventReport (exito, no encontrado, count error, findActive error ignorado, capacity 0, findUser error ignorado), GetGlobalReport (exito, FindAll error, count error)                                                                                                            |

### Tests de controladores

| Archivo                      | Casos cubiertos                                                                                                                                                                                             |
| ---------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `middleware_test.go`         | AuthRequired (sin header, formato invalido, token malformado), AdminRequired (admin pasa, client 403, sin rol 403)                                                                                          |
| `auth_controller_test.go`    | Register (201, 400 campos faltantes, 400 email duplicado, 500 token error, 500 server error), Login (200 con token, 401 credenciales invalidas, 400 campo faltante), userFriendlyError (3 casos), helpers (isNotFound, optionalString, timePtrToString) |
| `event_controller_test.go`   | GetAll (200, filtros category/date/price, filtros invalidos, error servicio), GetByID (200, 404, 400, 500), Create (201, 201 evento gratis con precio 0, 400 precio negativo, 400 campos faltantes, 400 fecha invalida, 400 service error, 500), Update (200, 404, 400, 400 cancelled, 400 json invalido, 500), Delete (204, 404, 400), GetSaleStatus (200 presale, 400 invalid id, 404, cancelled, sold out, no presale, 500), salePhaseMessage (todas las fases) |
| `ticket_controller_test.go`  | Purchase (201, 400 evento cancelado, 400 sin capacidad, 400 no encontrado, 400 campo faltante, 400 venta no abierta, 400 falta codigo de preventa, 400 codigo de preventa invalido, 500), GetMyTickets (200, 500, empty), Cancel (204, 404, 400, 500), Transfer (200, 400 campo faltante, 400 invalid id, 404, 500), toTicketResponse (con/sin evento) |
| `admin_controller_test.go`   | GetReports (200, 500), GetEventReport (200, 404, 400, 500), toGlobalReportResponse, toEventReportResponse                                                                                                  |

---

## Bonus Track

**Preventa con codigo de acceso y fases de venta**: un administrador puede activar una preventa para un evento, definiendo un codigo de acceso, la fecha de inicio de la preventa y la fecha de inicio de la venta general (inicio de preventa < venta general < fecha del evento). Segun el momento en que se intenta comprar, el evento esta en una de estas fases:

| Fase           | Quien puede comprar                                   |
| -------------- | ----------------------------------------------------- |
| `not_yet_open` | Nadie (la preventa todavia no empezo)                 |
| `presale`      | Solo quien ingresa el codigo de preventa correcto     |
| `public`       | Cualquier usuario registrado                          |
| `no_presale`   | Cualquier usuario registrado (evento sin preventa)    |

Donde esta implementado:

- **Dominio**: `backend/domain/event.go` — campos `PresaleActive`, `PresaleCode`, `PresaleStartDate`, `GeneralSaleDate` y el metodo `CurrentSalePhase(now)`, que calcula la fase actual.
- **Service**: `TicketService.Purchase` valida la fase y el codigo antes de vender (`ErrSalesNotOpen`, `ErrPresaleCodeRequired`, `ErrInvalidPresaleCode`). `EventService` valida la configuracion de la preventa al crear o editar un evento.
- **Endpoint**: `GET /api/events/:id/sale-status` devuelve la fase de venta actual del evento.
- **Frontend**: el detalle del evento pide el codigo cuando el evento esta en fase `presale`, y el panel de administracion permite configurar la preventa en el formulario de eventos.

---

## Autores

BRUA, Jonathan — HERNANDEZ, Juan — LINDON, Maria Victoria — TERRERA, Athina

Proyecto integrador 2026 - Universidad Catolica de Cordoba (UCC)
