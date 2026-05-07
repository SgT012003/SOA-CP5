package models

import (
	"time"
)

type Guest struct {
	ID        string    `db:"id" json:"id"`
	FullName  string    `db:"full_name" json:"full_name"`
	Document  string    `db:"document" json:"document"`
	Email     string    `db:"email" json:"email"`
	Phone     string    `db:"phone" json:"phone"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Room struct {
	ID            string  `db:"id" json:"id"`
	Number        string  `db:"number" json:"number"`
	Type          string  `db:"type" json:"type"`
	Capacity      int     `db:"capacity" json:"capacity"`
	PricePerNight float64 `db:"price_per_night" json:"price_per_night"`
	Status        string  `db:"status" json:"status"` // ATIVO, INATIVO
}

type Reservation struct {
	ID               string     `db:"id" json:"id"`
	GuestID          string     `db:"guest_id" json:"guest_id"`
	RoomID           string     `db:"room_id" json:"room_id"`
	CheckinExpected  time.Time  `db:"checkin_expected" json:"checkin_expected"`
	CheckoutExpected time.Time  `db:"checkout_expected" json:"checkout_expected"`
	CheckinAt        *time.Time `db:"checkin_at" json:"checkin_at"`
	CheckoutAt       *time.Time `db:"checkout_at" json:"checkout_at"`
	Status           string     `db:"status" json:"status"` // CREATED, CHECKED_IN, CHECKED_OUT, CANCELED
	EstimatedAmount  float64    `db:"estimated_amount" json:"estimated_amount"`
	FinalAmount      *float64   `db:"final_amount" json:"final_amount"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// DTOs
type CreateReservationRequest struct {
	GuestID          string    `json:"guest_id" binding:"required"`
	RoomID           string    `json:"room_id" binding:"required"`
	CheckinExpected  time.Time `json:"checkin_expected" binding:"required"`
	CheckoutExpected time.Time `json:"checkout_expected" binding:"required"`
	GuestsCount      int       `json:"guests_count" binding:"required,min=1"`
}

type UpdateReservationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=CANCELED CHECKED_IN CHECKED_OUT"`
}

type CreateGuestRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Document string `json:"document" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
}

type UpdateGuestRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type CreateRoomRequest struct {
	Number        string  `json:"number" binding:"required"`
	Type          string  `json:"type" binding:"required,oneof=STANDARD DELUXE SUITE"`
	Capacity      int     `json:"capacity" binding:"required,min=1"`
	PricePerNight float64 `json:"price_per_night" binding:"required,gt=0"`
}

type UpdateRoomRequest struct {
	Type          string  `json:"type" binding:"required,oneof=STANDARD DELUXE SUITE"`
	Capacity      int     `json:"capacity" binding:"required,min=1"`
	PricePerNight float64 `json:"price_per_night" binding:"required,gt=0"`
}
