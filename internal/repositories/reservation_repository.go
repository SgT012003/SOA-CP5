package repositories

import (
	"hotel-api/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReservationRepository interface {
	Create(reservation *models.Reservation) error
	CheckOverlapping(roomID string, checkin, checkout time.Time) (bool, error)
	GetByID(id string) (*models.Reservation, error)
	Update(reservation *models.Reservation) error
}

type reservationRepository struct {
	db *sqlx.DB
}

func NewReservationRepository(db *sqlx.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(reservation *models.Reservation) error {
	query := `
		INSERT INTO reservations (id, guest_id, room_id, checkin_expected, checkout_expected, status, estimated_amount, created_at, updated_at)
		VALUES (:id, :guest_id, :room_id, :checkin_expected, :checkout_expected, :status, :estimated_amount, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, reservation)
	return err
}

func (r *reservationRepository) CheckOverlapping(roomID string, checkin, checkout time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) FROM reservations 
		WHERE room_id = $1 
		AND status != 'CANCELED'
		AND (
			(checkin_expected < $3 AND checkout_expected > $2)
		)
	`
	var count int
	err := r.db.Get(&count, query, roomID, checkin, checkout)
	return count > 0, err
}

func (r *reservationRepository) GetByID(id string) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Get(&reservation, "SELECT * FROM reservations WHERE id = $1", id)
	return &reservation, err
}

func (r *reservationRepository) Update(reservation *models.Reservation) error {
	query := `
		UPDATE reservations 
		SET checkin_at = :checkin_at, checkout_at = :checkout_at, status = :status, final_amount = :final_amount, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExec(query, reservation)
	return err
}
