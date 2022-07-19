CREATE TABLE IF NOT EXISTS users(

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email text not null,
  username text not null,
  password text not null,
  created_by uuid references users(id),
  created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  archived_at   TIMESTAMP WITH TIME ZONE default null
);