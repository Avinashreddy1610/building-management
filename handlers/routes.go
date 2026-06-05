package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// Env hold our db connection for handlers
type Env struct {
	DB *sql.DB
}

// NewEnv make new env with db injected inside
func NewEnv(db *sql.DB) *Env {
	return &Env{DB: db}
}

// RegisterRoutes setup all the url routes
func (env *Env) RegisterRoutes(app *fiber.App) {

	buildings := app.Group("/buildings")
	buildings.Get("/", env.GetBuildings)
	buildings.Get("/:id", env.GetBuilding)
	buildings.Post("/", env.UpsertBuilding)
	buildings.Delete("/:id", env.DeleteBuilding)

	apartments := app.Group("/apartments")
	apartments.Get("/", env.GetApartments)
	apartments.Get("/:id", env.GetApartment)
	apartments.Get("/building/:buildingId", env.GetApartmentsByBuilding)
	apartments.Post("/", env.UpsertApartment)
	apartments.Delete("/:id", env.DeleteApartment)
}