-- Rekap Jam Pemakaian Ruangan

select
    r.id as room_id,
    r.name,
    round(
        coalesce(
            sum(extract(EPOCH FROM (b.end_time - b.start_time)) / 3600.0),
            0
        )::numeric,
        2
    ) AS total_hours
from rooms r
left join bookings b on b.room_id = r.id
    and status = 'CONFIRMED'
    and b.start_time >=date_trunc('month', now()
    and b.end_time < date_trunc('month', now()
    and b.start_time < date_trunc('month', now() + interval '1 month')

group by r.id, r.name
order by total_hours, r.name asc

-- Pencarian Ruangan Tersedia
select r.id, r.name, r.capacity
from rooms r
where r.is_active = true
and not exists (
    select 1
    from bookings b
    where b.room_id = r.id
    and b.status = 'CONFIRMED'
    and b.start_time < '2026-12-01T12:00:00Z'
    and b.end_time > '2026-12-01T10:00:00Z'
)

order by r.name asc
