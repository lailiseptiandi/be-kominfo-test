-- Rekap Jam Pemakaian Ruangan
SELECT
    r.id AS room_id,
    r.name,
    ROUND(
        COALESCE(
            SUM(
                EXTRACT(EPOCH FROM (b.end_time - b.start_time)) / 3600.0
            ),
            0
        )::numeric,
        2
    ) AS total_hours
FROM rooms r
LEFT JOIN bookings b
    ON b.room_id = r.id
    AND b.status = 'CONFIRMED'
    AND b.start_time >= date_trunc('month', NOW())
    AND b.start_time < date_trunc('month', NOW() + INTERVAL '1 month')
GROUP BY r.id, r.name
ORDER BY total_hours ASC, r.name ASC;

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
