CREATE TYPE user_roles as ENUM ('user','admin','subadmin');
CREATE TABLE roles
(
    user_id     UUID NOT NULL REFERENCES users (id),
    role        user_roles               DEFAULT 'user'::user_roles,
    username    text not null,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

create unique index on roles (username) where archived_at is null