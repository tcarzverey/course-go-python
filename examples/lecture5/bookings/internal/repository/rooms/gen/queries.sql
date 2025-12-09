-- name: CreateRoom :one
insert into rooms (name, capacity)
values (@name, @capacity)
returning id, name, capacity, created_at, updated_at;

-- name: ListRooms :many
select
    id,
    name,
    capacity,
    created_at,
    updated_at
from rooms
order by id;

-- name: GetRoom :one
select
    id,
    name,
    capacity,
    created_at,
    updated_at
from rooms
where id = @id;