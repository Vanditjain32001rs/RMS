create table if not exists dishes(
    id uuid primary key default gen_random_uuid(),
    name text not null,
    price double precision not null,
    restaurant_id uuid not null references restaurant(id),
    created_by uuid references users(id),
    created_at timestamp with time zone default now(),
    archived_at timestamp with time zone default null
);

ALTER TABLE dishes ALTER COLUMN price TYPE float;