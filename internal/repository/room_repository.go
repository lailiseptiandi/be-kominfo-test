package repository

import (
	"context"

	"github.com/lailiseptiandi/be-kominfo-test/internal/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, data models.Room) (*models.Room, error)
	FindByID(ctx context.Context, id int) (*models.Room, error)
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
