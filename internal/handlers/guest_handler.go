package handlers

import (
	"net/http"

	"hotel-api/internal/models"
	"hotel-api/internal/services"
	"hotel-api/pkg/errors"

	"github.com/gin-gonic/gin"
)

type GuestHandler struct {
	service services.GuestService
}

func NewGuestHandler(service services.GuestService) *GuestHandler {
	return &GuestHandler{service: service}
}

// Create godoc
// @Summary Cria um hóspede
// @Tags guests
// @Accept json
// @Produce json
// @Param request body models.CreateGuestRequest true "Dados do Hóspede"
// @Success 201 {object} models.Guest
// @Failure 400 {object} middlewares.ErrorResponse
// @Failure 409 {object} middlewares.ErrorResponse
// @Router /guests [post]
func (h *GuestHandler) Create(c *gin.Context) {
	var req models.CreateGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrValidation)
		return
	}

	guest, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, guest)
}

// GetAll godoc
// @Summary Lista todos os hóspedes
// @Tags guests
// @Produce json
// @Success 200 {array} models.Guest
// @Router /guests [get]
func (h *GuestHandler) GetAll(c *gin.Context) {
	guests, err := h.service.GetAll()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, guests)
}

// GetByID godoc
// @Summary Busca um hóspede por ID
// @Tags guests
// @Produce json
// @Param id path string true "ID do Hóspede"
// @Success 200 {object} models.Guest
// @Failure 404 {object} middlewares.ErrorResponse
// @Router /guests/{id} [get]
func (h *GuestHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	guest, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, guest)
}

// Update godoc
// @Summary Atualiza um hóspede
// @Tags guests
// @Accept json
// @Produce json
// @Param id path string true "ID do Hóspede"
// @Param request body models.UpdateGuestRequest true "Dados para atualizar"
// @Success 200 {object} models.Guest
// @Failure 400 {object} middlewares.ErrorResponse
// @Failure 404 {object} middlewares.ErrorResponse
// @Router /guests/{id} [put]
func (h *GuestHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrValidation)
		return
	}

	guest, err := h.service.Update(id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, guest)
}

// Delete godoc
// @Summary Remove um hóspede
// @Tags guests
// @Param id path string true "ID do Hóspede"
// @Success 204
// @Failure 404 {object} middlewares.ErrorResponse
// @Router /guests/{id} [delete]
func (h *GuestHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
