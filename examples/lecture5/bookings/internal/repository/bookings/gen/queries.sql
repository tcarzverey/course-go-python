-- name: CreateBooking :one
insert into bookings (room_id, user_id, start_time, end_time)
values (@room_id, @user_id, @start_time, @end_time)
returning
    id,
    room_id,
    user_id,
    start_time,
    end_time,
    status,
    created_at,
    updated_at;

-- name: ListBookings :many
select
    id,
    room_id,
    user_id,
    start_time,
    end_time,
    status,
    created_at,
    updated_at
from bookings
order by id;

-- name: GetBooking :one
select
    id,
    room_id,
    user_id,
    start_time,
    end_time,
    status,
    created_at,
    updated_at
from bookings
where id = @id;

-- name: CancelBooking :one
update bookings
set
    status = 'cancelled',
    updated_at = now()
where id = @id
returning
    id,
    room_id,
    user_id,
    start_time,
    end_time,
    status,
    created_at,
    updated_at;

-- name: ListBookingsByPeriod :many
select
    id,
    room_id,
    user_id,
    start_time,
    end_time,
    status,
    created_at,
    updated_at
from bookings
where
    start_time < @start
  and end_time > @finish
order by room_id, start_time;