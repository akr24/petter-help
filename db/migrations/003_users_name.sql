-- Display name for accounts. Other user associations (lifestyle profiles,
-- purveyor details, listings) belong in their own tables keyed by users.id.

ALTER TABLE users ADD COLUMN IF NOT EXISTS name TEXT NOT NULL DEFAULT '';
