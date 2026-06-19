TRUNCATE TABLE
    soldat.verletzung,
    soldat.ausruestung,
    soldat.soldat
RESTART IDENTITY CASCADE;

INSERT INTO soldat.soldat
    (version, vorname, nachname, geburtsdatum, geschlecht, rang, username, erzeugt, aktualisiert)
VALUES
    (0, 'Eren', 'Jaeger', '2000-01-01', 'MAENNLICH', 'SOLDAT', 'eren-test', NOW(), NOW()),
    (0, 'Mikasa', 'Ackermann', '2000-02-10', 'WEIBLICH', 'ELITE-SOLDAT', 'mikasa-test', NOW(), NOW()),
    (0, 'Armin', 'Arlert', '2000-11-03', 'MAENNLICH', 'REKRUT', 'armin-test', NOW(), NOW()),
    (0, 'Levi', 'Ackermann', '1985-12-25', 'MAENNLICH', 'CAPTAIN', 'levi-test', NOW(), NOW()),
    (0, 'Historia', 'Reiss', '2000-01-15', 'WEIBLICH', 'SOLDAT', 'historia-test', NOW(), NOW());

INSERT INTO soldat.ausruestung
    (waffe, seriennummer, soldat_id)
VALUES
    ('ODM_GEAR', 'AOT-SEED-001', (SELECT id FROM soldat.soldat WHERE username = 'eren-test')),
    ('Klinge', 'AOT-SEED-002', (SELECT id FROM soldat.soldat WHERE username = 'mikasa-test')),
    ('ODM_GEAR', 'AOT-SEED-003', (SELECT id FROM soldat.soldat WHERE username = 'armin-test')),
    ('Klinge', 'AOT-SEED-004', (SELECT id FROM soldat.soldat WHERE username = 'levi-test')),
    ('Schrotflinte', 'AOT-SEED-005', (SELECT id FROM soldat.soldat WHERE username = 'historia-test'));

INSERT INTO soldat.verletzung
    (verletzungsbezeichnung, behandelt, schweregrad, verletzungsdatum, soldat_id)
VALUES
    ('Schnittverletzung', true, 'LEICHT', '2025-01-10', (SELECT id FROM soldat.soldat WHERE username = 'eren-test')),
    ('Prellung', false, 'MITTEL', '2025-02-12', (SELECT id FROM soldat.soldat WHERE username = 'mikasa-test')),
    ('Knochenbruch', true, 'SCHWER', '2025-03-20', (SELECT id FROM soldat.soldat WHERE username = 'levi-test'));
