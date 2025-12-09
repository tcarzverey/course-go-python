-- +goose Up
create table rooms
(
    id         int generated always as identity primary key,
    name       text        not null,
    capacity   int         not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table bookings
(
    id         int generated always as identity primary key,
    room_id    int references rooms (id) not null,
    user_id    text                      not null,
    capacity   int                       not null,
    status     int                       not null,
    start_time timestamptz               not null,
    end_time   timestamptz               not null,
    created_at timestamptz               not null default now(),
    updated_at timestamptz               not null default now()
);
