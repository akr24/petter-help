-- Initial schema. Runs automatically the first time the Postgres container
-- initialises its data volume (see docker-compose.db.yml).
--
-- The dogs table is deliberately simple: the real dog categories and the
-- questionnaire traits used for matching are still open questions.

CREATE TABLE IF NOT EXISTS dogs (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    breed       TEXT NOT NULL,
    size        TEXT NOT NULL,
    personality TEXT NOT NULL DEFAULT '',
    needs       TEXT NOT NULL DEFAULT '',
    purveyor    TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Sample listings so the API returns something on a fresh database.
INSERT INTO dogs (name, breed, size, personality, needs, purveyor) VALUES
    ('Biscuit', 'Beagle',           'medium', 'curious, vocal, loves walks',         'daily exercise, secure yard',   'Cambridge Animal Rescue'),
    ('Maple',   'Labrador mix',     'large',  'gentle, patient, great with kids',    'joint supplements',             'MSPCA Boston'),
    ('Pixel',   'Chihuahua',        'small',  'shy at first, bonded to one person',  'quiet home, no small children', 'Private rehome');
