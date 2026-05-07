package services

import (
	"net/http"

	"hotel-api/internal/models"
	"hotel-api/internal/repositories"
	"hotel-api/pkg/errors"

	"github.com/google/uuid"
)

type RoomService interface {
	Create(req *models.CreateRoomRequest) (*models.Room, error)
	GetByID(id string) (*models.Room, error)
	GetAll() ([]models.Room, error)
	Update(id string, req *models.UpdateRoomRequest) (*models.Room, error)
	Delete(id string) error
}

type roomService struct {
	roomRepo repositories.RoomRepository
}

func NewRoomService(roomRepo repositories.RoomRepository) RoomService {
	return &roomService{roomRepo: roomRepo}
}

func (s *roomService) Create(req *models.CreateRoomRequest) (*models.Room, error) {
	existing, err := s.roomRepo.GetByNumber(req.Number)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.NewAppError(http.StatusConflict, "RoomNumberAlreadyExists", "Já existe um quarto com este número.")
	}

	room := &models.Room{
		ID:            uuid.New().String(),
		Number:        req.Number,
		Type:          req.Type,
		Capacity:      req.Capacity,
		PricePerNight: req.PricePerNight,
		Status:        "ATIVO",
	}

	err = s.roomRepo.Create(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) GetByID(id string) (*models.Room, error) {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if room == nil || room.Status != "ATIVO" {
		return nil, errors.ErrNotFound
	}
	return room, nil
}

func (s *roomService) GetAll() ([]models.Room, error) {
	return s.roomRepo.GetAll()
}

func (s *roomService) Update(id string, req *models.UpdateRoomRequest) (*models.Room, error) {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if room == nil || room.Status != "ATIVO" {
		return nil, errors.ErrNotFound
	}

	room.Type = req.Type
	room.Capacity = req.Capacity
	room.PricePerNight = req.PricePerNight

	err = s.roomRepo.Update(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) Delete(id string) error {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		return err
	}
	if room == nil || room.Status != "ATIVO" {
		return errors.ErrNotFound
	}

	return s.roomRepo.UpdateStatus(id, "INATIVO")
}
