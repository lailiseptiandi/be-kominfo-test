package services

import (
	"context"

	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/requests"
	"github.com/lailiseptiandi/be-kominfo-test/internal/models"
	"github.com/lailiseptiandi/be-kominfo-test/internal/repository"
)

type RoomService interface {
	Create(ctx context.Context, req requests.RoomRequest) (*models.Room, error)
	GetAvailable(ctx context.Context, req requests.RoomRequestParam) ([]*models.Room, error)
}

type roomService struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) *roomService {
	return &roomService{roomRepo: roomRepo}
}

func (r *roomService) Create(ctx context.Context, req requests.RoomRequest) (*models.Room, error) {
	data := models.Room{
		Name:     req.Name,
		Capacity: req.Capacity,
	}

	room, err := r.roomRepo.Create(ctx, data)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *roomService) GetAvailable(ctx context.Context, req requests.RoomRequestParam) ([]*models.Room, error) {
	data, err := r.roomRepo.GetAvailable(ctx, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}
