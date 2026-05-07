package services

import (
	"net/http"
	"time"

	"hotel-api/internal/models"
	"hotel-api/internal/repositories"
	"hotel-api/pkg/errors"

	"github.com/google/uuid"
)

type GuestService interface {
	Create(req *models.CreateGuestRequest) (*models.Guest, error)
	GetByID(id string) (*models.Guest, error)
	GetAll() ([]models.Guest, error)
	Update(id string, req *models.UpdateGuestRequest) (*models.Guest, error)
	Delete(id string) error
}

type guestService struct {
	guestRepo repositories.GuestRepository
}

func NewGuestService(guestRepo repositories.GuestRepository) GuestService {
	return &guestService{guestRepo: guestRepo}
}

func (s *guestService) Create(req *models.CreateGuestRequest) (*models.Guest, error) {
	existing, err := s.guestRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.NewAppError(http.StatusConflict, "EmailAlreadyInUse", "Este email já está em uso.")
	}

	existing, err = s.guestRepo.GetByDocument(req.Document)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.NewAppError(http.StatusConflict, "DocumentAlreadyInUse", "Este documento já está em uso.")
	}

	guest := &models.Guest{
		ID:        uuid.New().String(),
		FullName:  req.FullName,
		Document:  req.Document,
		Email:     req.Email,
		Phone:     req.Phone,
		CreatedAt: time.Now(),
	}

	err = s.guestRepo.Create(guest)
	if err != nil {
		return nil, err
	}

	return guest, nil
}

func (s *guestService) GetByID(id string) (*models.Guest, error) {
	guest, err := s.guestRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if guest == nil {
		return nil, errors.ErrNotFound
	}
	return guest, nil
}

func (s *guestService) GetAll() ([]models.Guest, error) {
	return s.guestRepo.GetAll()
}

func (s *guestService) Update(id string, req *models.UpdateGuestRequest) (*models.Guest, error) {
	guest, err := s.guestRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if guest == nil {
		return nil, errors.ErrNotFound
	}

	guest.FullName = req.FullName
	guest.Phone = req.Phone

	err = s.guestRepo.Update(guest)
	if err != nil {
		return nil, err
	}

	return guest, nil
}

func (s *guestService) Delete(id string) error {
	guest, err := s.guestRepo.GetByID(id)
	if err != nil {
		return err
	}
	if guest == nil {
		return errors.ErrNotFound
	}
	return s.guestRepo.Delete(id)
}
