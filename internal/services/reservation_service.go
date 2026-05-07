package services

import (
	"math"
	"net/http"
	"time"

	"hotel-api/internal/models"
	"hotel-api/internal/repositories"
	"hotel-api/pkg/errors"

	"github.com/google/uuid"
)

type ReservationService interface {
	CreateReservation(req *models.CreateReservationRequest) (*models.Reservation, error)
	UpdateStatus(id string, req *models.UpdateReservationStatusRequest) error
}

type reservationService struct {
	reservationRepo repositories.ReservationRepository
	roomRepo        repositories.RoomRepository
	guestRepo       repositories.GuestRepository
}

func NewReservationService(
	reservationRepo repositories.ReservationRepository,
	roomRepo repositories.RoomRepository,
	guestRepo repositories.GuestRepository,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		roomRepo:        roomRepo,
		guestRepo:       guestRepo,
	}
}

func (s *reservationService) CreateReservation(req *models.CreateReservationRequest) (*models.Reservation, error) {
	// A data de check-out deve ser estritamente maior que a data de check-in.
	if !req.CheckoutExpected.After(req.CheckinExpected) {
		return nil, errors.ErrInvalidDateRange
	}

	guest, err := s.guestRepo.GetByID(req.GuestID)
	if err != nil {
		return nil, err
	}
	if guest == nil {
		return nil, errors.ErrNotFound
	}

	room, err := s.roomRepo.GetByID(req.RoomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.ErrNotFound
	}

	// O quarto deve estar ativo
	if room.Status != "ATIVO" {
		return nil, errors.ErrRoomUnavailable
	}

	// O número de hóspedes não pode exceder a capacity do quarto.
	if req.GuestsCount > room.Capacity {
		return nil, errors.ErrCapacityExceeded
	}

	// Validar sobreposição de datas.
	overlapping, err := s.reservationRepo.CheckOverlapping(req.RoomID, req.CheckinExpected, req.CheckoutExpected)
	if err != nil {
		return nil, err
	}
	if overlapping {
		return nil, errors.ErrRoomUnavailable
	}

	// Calcular valor estimado
	duration := req.CheckoutExpected.Sub(req.CheckinExpected)
	days := math.Ceil(duration.Hours() / 24)
	if days < 1 {
		days = 1
	}
	estimatedAmount := days * room.PricePerNight

	reservation := &models.Reservation{
		ID:               uuid.New().String(),
		GuestID:          guest.ID,
		RoomID:           room.ID,
		CheckinExpected:  req.CheckinExpected,
		CheckoutExpected: req.CheckoutExpected,
		Status:           "CREATED",
		EstimatedAmount:  estimatedAmount,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.reservationRepo.Create(reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *reservationService) UpdateStatus(id string, req *models.UpdateReservationStatusRequest) error {
	reservation, err := s.reservationRepo.GetByID(id)
	if err != nil {
		return err
	}
	if reservation == nil {
		return errors.ErrNotFound
	}

	now := time.Now()

	switch req.Status {
	case "CHECKED_IN":
		if reservation.Status != "CREATED" {
			return errors.ErrInvalidState
		}
		// Janela de check-in: até 3 horas de antecedência
		minCheckinTime := reservation.CheckinExpected.Add(-3 * time.Hour)
		if now.Before(minCheckinTime) {
			return errors.NewAppError(http.StatusBadRequest, "CheckInTooEarly", "O check-in só pode ser efetuado com até 3 horas de antecedência da data/hora prevista.")
		}
		reservation.Status = "CHECKED_IN"
		reservation.CheckinAt = &now
		reservation.UpdatedAt = now

	case "CHECKED_OUT":
		if reservation.Status != "CHECKED_IN" {
			return errors.ErrInvalidState
		}
		reservation.Status = "CHECKED_OUT"
		reservation.CheckoutAt = &now
		reservation.UpdatedAt = now

		room, err := s.roomRepo.GetByID(reservation.RoomID)
		if err != nil {
			return err
		}
		duration := reservation.CheckoutAt.Sub(*reservation.CheckinAt)
		days := math.Ceil(duration.Hours() / 24)
		if days < 1 {
			days = 1
		}
		finalAmount := days * room.PricePerNight
		reservation.FinalAmount = &finalAmount

	case "CANCELED":
		if reservation.Status != "CREATED" {
			return errors.ErrInvalidState
		}
		reservation.Status = "CANCELED"
		reservation.UpdatedAt = now
	}

	return s.reservationRepo.Update(reservation)
}
