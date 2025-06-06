package handlers

import (
	"cleaning-app/db"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateReview(c *fiber.Ctx) error {
	userID := int(c.Locals("userID").(float64))
	orderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID заказа"})
	}

	// Проверка что заказ завершен и принадлежит пользователю
	var status string
	err = db.DB.QueryRow(
		c.Context(),
		"SELECT status FROM orders WHERE id = $1 AND user_id = $2",
		orderID, userID,
	).Scan(&status)

	if err != nil || status != "completed" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Нельзя оставить отзыв на этот заказ",
		})
	}

	// Проверка что отзыв еще не оставлен
	var exists bool
	db.DB.QueryRow(
		c.Context(),
		"SELECT EXISTS(SELECT 1 FROM reviews WHERE order_id = $1)",
		orderID,
	).Scan(&exists)

	if exists {
		return c.Status(400).JSON(fiber.Map{
			"error": "Вы уже оставили отзыв на этот заказ",
		})
	}

	var input struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверные данные"})
	}

	if input.Rating < 1 || input.Rating > 5 {
		return c.Status(400).JSON(fiber.Map{"error": "Рейтинг должен быть от 1 до 5"})
	}

	_, err = db.DB.Exec(
		c.Context(),
		"INSERT INTO reviews (order_id, user_id, rating, comment) VALUES ($1, $2, $3, $4)",
		orderID, userID, input.Rating, input.Comment,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при создании отзыва"})
	}

	return c.JSON(fiber.Map{"message": "Отзыв успешно добавлен"})
}

func GetAllReviews(c *fiber.Ctx) error {
	rows, err := db.DB.Query(
		c.Context(),
		`SELECT 
            r.id, 
            r.rating, 
            r.comment, 
            r.created_at,
            u.name as user_name,
            s.name as service_name
         FROM reviews r
         JOIN users u ON r.user_id = u.id
         JOIN orders o ON r.order_id = o.id
         JOIN services s ON o.service_id = s.id
         ORDER BY r.created_at DESC`,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при получении отзывов"})
	}
	defer rows.Close()

	var reviews []map[string]interface{}
	for rows.Next() {
		var review struct {
			ID          int       `json:"id"`
			Rating      int       `json:"rating"`
			Comment     string    `json:"comment"`
			CreatedAt   time.Time `json:"created_at"`
			UserName    string    `json:"user_name"`
			ServiceName string    `json:"service_name"`
		}
		if err := rows.Scan(
			&review.ID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UserName,
			&review.ServiceName,
		); err != nil {
			continue
		}
		reviews = append(reviews, map[string]interface{}{
			"id":           review.ID,
			"rating":       review.Rating,
			"comment":      review.Comment,
			"created_at":   review.CreatedAt,
			"user_name":    review.UserName,
			"service_name": review.ServiceName,
		})
	}

	return c.JSON(reviews)
}
