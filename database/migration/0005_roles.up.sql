CREATE TYPE user_roles as ENUM ('user','admin','subadmin');
CREATE TABLE roles(
    user_id UUID NOT NULL REFERENCES users(id),
    role user_roles DEFAULT 'user',
    username text not null unique,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

ALTER TABLE roles drop constraint roles_username_key;
