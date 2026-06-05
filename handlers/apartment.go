package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"building-management/models"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/gofiber/fiber/v2"
)

type apartmentPayload struct {
	ID         int    `json:"id"`
	BuildingID int    `json:"building_id"`
	Number     string `json:"number"`
	Floor      int    `json:"floor"`
	SqMeters   int    `json:"sq_meters"`
}

func (env *Env) GetApartments(c *fiber.Ctx) error {
	ctx := c.UserContext()
	list, err := models.Apartments().All(ctx, env.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve apartments"})
	}
	return c.JSON(list)
}

func (env *Env) GetApartment(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	apt, err := models.FindApartment(ctx, env.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "apartment not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal database error"})
	}

	return c.JSON(apt)
}

func (env *Env) GetApartmentsByBuilding(c *fiber.Ctx) error {
	ctx := c.UserContext()
	buildingID, err := strconv.Atoi(c.Params("buildingId"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid building ID format"})
	}

	exists, err := models.BuildingExists(ctx, env.DB, buildingID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check building existence"})
	}
	if !exists {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "referenced building not found"})
	}

	list, err := models.Apartments(models.ApartmentWhere.BuildingID.EQ(buildingID)).All(ctx, env.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve building apartments"})
	}

	return c.JSON(list)
}

func (env *Env) UpsertApartment(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var payload apartmentPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse request payload"})
	}

	if payload.BuildingID == 0 || payload.Number == "" || payload.Floor == 0 || payload.SqMeters == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "building_id, number, floor, and sq_meters are required",
		})
	}

	buildingExists, err := models.BuildingExists(ctx, env.DB, payload.BuildingID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "validation query failed"})
	}
	if !buildingExists {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "referenced building does not exist"})
	}

	var apt *models.Apartment

	if payload.ID > 0 {
		apt, err = models.FindApartment(ctx, env.DB, payload.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "database query failure"})
		}
	}

	if apt == nil {
		apt, err = models.Apartments(
			models.ApartmentWhere.BuildingID.EQ(payload.BuildingID),
			models.ApartmentWhere.Number.EQ(payload.Number),
		).One(ctx, env.DB)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "database query failure"})
		}
	}

	if apt != nil {
		apt.BuildingID = payload.BuildingID
		apt.Number = payload.Number
		apt.Floor = payload.Floor
		apt.SQMeters = payload.SqMeters
		if _, err := apt.Update(ctx, env.DB, boil.Infer()); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update apartment"})
		}
		return c.JSON(apt)
	}

	newApt := &models.Apartment{
		BuildingID: payload.BuildingID,
		Number:     payload.Number,
		Floor:      payload.Floor,
		SQMeters:   payload.SqMeters,
	}

	if err := newApt.Insert(ctx, env.DB, boil.Infer()); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create apartment"})
	}

	return c.Status(http.StatusCreated).JSON(newApt)
}

func (env *Env) DeleteApartment(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	apt, err := models.FindApartment(ctx, env.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "apartment not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "database lookup failure"})
	}

	if _, err := apt.Delete(ctx, env.DB); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete apartment"})
	}

	return c.JSON(fiber.Map{"message": "apartment successfully deleted"})
}