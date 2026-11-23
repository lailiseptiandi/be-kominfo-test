package repository

import (
	"context"
	"errors"

	"github.com/lailiseptiandi/be-kominfo-test/internal/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(ctx context.Context, req models.Booking) (*models.Booking, error)
	FindByID(ctx context.Context, id int) (*models.Booking, error)
	Cancel(ctx context.Context, id int) (*models.Booking, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, req models.Booking) (*models.Booking, error) {
	var room models.Room
	if err := r.db.WithContext(ctx).First(&room, req.RoomID).Error; err != nil {
		return nil, err
	}

	var count int64
	err := r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("room_id = ? AND start_time < ? AND end_time > ?", req.RoomID, req.EndTime, req.StartTime).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("room is already booked for this time")
	}

	err = r.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, int(req.ID))
}

func (r *bookRepository) FindByID(ctx context.Context, id int) (*models.Booking, error) {
	var data models.Booking
	err := r.db.WithContext(ctx).Preload("Room").Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *bookRepository) Cancel(ctx context.Context, id int) (*models.Booking, error) {
	result := r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("id = ?", id).
		Update("status", "CANCELLED")
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.FindByID(ctx, id)
}
