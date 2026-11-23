package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/requests"
	"github.com/lailiseptiandi/be-kominfo-test/internal/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, data models.Room) (*models.Room, error)
	FindByID(ctx context.Context, id int) (*models.Room, error)
	GetAvailable(ctx context.Context, req requests.RoomRequestParam) ([]*models.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *roomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room models.Room) (*models.Room, error) {
	err := r.db.WithContext(ctx).FirstOrCreate(&room).Error
	if err != nil {
		return nil, err
	}

	roomDetail, err := r.FindByID(ctx, int(room.ID))
	if err != nil {
		return nil, err
	}
	return roomDetail, nil
}

func (r *roomRepository) FindByID(ctx context.Context, id int) (*models.Room, error) {
	var room models.Room
	err := r.db.WithContext(ctx).Model(&models.Room{}).Where("id = ?", id).Find(&room).Error
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *roomRepository) GetAvailable(ctx context.Context, req requests.RoomRequestParam) ([]*models.Room, error) {
	var data []*models.Room
	query := r.db.WithContext(ctx).Model(&models.Room{})

	if req.Start != "" && req.End != "" {
		startTime, err := time.Parse(time.RFC3339, req.Start)
		if err != nil {
			return nil, errors.New("format start harus ISO-8601 (RFC3339)")
		}
		endTime, err := time.Parse(time.RFC3339, req.End)
		if err != nil {
			return nil, errors.New("format end harus ISO-8601 (RFC3339)")
		}

		subQuery := r.db.Table("bookings").
			Select("1").
			Where("bookings.room_id = rooms.id").
			Where("bookings.start_time < ? AND bookings.end_time > ?", endTime, startTime)

		query = query.Where("NOT EXISTS (?)", subQuery)
	}

	err := query.Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil

}
