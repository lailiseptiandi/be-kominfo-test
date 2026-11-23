# Booking API

## Jalankan dengan Docker

```sh
docker compose up --build
```

API berjalan di `http://localhost:9000`, PostgreSQL di port `5433`.

## Booking integration test

Jalankan PostgreSQL dan pastikan database memiliki skema dari `pkg/database/migration/init.sql` serta minimal satu user. `docker compose up` membuat skema tersebut pada inisialisasi database:

```sh
docker compose up -d postgres
```

Jalankan integration test dari root project:

```sh
DB_HOST=localhost DB_PORT=5433 DB_NAME=db_kominfo_test DB_USERNAME=app DB_PASSWORD=app JWT_SECRET=local-dev-secret go test ./internal/handlers -run '^TestCreateBookingRejectsConflictingConcurrentRequests$' -v
```

Test mengirim dua request booking secara bersamaan untuk room dan jadwal yang sama, lalu memastikan tepat satu booking tersimpan di database. Test akan dilewati jika `DB_HOST` tidak diatur.

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
