package routes

import (
	"git/internal/controllers"
	validation "git/internal/validator"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/user/signup", validation.ValidateUser, controllers.Register)
}
