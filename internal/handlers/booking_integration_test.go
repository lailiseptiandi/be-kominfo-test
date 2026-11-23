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
		t.Skip("set DB_HOST and related DB_* variables to run integration test")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		dbHost, env("DB_USERNAME", "app"), env("DB_PASSWORD", "app"),
		env("DB_NAME", "db_kominfo_test"), env("DB_PORT", "5433"))
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	var user models.User
	if err := db.First(&user).Error; err != nil {
		t.Fatalf("load test user: %v", err)
	}

	room := models.Room{Name: fmt.Sprintf("booking-test-%d", time.Now().UnixNano()), Capacity: 10}
	if err := db.Create(&room).Error; err != nil {
		t.Fatalf("create test room: %v", err)
	}
	t.Cleanup(func() { db.Delete(&room) })

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		t.Fatal("JWT_SECRET must be set before running integration test")
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, appjwt.Claims{
		UserID: uint64(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	routers.InitRoutes(db, r)
	start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Minute)
	title := fmt.Sprintf("Integration test %d", time.Now().UnixNano())
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
			t.Errorf("unexpected response status: %d", status)
		}
	}
	if created != 1 || rejected != 1 {
		t.Fatalf("expected one booking and one conflict, got created=%d rejected=%d", created, rejected)
	}

	var savedBookings int64
	if err := db.Model(&models.Booking{}).
		Where("room_id = ? AND title = ?", room.ID, title).
		Count(&savedBookings).Error; err != nil {
		t.Fatalf("check saved booking: %v", err)
	}
	if savedBookings != 1 {
		t.Fatalf("expected one booking saved in database, got %d", savedBookings)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
