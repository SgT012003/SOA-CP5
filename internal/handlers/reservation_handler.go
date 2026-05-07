package handlers

import (
	"net/http"

	"hotel-api/internal/models"
	"hotel-api/internal/services"
	"hotel-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	service services.ReservationService
}

func NewReservationHandler(service services.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

// CreateReservation godoc
// @Summary Cria uma reserva
// @Description Cria uma nova reserva validando disponibilidade e regras de negócio
// @Tags reservations
// @Accept json
// @Produce json
// @Param request body models.CreateReservationRequest true "Dados da Reserva"
// @Success 201 {object} models.Reservation
// @Failure 400 {object} middlewares.ErrorResponse
// @Failure 409 {object} middlewares.ErrorResponse
// @Router /reservations [post]
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req models.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrValidation)
		return
	}

	reservation, err := h.service.CreateReservation(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// UpdateStatus godoc
// @Summary Atualiza o status da reserva
// @Description Realiza transições de status da reserva (ex: CHECKED_IN, CHECKED_OUT, CANCELED)
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "ID da Reserva"
// @Param request body models.UpdateReservationStatusRequest true "Novo Status"
// @Success 200
// @Failure 400 {object} middlewares.ErrorResponse
// @Failure 404 {object} middlewares.ErrorResponse
// @Failure 409 {object} middlewares.ErrorResponse
// @Router /reservations/{id}/status [patch]
func (h *ReservationHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateReservationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrValidation)
		return
	}

	if err := h.service.UpdateStatus(id, &req); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}
