package main

import (
	"log"
	"os"

	"hotel-api/internal/handlers"
	"hotel-api/internal/middlewares"
	"hotel-api/internal/repositories"
	"hotel-api/internal/services"
	"hotel-api/pkg/database"

	_ "hotel-api/docs" // Swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Hotel Reservation API
// @version 1.0
// @description API for managing hotel reservations.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://hotel_user:hotel_password@localhost:5432/hotel_db?sslmode=disable"
	}

	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Repositories
	guestRepo := repositories.NewGuestRepository(db)
	roomRepo := repositories.NewRoomRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)

	// Services
	guestService := services.NewGuestService(guestRepo)
	roomService := services.NewRoomService(roomRepo)
	reservationService := services.NewReservationService(reservationRepo, roomRepo, guestRepo)

	// Handlers
	guestHandler := handlers.NewGuestHandler(guestService)
	roomHandler := handlers.NewRoomHandler(roomService)
	reservationHandler := handlers.NewReservationHandler(reservationService)

	// Router
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	v1 := r.Group("/api/v1")
	{
		guests := v1.Group("/guests")
		{
			guests.POST("", guestHandler.Create)
			guests.GET("", guestHandler.GetAll)
			guests.GET("/:id", guestHandler.GetByID)
			guests.PUT("/:id", guestHandler.Update)
			guests.DELETE("/:id", guestHandler.Delete)
		}

		rooms := v1.Group("/rooms")
		{
			rooms.POST("", roomHandler.Create)
			rooms.GET("", roomHandler.GetAll)
			rooms.GET("/:id", roomHandler.GetByID)
			rooms.PUT("/:id", roomHandler.Update)
			rooms.DELETE("/:id", roomHandler.Delete)
		}

		reservations := v1.Group("/reservations")
		{
			reservations.POST("", reservationHandler.CreateReservation)
			reservations.PATCH("/:id/status", reservationHandler.UpdateStatus)
		}
	}

	log.Println("Server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
