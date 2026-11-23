package services

import (
	"context"
	"errors"
	"time"

	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/requests"
	"github.com/lailiseptiandi/be-kominfo-test/internal/models"
	"github.com/lailiseptiandi/be-kominfo-test/internal/repository"
)

type BookingService interface {
	Create(ctx context.Context, userID uint, req requests.BookRequest) (*models.Booking, error)
	Cancel(ctx context.Context, id int) (*models.Booking, error)
}

type bookingService struct {
	bookingRepo repository.BookingRepository
}

func NewBookingService(bookingRepo repository.BookingRepository) *bookingService {
	return &bookingService{bookingRepo: bookingRepo}
}

func (s *bookingService) Create(ctx context.Context, userID uint, req requests.BookRequest) (*models.Booking, error) {
	if req.RoomID <= 0 || req.Title == "" {
		return nil, errors.New("room_id dan title wajib diisi")
	}
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}
	if !endTime.After(startTime) {
		return nil, err
	}

	booking := models.Booking{
		RoomID:    uint(req.RoomID),
		UserID:    userID,
		Title:     req.Title,
		StartTime: startTime.UTC(),
		EndTime:   endTime.UTC(),
		Status:    "CONFIRMED",
	}
	return s.bookingRepo.Create(ctx, booking)
}

func (s *bookingService) Cancel(ctx context.Context, id int) (*models.Booking, error) {
	if id <= 0 {
		return nil, errors.New("id booking tidak valid")
	}
	return s.bookingRepo.Cancel(ctx, id)
}
