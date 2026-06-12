# Ceibo Tickets — Backend

API REST para gestion de eventos y venta de entradas.

## Stack

Go 1.26 + Gin + GORM + MySQL + JWT + bcrypt

## Inicio rapido

```bash
cp .env.example .env   # configurar credenciales
go mod tidy
go run main.go         # requiere TLS (cert.pem + key.pem)
```

## Estructura

```
main.go         # Entry point, DI, rutas, middlewares
clients/        # Email client (log / SMTP)
controllers/    # Handlers HTTP + middleware + DTOs
dao/            # Interfaces GORM + implementaciones
domain/         # Entidades (User, Event, Ticket)
logger/         # Logger estructurado JSON
services/       # Logica de negocio + tests
```

Ver [docs/architecture.md](../docs/architecture.md) para detalle de diseno y patrones.

## Comandos utiles

```bash
go run main.go              # iniciar servidor
go test ./... -v -cover     # ejecutar tests
go build -o bin/server .    # compilar binario
```

## Convenciones

- Los controllers solo manejan HTTP y delegan a services
- Los services contienen la logica de negocio y errores sentinela
- Los DAOs se definen como interfaces para poder mockear en tests
- Las transacciones se manejan a nivel de service con `WithTransaction`

## Documentacion relacionada

| Que buscas? | Donde? |
|-------------|--------|
| Endpoints y ejemplos | [API Reference](../docs/api-reference.md) |
| Modelo de datos | [Base de Datos](../docs/database.md) |
| Patrones de diseno | [Arquitectura](../docs/architecture.md) |
| Tests y cobertura | [Testing](../docs/testing.md) |
