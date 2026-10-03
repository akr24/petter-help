-- Seeker lifestyle profile: one row per seeker, keyed by users.id.
--
-- These are the questionnaire answers the matching algorithm will score
-- against dog listings. Every column maps to something that changes which
-- dog is a good fit, not just which dog the seeker likes the look of.
-- Values are constrained text rather than Postgres enums so adding an
-- option later is an ALTER on a CHECK, not a type migration.
--
-- Column                 Why it matters for matching
-- home_type / has_yard   space and outdoor access (large or high-energy dogs)
-- children               tolerance for kids; "young" is under 6
-- has_dogs / has_cats    dog must get along with resident animals
-- activity_level         seeker's own exercise habits vs the dog's energy
-- hours_alone            separation tolerance; a big one for rescues
-- experience             first-timers do better with easy-going dogs
-- training_commitment    willingness to work through behaviour or puppy training
-- grooming_commitment    long coats and heavy shedders need time or money
-- needs_hypoallergenic   hard filter for allergy households
-- size_preferences       empty array means no preference
-- age_preferences        empty array means no preference
-- notes                  free text the purveyor can read; not scored

CREATE TABLE IF NOT EXISTS seeker_profiles (
    user_id              BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,

    home_type            TEXT NOT NULL CHECK (home_type IN ('apartment', 'house', 'other')),
    has_yard             BOOLEAN NOT NULL DEFAULT false,
    children             TEXT NOT NULL CHECK (children IN ('none', 'young', 'older')),
    has_dogs             BOOLEAN NOT NULL DEFAULT false,
    has_cats             BOOLEAN NOT NULL DEFAULT false,

    activity_level       TEXT NOT NULL CHECK (activity_level IN ('low', 'moderate', 'high')),
    hours_alone          SMALLINT NOT NULL CHECK (hours_alone BETWEEN 0 AND 24),
    experience           TEXT NOT NULL CHECK (experience IN ('first_time', 'some', 'experienced')),
    training_commitment  TEXT NOT NULL CHECK (training_commitment IN ('low', 'moderate', 'high')),
    grooming_commitment  TEXT NOT NULL CHECK (grooming_commitment IN ('low', 'moderate', 'high')),
    needs_hypoallergenic BOOLEAN NOT NULL DEFAULT false,

    size_preferences     TEXT[] NOT NULL DEFAULT '{}'
                         CHECK (size_preferences <@ ARRAY['small', 'medium', 'large']),
    age_preferences      TEXT[] NOT NULL DEFAULT '{}'
                         CHECK (age_preferences <@ ARRAY['puppy', 'adult', 'senior']),

    notes                TEXT NOT NULL DEFAULT '',

    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);
