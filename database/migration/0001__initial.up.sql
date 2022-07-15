CREATE TABLE IF NOT EXISTS users(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email text not null unique,
  username text not null unique,
  password text not null,
  created_by uuid references users(id) default 00000000-0000-0000-0000-000000000000,
  created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  archived_at   TIMESTAMP WITH TIME ZONE default null
);