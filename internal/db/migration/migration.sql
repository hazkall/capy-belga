

CREATE SCHEMA IF NOT EXISTS capybelga;

SET search_path TO capybelga;


CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clubs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    aquisition_channel VARCHAR(100),
    aquisition_location VARCHAR(100),
    plan_type VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE IF NOT EXISTS user_club (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    club_id INTEGER NOT NULL REFERENCES clubs(id),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_clubs_name ON clubs(name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_club_userid_clubid ON user_club(user_id, club_id);
CREATE INDEX IF NOT EXISTS idx_user_club_user_id ON user_club(user_id);
CREATE INDEX IF NOT EXISTS idx_user_club_club_id ON user_club(club_id);

ALTER TABLE clubs ADD CONSTRAINT unique_club_name UNIQUE (name);
ALTER TABLE users ADD CONSTRAINT unique_user_email UNIQUE (email);
ALTER TABLE user_club ADD CONSTRAINT unique_user_club UNIQUE (user_id, club_id);



CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_updated_at_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trigger_set_updated_at_clubs
BEFORE UPDATE ON clubs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trigger_set_updated_at_user_club
BEFORE UPDATE ON user_club
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

ANALYZE;
