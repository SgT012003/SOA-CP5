package repositories

import (
	"database/sql"
	"hotel-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type RoomRepository interface {
	GetByID(id string) (*models.Room, error)
	GetAll() ([]models.Room, error)
	Create(room *models.Room) error
	Update(room *models.Room) error
	UpdateStatus(id, status string) error
	GetByNumber(number string) (*models.Room, error)
}

type roomRepository struct {
	db *sqlx.DB
}

func NewRoomRepository(db *sqlx.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) GetByID(id string) (*models.Room, error) {
	var room models.Room
	err := r.db.Get(&room, "SELECT * FROM rooms WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &room, err
}

func (r *roomRepository) GetAll() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Select(&rooms, "SELECT * FROM rooms WHERE status = 'ATIVO'")
	return rooms, err
}

func (r *roomRepository) Create(room *models.Room) error {
	query := `INSERT INTO rooms (id, number, type, capacity, price_per_night, status) VALUES (:id, :number, :type, :capacity, :price_per_night, :status)`
	_, err := r.db.NamedExec(query, room)
	return err
}

func (r *roomRepository) Update(room *models.Room) error {
	query := `UPDATE rooms SET type = :type, capacity = :capacity, price_per_night = :price_per_night WHERE id = :id`
	_, err := r.db.NamedExec(query, room)
	return err
}

func (r *roomRepository) UpdateStatus(id, status string) error {
	_, err := r.db.Exec("UPDATE rooms SET status = $1 WHERE id = $2", status, id)
	return err
}

func (r *roomRepository) GetByNumber(number string) (*models.Room, error) {
	var room models.Room
	err := r.db.Get(&room, "SELECT * FROM rooms WHERE number = $1 AND status = 'ATIVO'", number)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &room, err
}
