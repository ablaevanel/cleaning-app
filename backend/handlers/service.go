package handlers

import (
	"cleaning-app/db"
	"cleaning-app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateService(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	if claims["role"] != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Доступ запрещён"})
	}

	var s models.Service
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	err := db.DB.QueryRow(
		c.Context(),
		"INSERT INTO services (name, description, price, duration_minutes, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		s.Name, s.Description, s.Price, s.DurationMinutes, s.ImageURL,
	).Scan(&s.ID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка создания услуги"})
	}

	return c.JSON(s)
}

func GetAllServices(c *fiber.Ctx) error {
	rows, err := db.DB.Query(c.Context(), "SELECT id, name, description, price, duration_minutes, image_url FROM services")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения услуг"})
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Price, &s.DurationMinutes, &s.ImageURL); err != nil {
			continue
		}
		services = append(services, s)
	}

	if err = rows.Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при чтении данных"})
	}

	return c.JSON(services)
}
