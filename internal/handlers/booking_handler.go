package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/requests"
	"github.com/lailiseptiandi/be-kominfo-test/internal/services"
	"github.com/lailiseptiandi/be-kominfo-test/pkg/utils/jwt"
	"github.com/lailiseptiandi/be-kominfo-test/pkg/utils/response"
	"gorm.io/gorm"
)

type BookingHandler struct {
	bookingService services.BookingService
}

func NewBookingHandler(booking services.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: booking}
}

func (h *BookingHandler) Create(ctx *gin.Context) {
	var req requests.BookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error(), "")
		return
	}

	claims, ok := ctx.Get("credential_user")
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized: Invalid token", "")
		return
	}
	user, ok := claims.(*jwt.Claims)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized: Invalid token", "")
		return
	}

	booking, err := h.bookingService.Create(ctx, uint(user.UserID), req)
	if err != nil {
		response.BadRequest(ctx, err.Error(), "")
		return
	}

	response.Created(ctx, booking)
}

func (h *BookingHandler) Cancel(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.BadRequest(ctx, "id booking tidak valid", "")
		return
	}

	booking, err := h.bookingService.Cancel(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(ctx, "booking tidak ditemukan", "")
			return
		}
		response.InternalServerError(ctx, err.Error(), "")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": booking.ID, "status": booking.Status})
}
