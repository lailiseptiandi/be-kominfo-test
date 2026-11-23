package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/requests"
	"github.com/lailiseptiandi/be-kominfo-test/internal/services"
	"github.com/lailiseptiandi/be-kominfo-test/pkg/utils/response"
)

type RoomHandler struct {
	roomService services.RoomService
}

func NewRoomHandler(room services.RoomService) *RoomHandler {
	return &RoomHandler{roomService: room}
}

func (h *RoomHandler) Create(ctx *gin.Context) {
	var req requests.RoomRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error(), "")
		return
	}

	room, err := h.roomService.Create(ctx, req)
	if err != nil {
		response.NotFound(ctx, err.Error(), "")
		return
	}

	response.Created(ctx, room)

}
