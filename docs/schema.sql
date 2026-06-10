CREATE DATABASE IF NOT EXISTS ceibo_db
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

USE ceibo_db;

CREATE TABLE IF NOT EXISTS users (
id              INT UNSIGNED NOT NULL AUTO_INCREMENT,
name            VARCHAR(150) NOT NULL,
email           VARCHAR(255) NOT NULL,
password_hash   VARCHAR(255) NOT NULL, -- SHA-256 / bcrypt; never plain-text
role            ENUM('client', 'admin') NOT NULL DEFAULT 'client',
created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
deleted_at      DATETIME NULL DEFAULT NULL, -- soft-delete compatible with GORM

```
PRIMARY KEY (id),
UNIQUE KEY uq_users_email (email),
INDEX idx_users_deleted_at (deleted_at)
```

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS events (
id                  INT UNSIGNED NOT NULL AUTO_INCREMENT,
title               VARCHAR(200) NOT NULL,
description         TEXT NULL,
image_url           VARCHAR(500) NULL,
category            VARCHAR(100) NULL,
location            VARCHAR(300) NULL,
event_date          DATETIME NOT NULL,
duration_minutes    INT NOT NULL DEFAULT 0,
capacity            INT NOT NULL DEFAULT 0,
tickets_sold        INT NOT NULL DEFAULT 0,
price               DECIMAL(10,2) NOT NULL DEFAULT 0.00,

```
-- Event status and presale configuration
status              ENUM('active', 'presale', 'sold_out', 'cancelled') NOT NULL DEFAULT 'active',
presale_active      BOOLEAN NOT NULL DEFAULT FALSE,
presale_code        VARCHAR(100) NULL,
presale_start_date  DATETIME NULL,
general_sale_date   DATETIME NULL,

created_by_id       INT UNSIGNED NOT NULL, -- FK -> users(id), administrator who created the event
created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
deleted_at          DATETIME NULL DEFAULT NULL, -- soft-delete compatible with GORM

PRIMARY KEY (id),

INDEX idx_events_status_date (status, event_date),
INDEX idx_events_event_date (event_date),
INDEX idx_events_category (category),
INDEX idx_events_deleted_at (deleted_at),
INDEX idx_events_created_by (created_by_id),

CONSTRAINT fk_events_created_by
    FOREIGN KEY (created_by_id)
    REFERENCES users (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

CONSTRAINT chk_events_capacity
    CHECK (tickets_sold >= 0 AND tickets_sold <= capacity)
```

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tickets (
id                  INT UNSIGNED NOT NULL AUTO_INCREMENT,
user_id             INT UNSIGNED NOT NULL, -- Current holder / buyer
event_id            INT UNSIGNED NOT NULL,
status              ENUM('active', 'cancelled', 'transferred') NOT NULL DEFAULT 'active',
purchase_price      DECIMAL(10,2) NOT NULL DEFAULT 0.00, -- Price at time of purchase
purchased_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
cancelled_at        DATETIME NULL DEFAULT NULL,
transferred_at      DATETIME NULL DEFAULT NULL,
transferred_to_id   INT UNSIGNED NULL DEFAULT NULL, -- Optional recipient of transfer
created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

```
PRIMARY KEY (id),

INDEX idx_tickets_user_id (user_id),
INDEX idx_tickets_event_id (event_id),
INDEX idx_tickets_event_status (event_id, status),
INDEX idx_tickets_status (status),
INDEX idx_tickets_transferred_to (transferred_to_id),

CONSTRAINT fk_tickets_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

CONSTRAINT fk_tickets_event
    FOREIGN KEY (event_id)
    REFERENCES events (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

CONSTRAINT fk_tickets_transferred_to
    FOREIGN KEY (transferred_to_id)
    REFERENCES users (id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
```

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
