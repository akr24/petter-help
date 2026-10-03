-- Accounts for seekers and purveyors.
--
-- password_hash holds a PHC-encoded argon2id string ("$argon2id$v=19$m=...$salt$hash"),
-- so the salt and cost parameters are stored alongside the hash.
-- Emails are stored lower-cased by the application; the unique index enforces it.

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL CHECK (role IN ('seeker', 'purveyor')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_email_lower CHECK (email = lower(email))
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users (email);
