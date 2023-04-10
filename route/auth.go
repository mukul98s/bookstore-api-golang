package route

import (
	"bookstore/controller"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router) {
	router.Post("/login", controller.Login)
	router.Post("/signup", controller.Signup)
}
