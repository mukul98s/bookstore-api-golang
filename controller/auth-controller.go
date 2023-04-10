package controller

import (
	"bookstore/helper"
	"bookstore/model"

	"github.com/gofiber/fiber/v2"
)

func Login(ctx *fiber.Ctx) error {
	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "User Logged in successfully",
	})
}

func Signup(ctx *fiber.Ctx) error {
	user := new(model.User)

	err := ctx.BodyParser(user)
	helper.CheckError(err, "Failed to parse Body")

	errors := helper.ValidateUser(*user)
	if errors != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid Data",
			"data":    errors,
		})
	}

	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "User Created successfully",
		"data":    user,
	})
}
