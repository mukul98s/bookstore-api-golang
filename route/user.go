package route

import (
	"bookstore/controller"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router) {
	router.Get("/", controller.GetUser)
	router.Put("/", controller.UpdateUser)
	router.Delete("/", controller.DeleteUser)
}
