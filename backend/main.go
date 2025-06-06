package main

import (
	"cleaning-app/db"
	"cleaning-app/handlers"
	"cleaning-app/middleware"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize database connection
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "Cleaning App API",
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           300,
	}))

	// Logger middleware
	app.Use(logger.New())

	// Recover middleware
	app.Use(recover.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     db.IsConnected(),
			"time":   time.Now().Format(time.RFC3339),
		})
	})

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

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
