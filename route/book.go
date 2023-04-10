package route

import (
	"bookstore/controller"

	"github.com/gofiber/fiber/v2"
)

func BookRoutes(router fiber.Router) {
	router.Get("/", controller.GetBooks)
	router.Post("/", controller.AddBook)
	router.Put("/", controller.UpdateBook)
	router.Delete("/", controller.DeleteBook)
}
