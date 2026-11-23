package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/requests"
	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/responses"
	"github.com/lailiseptiandi/be-kominfo-test/internal/services"
	"github.com/lailiseptiandi/be-kominfo-test/pkg/utils/response"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(auth services.AuthService) *AuthHandler {
	return &AuthHandler{authService: auth}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req requests.UserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error(), "")
		return
	}

	userToken, err := h.authService.Login(ctx, req)
	if err != nil {
		response.NotFound(ctx, err.Error(), "")
		return
	}

	userFormatter := (&responses.UserResponse{}).FormatterLoginUser(*userToken)
	response.Created(ctx, userFormatter)

}
