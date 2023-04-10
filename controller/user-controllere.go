package controller

import (
	"bookstore/database"
	"bookstore/helper"
	"bookstore/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// User Controller
func GetUser(ctx *fiber.Ctx) error {
	rows, err := database.DB.Query("Select * from users")
	// defer database.DB.Close()
	helper.CheckError(err, "Failed to get user information")

	fmt.Println(rows)

	return ctx.JSON(&model.User{
		Id:    "1",
		Name:  "Mukul Sharma",
		Email: "mymukul@112@gmail.com",
		Phone: "9610898865",
	})
}

func UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("user_id")
	fmt.Println(id)

	return ctx.JSON(fiber.Map{
		id: id,
	})
}

func DeleteUser(ctx *fiber.Ctx) error {
	return ctx.JSON(&fiber.Map{
		"message": "Delete user",
	})
}
