package handlers

import (
	"net/http"

	"hotel-api/internal/models"
	"hotel-api/internal/services"
	"hotel-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	service services.RoomService
}

func NewRoomHandler(service services.RoomService) *RoomHandler {
	return &RoomHandler{service: service}
}

// Create godoc
// @Summary Cria um quarto
// @Tags rooms
// @Accept json
// @Produce json
// @Param request body models.CreateRoomRequest true "Dados do Quarto"
// @Success 201 {object} models.Room
// @Failure 400 {object} middlewares.ErrorResponse
// @Failure 409 {object} middlewares.ErrorResponse
// @Router /rooms [post]
func (h *RoomHandler) Create(c *gin.Context) {
	var req models.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrValidation)
		return
	}

	room, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, room)
}

// GetAll godoc
// @Summary Lista todos os quartos ativos
// @Tags rooms
// @Produce json
// @Success 200 {array} models.Room
// @Router /rooms [get]
func (h *RoomHandler) GetAll(c *gin.Context) {
	rooms, err := h.service.GetAll()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// GetByID godoc
// @Summary Busca um quarto por ID
// @Tags rooms
// @Produce json
// @Param id path string true "ID do Quarto"
// @Success 200 {object} models.Room
// @Failure 404 {object} middlewares.ErrorResponse
// @Router /rooms/{id} [get]
func (h *RoomHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	room, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, room)
}

// Update godoc
// @Summary Atualiza um quarto
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "ID do Quarto"
// @Param request body models.UpdateRoomRequest true "Dados para atualizar"
// @Success 200 {object} models.Room
// @Failure 400 {object} middlewares.ErrorResponse
// @Failure 404 {object} middlewares.ErrorResponse
// @Router /rooms/{id} [put]
func (h *RoomHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrValidation)
		return
	}

	room, err := h.service.Update(id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, room)
}

// Delete godoc
// @Summary Inativa um quarto (soft delete)
// @Tags rooms
// @Param id path string true "ID do Quarto"
// @Success 204
// @Failure 404 {object} middlewares.ErrorResponse
// @Router /rooms/{id} [delete]
func (h *RoomHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
