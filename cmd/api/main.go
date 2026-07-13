package main

import (
	"log"
	"os"

	_ "git/docs"
	"git/internal/database"
	"git/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title Realtime API
// @version 1.0
// @description Realtime backend API documentation.
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if !loadEnv() {
		log.Println(".env dosyasi bulunamadi, sistem environment degerleri kullanilacak")
	}

	database.Connect()

	app := fiber.New()
	app.Get("/", healthHandler)
	routes.SetupRoutes(app)
	app.Get("/swagger/*", swagger.HandlerDefault)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

func loadEnv() bool {
	paths := []string{".env", "../../.env"}
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			return true
		}
	}
	return false
}

// healthHandler godoc
// @Summary API health check
// @Description Returns a simple welcome message when the API is running.
// @Tags health
// @Produce plain
// @Success 200 {string} string "Welcome to the API"
// @Router / [get]
func healthHandler(c *fiber.Ctx) error {
	return c.SendString("Welcome to the API")
}
