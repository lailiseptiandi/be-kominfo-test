package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lailiseptiandi/be-kominfo-test/internal/models"
	"github.com/lailiseptiandi/be-kominfo-test/internal/routers"
	appjwt "github.com/lailiseptiandi/be-kominfo-test/pkg/utils/jwt"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateBookingRejectsConflictingConcurrentRequests(t *testing.T) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		t.Skip("atur DB_HOST dan variabel DB_* lainnya untuk menjalankan tes integrasi")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, env("DB_USERNAME", "app"), env("DB_PASSWORD", "app"),
		env("DB_NAME", "db_kominfo_test"), env("DB_PORT", "5433"))
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("gagal terhubung ke database tes: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("gagal mengambil koneksi database: %v", err)
	}
	defer sqlDB.Close()

	var user models.User
	if err := db.First(&user).Error; err != nil {
		t.Fatalf("gagal mengambil pengguna tes: %v", err)
	}

	room := models.Room{Name: fmt.Sprintf("booking-test-%d", time.Now().UnixNano()), Capacity: 10}
	if err := db.Create(&room).Error; err != nil {
		t.Fatalf("gagal membuat ruangan tes: %v", err)
	}
	t.Cleanup(func() { db.Delete(&room) })

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		t.Fatal("JWT_SECRET harus diatur sebelum menjalankan tes integrasi")
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, appjwt.Claims{
		UserID: uint64(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("gagal membuat token tes: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	routers.InitRoutes(db, r)
	start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Minute)
	title := fmt.Sprintf("Tes integrasi %d", time.Now().UnixNano())
	payload, _ := json.Marshal(map[string]any{
		"room_id": room.ID, "title": title,
		"start_time": start.Format(time.RFC3339), "end_time": start.Add(time.Hour).Format(time.RFC3339),
	})

	statuses := make(chan int, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)
			statuses <- res.Code
		}()
	}
	wg.Wait()
	close(statuses)

	created, rejected := 0, 0
	for status := range statuses {
		if status == http.StatusCreated {
			created++
		} else if status == http.StatusBadRequest {
			rejected++
		} else {
			t.Errorf("menerima status respons yang tidak sesuai: %d", status)
		}
	}
	if created != 1 || rejected != 1 {
		t.Fatalf("seharusnya satu booking berhasil dan satu ditolak, hasil: berhasil=%d ditolak=%d", created, rejected)
	}

	var savedBookings int64
	if err := db.Model(&models.Booking{}).
		Where("room_id = ? AND title = ?", room.ID, title).
		Count(&savedBookings).Error; err != nil {
		t.Fatalf("gagal memeriksa booking yang tersimpan: %v", err)
	}
	if savedBookings != 1 {
		t.Fatalf("seharusnya ada satu booking tersimpan di database, ditemukan %d", savedBookings)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
