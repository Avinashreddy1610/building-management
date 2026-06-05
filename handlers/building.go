package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"building-management/models"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/gofiber/fiber/v2"
)

type buildingPayload struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (env *Env) GetBuildings(c *fiber.Ctx) error {
	ctx := c.UserContext()
	includeApts := c.Query("include_apartments") == "true"

	var mods []qm.QueryMod
	if includeApts {
		mods = append(mods, qm.Load(models.BuildingRels.Apartments))
	}

	list, err := models.Buildings(mods...).All(ctx, env.DB)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list buildings"})
	}

	return c.JSON(list)
}

func (env *Env) GetBuilding(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	b, err := models.FindBuilding(ctx, env.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "building not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal database error"})
	}

	return c.JSON(b)
}

func (env *Env) UpsertBuilding(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var payload buildingPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse request payload"})
	}

	if payload.Name == "" || payload.Address == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "name and address are required"})
	}

	var target *models.Building
	var err error

	if payload.ID > 0 {
		target, err = models.FindBuilding(ctx, env.DB, payload.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "database query failure"})
		}
	}

	if target == nil {
		target, err = models.Buildings(models.BuildingWhere.Name.EQ(payload.Name)).One(ctx, env.DB)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "database query failure"})
		}
	}

	if target != nil {
		target.Name = payload.Name
		target.Address = null.StringFrom(payload.Address)
		if _, err := target.Update(ctx, env.DB, boil.Infer()); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update building"})
		}
		return c.JSON(target)
	}

	newB := &models.Building{
		Name:    payload.Name,
		Address: null.StringFrom(payload.Address),
	}

	if err := newB.Insert(ctx, env.DB, boil.Infer()); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create building"})
	}

	return c.Status(http.StatusCreated).JSON(newB)
}

func (env *Env) DeleteBuilding(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	b, err := models.FindBuilding(ctx, env.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "building not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "database lookup failure"})
	}

	if _, err := b.Delete(ctx, env.DB); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete building"})
	}

	return c.JSON(fiber.Map{"message": "building successfully deleted"})
}