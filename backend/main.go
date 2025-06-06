package main

import (
	"cleaning-app/db"
	"cleaning-app/handlers"
	"cleaning-app/middleware"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Cleaning App API",
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	// Logger middleware
	app.Use(logger.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     db.IsConnected(),
		})
	})

	db.Connect()

	// API routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	app.Get("/services", handlers.GetAllServices)
	app.Get("/reviews", handlers.GetAllReviews)

	app.Post("/orders", middleware.Protected(), handlers.CreateOrder)
	app.Get("/orders", middleware.Protected(), handlers.GetMyOrders)
	app.Delete("/orders/:id", middleware.Protected(), handlers.DeleteOrder)
	app.Post("/orders/:id/reviews", middleware.Protected(), handlers.CreateReview)

	app.Get("/admin/orders", middleware.Protected(), middleware.AdminOnly(), handlers.GetAllOrders)
	app.Patch("/admin/orders/:id", middleware.Protected(), middleware.AdminOnly(), handlers.UpdateOrderStatus)
	app.Delete("/admin/orders/:id", middleware.Protected(), middleware.AdminOnly(), handlers.DeleteOrder)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen(":" + port))
}
