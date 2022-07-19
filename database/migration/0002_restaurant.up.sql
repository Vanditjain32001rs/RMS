create table if not exists restaurant
(
    id          uuid primary key         default gen_random_uuid(),
    name        text             not null,
    created_by  uuid             not null references users (id),
    latitude    double precision not null,
    longitude   double precision not null,
    created_at  timestamp with time zone default now(),
    archived_at timestamp with time zone default null
);