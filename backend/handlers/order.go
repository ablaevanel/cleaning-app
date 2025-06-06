package handlers

import (
	"cleaning-app/db"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) error {
	userID := int(c.Locals("userID").(float64))

	var input struct {
		ServiceID int `json:"service_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверные данные"})
	}

	var orderID int
	err := db.DB.QueryRow(
		c.Context(),
		"INSERT INTO orders (user_id, service_id) VALUES ($1, $2) RETURNING id",
		userID, input.ServiceID,
	).Scan(&orderID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Не удалось создать заказ"})
	}

	return c.JSON(fiber.Map{"message": "Заказ создан", "order_id": orderID})
}

func GetAllOrders(c *fiber.Ctx) error {
	rows, err := db.DB.Query(c.Context(), `
		SELECT orders.id, users.name, services.name, orders.status
		FROM orders
		JOIN users ON orders.user_id = users.id
		JOIN services ON orders.service_id = services.id
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при получении заказов"})
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var userName, serviceName, status string
		rows.Scan(&id, &userName, &serviceName, &status)
		result = append(result, map[string]interface{}{
			"id":      id,
			"user":    userName,
			"service": serviceName,
			"status":  status,
		})
	}

	return c.JSON(result)
}

func UpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")
	var input struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверные данные"})
	}

	_, err := db.DB.Exec(c.Context(),
		"UPDATE orders SET status=$1 WHERE id=$2",
		input.Status, orderID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Не удалось обновить статус"})
	}

	return c.JSON(fiber.Map{"message": "Статус обновлён"})
}

func GetMyOrders(c *fiber.Ctx) error {
	userID := int(c.Locals("userID").(float64))

	rows, err := db.DB.Query(c.Context(), `
        SELECT 
            o.id, 
            s.name as service_name, 
            o.status, 
            o.created_at
        FROM orders o
        JOIN services s ON o.service_id = s.id
        WHERE o.user_id = $1
        ORDER BY o.created_at DESC
    `, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при получении заказов"})
	}
	defer rows.Close()

	type OrderResponse struct {
		ID        int       `json:"id"`
		Service   string    `json:"service"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	var orders []OrderResponse
	for rows.Next() {
		var order OrderResponse
		if err := rows.Scan(
			&order.ID,
			&order.Service,
			&order.Status,
			&order.CreatedAt,
		); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	if orders == nil {
		orders = []OrderResponse{}
	}

	return c.JSON(orders)
}

func DeleteOrder(c *fiber.Ctx) error {
	userID := int(c.Locals("userID").(float64))
	orderID := c.Params("id")

	var status string
	err := db.DB.QueryRow(
		c.Context(),
		"SELECT status FROM orders WHERE id = $1 AND user_id = $2",
		orderID, userID,
	).Scan(&status)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Заказ не найден или не принадлежит вам"})
	}

	if status != "created" {
		return c.Status(400).JSON(fiber.Map{"error": "Можно удалять только заказы со статусом 'created'"})
	}

	_, err = db.DB.Exec(
		c.Context(),
		"DELETE FROM orders WHERE id = $1",
		orderID,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при удалении заказа"})
	}

	return c.JSON(fiber.Map{"message": "Заказ успешно удален"})
}

func DeleteOrderAdmin(c *fiber.Ctx) error {
	orderID := c.Params("id")

	_, err := db.DB.Exec(c.Context(), "DELETE FROM orders WHERE id = $1", orderID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Не удалось удалить заказ"})
	}

	return c.JSON(fiber.Map{"message": "Заказ удалён"})
}
