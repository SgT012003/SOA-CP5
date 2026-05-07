package repositories

import (
	"database/sql"
	"hotel-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type GuestRepository interface {
	GetByID(id string) (*models.Guest, error)
	GetAll() ([]models.Guest, error)
	Create(guest *models.Guest) error
	Update(guest *models.Guest) error
	Delete(id string) error
	GetByEmail(email string) (*models.Guest, error)
	GetByDocument(document string) (*models.Guest, error)
}

type guestRepository struct {
	db *sqlx.DB
}

func NewGuestRepository(db *sqlx.DB) GuestRepository {
	return &guestRepository{db: db}
}

func (r *guestRepository) GetByID(id string) (*models.Guest, error) {
	var guest models.Guest
	err := r.db.Get(&guest, "SELECT * FROM guests WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &guest, err
}

func (r *guestRepository) GetAll() ([]models.Guest, error) {
	var guests []models.Guest
	err := r.db.Select(&guests, "SELECT * FROM guests")
	return guests, err
}

func (r *guestRepository) Create(guest *models.Guest) error {
	query := `INSERT INTO guests (id, full_name, document, email, phone, created_at) VALUES (:id, :full_name, :document, :email, :phone, :created_at)`
	_, err := r.db.NamedExec(query, guest)
	return err
}

func (r *guestRepository) Update(guest *models.Guest) error {
	query := `UPDATE guests SET full_name = :full_name, phone = :phone WHERE id = :id`
	_, err := r.db.NamedExec(query, guest)
	return err
}

func (r *guestRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM guests WHERE id = $1", id)
	return err
}

func (r *guestRepository) GetByEmail(email string) (*models.Guest, error) {
	var guest models.Guest
	err := r.db.Get(&guest, "SELECT * FROM guests WHERE email = $1", email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &guest, err
}

func (r *guestRepository) GetByDocument(document string) (*models.Guest, error) {
	var guest models.Guest
	err := r.db.Get(&guest, "SELECT * FROM guests WHERE document = $1", document)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &guest, err
}
