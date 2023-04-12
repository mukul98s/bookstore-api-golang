package route

import (
	"bookstore/controller"

	"github.com/gofiber/fiber/v2"
)

func BookRoutes(router fiber.Router) {
	router.Get("/", controller.GetBooks)
	router.Post("/", controller.AddBook)
	router.Put("/:book_id", controller.UpdateBook)
	router.Delete("/:book_id", controller.DeleteBook)
}
