USE ceibo_db;

-- Admin users (FK reference for events.created_by_id)
-- Contraseñas personalizadas por usuario ([firstname]123)
INSERT INTO users (name, email, password_hash, role) VALUES
    ('Carlos Méndez',   'carlos@ceibo.com',    '$2a$12$uKj2TFzunaEH5NPctPCRnuIvIiG8b87P5R.Gzt0UWooqeizDoDo2a', 'admin'),
    ('Ana Lucía Rivas', 'ana@ceibo.com',       '$2a$12$ujbVBRdHeh5uReTtcsjiL.o2m2iKPx9TbVw5MEc58cqx5dSpTIIRK', 'admin'),
    ('Pedro Castillo',  'pedro@ceibo.com',     '$2a$12$SfFDw7c4bOWILelUe00dRulli6QflbSIWZmYG7K/0avuY/XTJnj0.', 'admin');

-- Client users
INSERT INTO users (name, email, password_hash, role) VALUES
    ('Sofía Martínez',  'sofia@gmail.com',     '$2a$12$86zzn/cua1JGhjUzTRZ2U.W3egchiuzZsnbAmBWFPirhGGu14JE3G', 'client'),
    ('Lautaro Gómez',   'lautaro@hotmail.com', '$2a$12$OVJlm2F0Wpc3C3opQnwAHuf2sR9962OwZMWQjMSViTQd/QSm770oO', 'client'),
    ('Camila Rodríguez','camila@yahoo.com',    '$2a$12$NvPB6MJak/aZjZKNqZ2QYOSZU84J1lcpVM.bc4cOu0.NVGizo8VIq', 'client'),
    ('Facundo Díaz',    'facundo@gmail.com',   '$2a$12$IKvwDbCtRPubqjdwCODfkOejiyW16XYOPrYgIhHJIu8UVKb.9hMmG', 'client'),
    ('Lucía Fernández', 'lucia@outlook.com',   '$2a$12$JuOMDDo8585ofwM41wk8j.qUkP6RC.pQdmoQM.wej1EYf5hatgtzy', 'client'),
    ('Mateo Álvarez',   'mateo@gmail.com',     '$2a$12$p7mVk5Qea4lDK97EkKmno.eJ/OqunfaugCBbOu8IgfaR7BX3T75Lm', 'client'),
    ('Valentina Torres','valentina@live.com',  '$2a$12$piCrLYOuuSHo17AgW5.EIeQurJz5gdOQE3CniHD.xGFsv77.rQCE.', 'client'),
    ('Santiago López',  'santiago@gmail.com',  '$2a$12$1qYnIYmvKilRCZv489EDv.7sSNnc.AqyNRrImgdiJ3MErR7Ii.EqS', 'client'),
    ('Florencia Acosta','florencia@yahoo.com', '$2a$12$k/ZtydObejACO30ZxaaPU.XkMIjbcs/oStNy9W5eXISe8vrnEAgjy', 'client'),
    ('Agustín Pereyra', 'agustin@gmail.com',   '$2a$12$oKBf3LzxrhAFQzld7twq9uv1FvJFAYCBqiO13RPYeENxCjtILAtR6', 'client');

-- Events (23)
INSERT INTO events (title, description, image_url, category, location, event_date, duration_minutes, capacity, tickets_sold, price, status, presale_active, presale_code, presale_start_date, general_sale_date, created_by_id)
VALUES
    ('Peña folclórica Los Amigos',           'Noche de zambas y chacareras con artistas invitados.',             '/images/pena-amigos.jpg',              'Peña',         'Peña Los Amigos, Palermo, CABA',           '2026-08-15 20:00:00', 300, 200,   165,   25000.00,'active',   0, NULL,                        NULL,                       NULL,                       1),
    ('Peña de la Tradición Salteña',         'Música y danzas típicas del norte argentino.',                     '/images/tradicion-saltena.jpg',        'Peña',         'Centro Cultural, Salta',                   '2026-09-10 21:00:00', 360, 350,   280,   30000.00,'active',   0, NULL,                        NULL,                       NULL,                       1),
    ('Peña del Carnavalito',                 'Carnavalito, erke y caja coplera en vivo.',                        '/images/carnavalito.jpg',              'Peña',         'Plaza 9 de Julio, Jujuy',                  '2026-07-22 19:00:00', 240, 500,   410,   35000.00,'active',   1, 'CARNAVAL2026',             '2026-06-01 00:00:00',     '2026-07-01 00:00:00',     2),
    ('Peña del Chamamé',                     'Acordeón, guitarra y sapucay en el litoral.',                     '/images/chamame.jpg',                  'Peña',         'Club Social, Corrientes',                  '2026-06-30 21:00:00', 300, 300,   300,   20000.00,'sold_out', 0, NULL,                        NULL,                       NULL,                       2),
    ('Peña de la Chacarera',                 'Bombo legüero y violín para bailar hasta el amanecer.',            '/images/chacarera.jpg',                'Peña',         'Peña El Sauce, Santiago del Estero',       '2026-10-05 21:00:00', 360, 150,   95,    15000.00,'active',   0, NULL,                        NULL,                       NULL,                       3),
    ('Peña de la Cosecha',                   'Festival de la vendimia con artistas cuyanos.',                   '/images/cosecha.jpg',                  'Peña',         'Plaza Independencia, Mendoza',             '2026-11-20 20:00:00', 420, 800,   520,   55000.00,'active',   1, 'CUYO2026',                 '2026-08-01 00:00:00',     '2026-09-01 00:00:00',     1),
    ('Peña del Poncho',                      'Encuentro de artesanos y músicos populares.',                     '/images/poncho.jpg',                   'Peña',         'Predio Ferial, La Rioja',                  '2026-09-25 18:00:00', 480, 2000,  2000,  45000.00,'sold_out', 0, NULL,                        NULL,                       NULL,                       3),
    ('Peña de la Milonga',                   'Milonga campera con guitarra y bandoneón.',                       '/images/milonga.jpg',                  'Peña',         'Sociedad Rural, Rosario',                  '2026-05-15 20:00:00', 240, 200,   80,    20000.00,'cancelled',0, NULL,                        NULL,                       NULL,                       2),
    ('Peña de la Cueca Cuyana',              'Música cuyana con guitarra, bombo y acordeón.',                   '/images/cueca-cuyana.jpg',             'Peña',         'Teatro Griego, San Juan',                  '2026-12-01 20:00:00', 240, 400,   0,     25000.00,'presale',  1, 'CUECA2026',                '2026-11-01 00:00:00',     '2026-11-15 00:00:00',     1),
    ('Peña Gaucha del Talar',                'Jineteadas, danzas y fogón criollo.',                             '/images/gaucha.jpg',                   'Peña',         'Parque Gaucho, Tandil',                    '2026-07-12 10:00:00', 600, 1500,  1100,  20000.00,'active',   0, NULL,                        NULL,                       NULL,                       3),
    ('Peña del Malambo',                     'Espectáculo de malambo con ballet folclórico.',                   '/images/malambo.jpg',                  'Peña',         'Teatro San Martín, Córdoba',               '2026-08-28 20:00:00', 180, 250,   250,   35000.00,'sold_out', 0, NULL,                        NULL,                       NULL,                       2),
    ('Peña de la Tonada',                    'Tonadas cuyanas al pie del cerro.',                              '/images/tonada.jpg',                   'Peña',         'Anfiteatro, San Luis',                     '2026-06-10 19:00:00', 270, 600,   0,     20000.00,'active',   0, NULL,                        NULL,                       NULL,                       1),
    ('Peña del Fogón Patagónico',            'Música sureña con contrabajo y acordeón a orillas del lago.',     '/images/fogon-patagonico.jpg',         'Peña',         'Costanera, Bariloche',                     '2026-10-30 21:00:00', 300, 180,   145,   40000.00,'active',   1, 'PATAGONIA2026',            '2026-09-01 00:00:00',     '2026-10-01 00:00:00',     2),
    ('Peña de la Lechiguana',                'Candombe y milonga con percusión afro-uruguaya.',                  '/images/lechiguana.jpg',               'Peña',         'Teatro Solís, Montevideo',                 '2026-11-15 22:00:00', 300, 350,   350,   20000.00,'sold_out', 0, NULL,                        NULL,                       NULL,                       3),
    ('Peña del Buen Mate',                   'Tarde de mate, torta frita y canciones de la tierra.',            '/images/buen-mate.jpg',                'Peña',         'Museo de la Yerba Mate, Misiones',         '2026-09-18 16:00:00', 240, 120,   68,    15000.00,'active',   0, NULL,                        NULL,                       NULL,                       2),
    ('Peña del Gato y el Escondido',         'Baile tradicional con parejas de danzas folclóricas.',             '/images/gato-escondido.jpg',           'Peña',         'Peña La Vieja Estación, Tucumán',          '2026-12-31 21:00:00', 360, 300,   220,   60000.00,'active',   1, 'FOLCLORE2026',             '2026-11-01 00:00:00',     '2026-12-01 00:00:00',     1),
    ('Peña del Locro',                       'Peña con locro, empanadas y vino patagónico.',                    '/images/locro.jpg',                    'Peña',         'Salón Comunitario, Neuquén',               '2026-07-30 12:00:00', 360, 100,   100,   18000.00,'sold_out', 0, NULL,                        NULL,                       NULL,                       3),
    ('Peña de la Baguala',                   'Cantos ancestrales del altiplano con caja y quena.',              '/images/baguala.jpg',                  'Peña',         'Cerro de los Siete Colores, Purmamarca',   '2027-01-15 18:00:00', 240, 300,   0,     25000.00,'presale',  1, 'BAGUALA2027',              '2026-12-01 00:00:00',     '2026-12-15 00:00:00',     2),
    ('Peña Federal de la Guitarra',          'Guitarra, bombo y violin de todos los rincones del país.',        '/images/guitarra.jpg',                 'Peña',         'Anfiteatro Municipal, Paraná',             '2026-08-05 20:00:00', 360, 1000,  0,     0.00,   'active',   0, NULL,                        NULL,                       NULL,                       1),
    ('Peña de la Zamba',                     'Noche entera dedicada a la zamba, la madre de nuestras danzas.',   '/images/zamba.jpg',                    'Peña',         'Peña El Cardón, Salta',                    '2027-02-20 21:00:00', 360, 80,    10,    70000.00,'active',   1, 'ZAMBA2027',                '2027-01-01 00:00:00',     '2027-01-15 00:00:00',     3),
    -- Festival events (21-23)
    ('Festival Nacional del Folclore',       'Tres días de música, danza y gastronomía típica de todas las provincias.', '/images/festival-folclore.jpg',       'Festival',     'Predio Ferial, Cosquín, Córdoba',          '2027-01-18 18:00:00', 4320, 50000, 12500, 100000.00,'active',   1, 'COSQUIN2027',             '2026-10-01 00:00:00',     '2026-12-01 00:00:00',     1),
    ('Festival de la Cerveza Artesanal',     'Música en vivo, food trucks y más de 50 cervecerías artesanales.',  '/images/cerveza-fest.jpg',             'Festival',     'Parque Norte, Córdoba',                     '2026-11-08 16:00:00', 2880, 30000, 8000,  80000.00,'active',   1, 'CERVEZA2026',             '2026-08-01 00:00:00',     '2026-10-01 00:00:00',     2),
    ('Festival del Día de la Primavera',     'Jornada al aire libre con bandas nacionales, ferias y actividades.',  '/images/primavera-fest.jpg',           'Festival',     'Parque 3 de Febrero, CABA',                '2026-09-21 12:00:00', 720, 20000, 15000, 90000.00, 'active',   0, NULL,                        NULL,                       NULL,                       3);

-- ============================================================
-- Tickets (35)
-- Sofía (4): 5 tickets  |  Lautaro (5): 3  |  Camila (6): 4
-- Facundo (7): 4        |  Lucía (8): 4     |  Mateo (9): 5
-- Valentina (10): 4     |  Santiago (11): 4 |  Florencia (12): 1
-- Agustín (13): 1
-- ============================================================
INSERT INTO tickets (user_id, event_id, status, purchase_price, purchased_at, cancelled_at, transferred_at, transferred_to_id)
VALUES
    -- Sofía (ID 4)
    (4,  1,  'active',     25000.00, '2026-07-01 14:30:00', NULL, NULL, NULL),
    (4,  4,  'active',     20000.00, '2026-05-10 10:00:00', NULL, NULL, NULL),
    (4,  9,  'active',     25000.00, '2026-10-15 09:00:00', NULL, NULL, NULL),
    (4, 13,  'active',     40000.00, '2026-09-20 18:45:00', NULL, NULL, NULL),
    (4, 21,  'active',    100000.00, '2026-11-01 10:00:00', NULL, NULL, NULL),
    -- Lautaro (ID 5)
    (5,  3,  'active',     35000.00, '2026-05-28 11:00:00', NULL, NULL, NULL),
    (5, 16,  'active',     60000.00, '2026-10-05 16:30:00', NULL, NULL, NULL),
    (5, 22,  'active',     80000.00, '2026-09-15 14:00:00', NULL, NULL, NULL),
    -- Camila (ID 6)
    (6,  6,  'active',     55000.00, '2026-07-10 08:00:00', NULL, NULL, NULL),
    (6, 17,  'active',     18000.00, '2026-06-15 12:00:00', NULL, NULL, NULL),
    (6, 20,  'active',     70000.00, '2026-12-20 20:00:00', NULL, NULL, NULL),
    (6, 21,  'active',    100000.00, '2026-11-20 14:00:00', NULL, NULL, NULL),
    -- Facundo (ID 7)
    (7,  7,  'active',     45000.00, '2026-08-01 15:00:00', NULL, NULL, NULL),
    (7, 11,  'active',     35000.00, '2026-07-20 09:30:00', NULL, NULL, NULL),
    (7, 14,  'active',     20000.00, '2026-10-01 11:00:00', NULL, NULL, NULL),
    (7, 23,  'active',     90000.00, '2026-08-10 09:00:00', NULL, NULL, NULL),
    -- Lucía (ID 8)
    (8,  5,  'active',     15000.00, '2026-09-01 10:00:00', NULL, NULL, NULL),
    (8, 10,  'active',     20000.00, '2026-05-30 13:00:00', NULL, NULL, NULL),
    (8,  8,  'cancelled',   0.00, '2026-04-20 17:00:00', '2026-05-01 09:00:00', NULL, NULL),
    (8, 22,  'active',     80000.00, '2026-09-01 16:00:00', NULL, NULL, NULL),
    -- Mateo (ID 9)
    (9,  1,  'active',     25000.00, '2026-07-15 20:00:00', NULL, NULL, NULL),
    (9,  6,  'active',     55000.00, '2026-08-10 14:00:00', NULL, NULL, NULL),
    (9, 12,  'active',     20000.00, '2026-05-05 10:30:00', NULL, NULL, NULL),
    (9, 20,  'active',     70000.00, '2026-12-22 19:00:00', NULL, NULL, NULL),
    (9, 23,  'active',     90000.00, '2026-08-25 11:00:00', NULL, NULL, NULL),
    -- Valentina (ID 10)
    (10, 15, 'active',     15000.00, '2026-08-20 11:15:00', NULL, NULL, NULL),
    (10, 19, 'active',       0.00, '2026-07-01 09:00:00', NULL, NULL, NULL),
    (10, 20, 'active',     70000.00, '2026-12-18 08:30:00', NULL, NULL, NULL),
    (10, 21, 'active',    100000.00, '2026-11-10 12:00:00', NULL, NULL, NULL),
    -- Santiago (ID 11)
    (11,  2,  'active',     30000.00, '2026-08-05 16:00:00', NULL, NULL, NULL),
    (11,  4,  'transferred',20000.00, '2026-05-12 18:00:00', NULL, '2026-06-01 12:00:00', 8),
    (11, 10,  'active',     20000.00, '2026-06-10 10:00:00', NULL, NULL, NULL),
    (11, 22,  'active',     80000.00, '2026-10-01 10:00:00', NULL, NULL, NULL),
    -- Florencia (ID 12)
    (12, 21,  'active',    100000.00, '2026-12-01 10:00:00', NULL, NULL, NULL),
    -- Agustín (ID 13)
    (13, 23,  'active',     90000.00, '2026-09-01 09:00:00', NULL, NULL, NULL);
