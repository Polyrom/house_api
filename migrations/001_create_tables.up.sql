-- create users table
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(255) NOT NULL CHECK(role IN ('client', 'moderator')),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- create tokens table
CREATE TABLE IF NOT EXISTS tokens (
  user_id UUID NOT NULL UNIQUE REFERENCES users(id),
  token UUID NOT NULL,
  expires_at TIMESTAMP NOT NULL
);
-- create houses table
CREATE TABLE IF NOT EXISTS houses (
  id SERIAL PRIMARY KEY CHECK (id >= 1),
  address VARCHAR(1000) NOT NULL,
  year INTEGER NOT NULL CHECK (id >= 0),
  developer VARCHAR(200),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_id UNIQUE (id)
);
-- create flats table
CREATE TABLE IF NOT EXISTS flats (
  id SERIAL PRIMARY KEY CHECK (id >= 1),
  house_id INTEGER NOT NULL REFERENCES houses(id),
  price INTEGER NOT NULL,
  rooms INTEGER NOT NULL CHECK (id >= 1),
  moderator UUID REFERENCES users(id),
  status VARCHAR(50) NOT NULL CHECK(
    status IN (
      'created',
      'approved',
      'declined',
      'on moderation'
    )
  ) DEFAULT 'created'
);
-- create function to update houses update_at any time a new flat is added
CREATE OR REPLACE FUNCTION update_house_updated_at() RETURNS TRIGGER AS $$ BEGIN
UPDATE houses
SET update_at = CURRENT_TIMESTAMP
WHERE id = NEW.house_id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- attach the trigger to the flats table
CREATE TRIGGER update_house_updated_at_trigger
AFTER
INSERT ON flats FOR EACH ROW EXECUTE PROCEDURE update_house_updated_at();