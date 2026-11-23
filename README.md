# Booking API

## Jalankan dengan Docker

```sh
docker compose up --build
```

API berjalan di `http://localhost:9000`, PostgreSQL di port `5433`.

## Endpoint

```sh
# Health check
curl http://localhost:9000/health

# Login
curl -X POST http://localhost:9000/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user1@test.com","password":"password"}'

# Gunakan token hasil login untuk endpoint berikutnya
TOKEN='hasil-login'

# Buat room
curl -X POST http://localhost:9000/api/v1/rooms \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Ruang Meeting","capacity":10}'

# Cek room tersedia
curl 'http://localhost:9000/api/v1/rooms/available?start=2026-12-01T10:00:00Z&end=2026-12-01T12:00:00Z' \
  -H "Authorization: Bearer $TOKEN"

# Buat booking
curl -X POST http://localhost:9000/api/v1/bookings \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"room_id":1,"title":"Sprint Planning","start_time":"2026-12-01T10:00:00Z","end_time":"2026-12-01T12:00:00Z"}'

# Batalkan booking
curl -X PATCH http://localhost:9000/api/v1/bookings/1/cancel \
  -H "Authorization: Bearer $TOKEN"
```
