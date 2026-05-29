-- =============================================================
-- Sistema de Gestión de Eventos y Entradas - Grupo 5
-- Facultad de Ingeniería - Desarrollo de Software 2026
-- =============================================================

CREATE DATABASE IF NOT EXISTS ticketek_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE ticketek_db;

-- =============================================================
-- TABLA: users
-- =============================================================
CREATE TABLE IF NOT EXISTS users (
    id            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    name          VARCHAR(150)     NOT NULL,
    email         VARCHAR(255)     NOT NULL,
    password_hash VARCHAR(255)     NOT NULL,                       -- SHA-256 / bcrypt; nunca plain-text
    role          ENUM('client', 'admin') NOT NULL DEFAULT 'client',
    created_at    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME         NULL DEFAULT NULL,              -- soft-delete (compatible con GORM)

    PRIMARY KEY (id),
    UNIQUE KEY uq_users_email (email),
    INDEX idx_users_role (role),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================================
-- TABLA: events
-- =============================================================
CREATE TABLE IF NOT EXISTS events (
    id                  BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    title               VARCHAR(200)     NOT NULL,
    description         TEXT             NULL,
    image_url           VARCHAR(500)     NULL,
    category            VARCHAR(100)     NULL,
    location            VARCHAR(300)     NULL,
    event_date          DATETIME         NOT NULL,
    duration_minutes    INT              NOT NULL DEFAULT 0,
    capacity            INT              NOT NULL DEFAULT 0,
    tickets_sold        INT              NOT NULL DEFAULT 0,
    price               DECIMAL(10,2)   NOT NULL DEFAULT 0.00,

    -- Funcionalidad Extra – Preventa / Reserva anticipada (Grupo 5)
    status              ENUM('active', 'presale', 'sold_out', 'cancelled') NOT NULL DEFAULT 'active',
    presale_code        VARCHAR(100)     NULL,                    -- Código de acceso anticipado
    presale_start_date  DATETIME         NULL,                    -- Inicio del período de preventa
    general_sale_date   DATETIME         NULL,                    -- Fecha apertura venta general

    created_by_id       BIGINT UNSIGNED  NOT NULL,               -- FK → users (administrador)
    created_at          DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME         NULL DEFAULT NULL,       -- soft-delete

    PRIMARY KEY (id),
    INDEX idx_events_status      (status),
    INDEX idx_events_event_date  (event_date),
    INDEX idx_events_category    (category),
    INDEX idx_events_deleted_at  (deleted_at),
    INDEX idx_events_created_by  (created_by_id),

    CONSTRAINT fk_events_created_by
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================================
-- TABLA: tickets
-- =============================================================
CREATE TABLE IF NOT EXISTS tickets (
    id              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
    user_id         BIGINT UNSIGNED  NOT NULL,                    -- Titular actual
    event_id        BIGINT UNSIGNED  NOT NULL,
    status          ENUM('active', 'cancelled', 'transferred') NOT NULL DEFAULT 'active',
    purchase_price  DECIMAL(10,2)   NOT NULL DEFAULT 0.00,        -- Precio al momento de la compra
    purchased_at    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cancelled_at    DATETIME         NULL DEFAULT NULL,
    transferred_at  DATETIME         NULL DEFAULT NULL,
    transferred_to_id BIGINT UNSIGNED NULL DEFAULT NULL,          -- A quién se transfirió
    created_at      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    INDEX idx_tickets_user_id    (user_id),
    INDEX idx_tickets_event_id   (event_id),
    INDEX idx_tickets_status     (status),
    INDEX idx_tickets_transferred_to (transferred_to_id),

    CONSTRAINT fk_tickets_user
        FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_tickets_event
        FOREIGN KEY (event_id) REFERENCES events (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_tickets_transferred_to
        FOREIGN KEY (transferred_to_id) REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================================
-- DATOS INICIALES (seeds opcionales para desarrollo)
-- =============================================================

-- Administrador por defecto
-- Contraseña: "admin1234" hasheada con SHA-256
INSERT INTO users (name, email, password_hash, role) VALUES
('Admin Sistema', 'admin@ticketek.com',
 '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', -- SHA-256 de "1234" (cambiar en producción)
 'admin');

-- Cliente de ejemplo
INSERT INTO users (name, email, password_hash, role) VALUES
('Cliente Demo', 'cliente@ticketek.com',
 '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4',
 'client');

-- Evento de ejemplo con preventa activa
INSERT INTO events (
    title, description, image_url, category, location,
    event_date, duration_minutes, capacity, tickets_sold, price,
    status, presale_code, presale_start_date, general_sale_date,
    created_by_id
) VALUES (
    'Lollapalooza Argentina 2026',
    'El festival de música más grande del año regresa al Hipódromo de San Isidro.',
    'https://example.com/lolla2026.jpg',
    'Festival',
    'Hipódromo de San Isidro, Buenos Aires',
    '2026-11-20 14:00:00',
    480,
    50000,
    0,
    15000.00,
    'presale',
    'PREVENTAVIP',
    '2026-06-01 10:00:00',
    '2026-07-01 10:00:00',
    1   -- created_by_id = Admin Sistema
);

-- Evento de ejemplo de venta general
INSERT INTO events (
    title, description, image_url, category, location,
    event_date, duration_minutes, capacity, tickets_sold, price,
    status, created_by_id
) VALUES (
    'Stand Up Comedy Night',
    'Una noche de comedia con los mejores humoristas del país.',
    'https://example.com/comedy.jpg',
    'Teatro',
    'Teatro Gran Rex, Buenos Aires',
    '2026-08-15 21:00:00',
    120,
    800,
    0,
    3500.00,
    'active',
    1
);
