package handlers

import (
	"cleaning-app/db"
	"cleaning-app/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Не удалось хешировать пароль"})
	}

	err = db.DB.QueryRow(
		c.Context(),
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, string(hashed),
	).Scan(&user.ID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка создания пользователя"})
	}

	return c.JSON(fiber.Map{"message": "Пользователь зарегистрирован"})
}

func Login(c *fiber.Ctx) error {
	type Credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var creds Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	var user models.User
	err := db.DB.QueryRow(
		c.Context(),
		"SELECT id, name, email, password, role FROM users WHERE email=$1",
		creds.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Пользователь не найден"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Неверный пароль"})
	}

	// Создание JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка генерации токена"})
	}

	return c.JSON(fiber.Map{"token": t})
}
