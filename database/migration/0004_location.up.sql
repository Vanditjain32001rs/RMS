CREATE TABLE location
(
    user_id     UUID  NOT NULL REFERENCES users (id),
    latitude    float NOT NULL,
    longitude   float NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);