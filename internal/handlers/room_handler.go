package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/requests"
	"github.com/lailiseptiandi/be-kominfo-test/internal/repository"
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
		if errors.Is(err, repository.ErrRoomNameExists) {
			response.Error(ctx, http.StatusConflict, err.Error(), "")
			return
		}
		response.NotFound(ctx, err.Error(), "")
		return
	}

	response.Created(ctx, room)

}

func (h *RoomHandler) GetAvailable(ctx *gin.Context) {
	req := requests.RoomRequestParam{
		Start: ctx.Query("start"),
		End:   ctx.Query("end"),
	}

	data, err := h.roomService.GetAvailable(ctx, req)
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error(), "")
		return
	}

	response.Success(ctx, data)

}
