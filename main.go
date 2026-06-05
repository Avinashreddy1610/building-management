package main

import (
	"log"

	"building-management/database"
	"building-management/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// init database with correct logins
	db, err := database.Connect("localhost", 5432, "postgres", "postgres", "bms")
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	defer db.Close()

	app := fiber.New()
	app.Use(logger.New())

	env := handlers.NewEnv(db)
	env.RegisterRoutes(app)

	log.Println("Starting Building Management System on port 3000...")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Server connection error: %v", err)
	}
}