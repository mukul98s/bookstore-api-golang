package controller

import "github.com/gofiber/fiber/v2"

func GetBooks(ctx *fiber.Ctx) error {
	return ctx.JSON(&fiber.Map{
		"message": "Get Books",
	})
}

func AddBook(ctx *fiber.Ctx) error {
	return ctx.JSON(&fiber.Map{
		"message": "Add Books",
	})
}

func UpdateBook(ctx *fiber.Ctx) error {
	return ctx.JSON(&fiber.Map{
		"message": "Update Books",
	})
}

func DeleteBook(ctx *fiber.Ctx) error {
	return ctx.JSON(&fiber.Map{
		"message": "Delete Books",
	})
}
