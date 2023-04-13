package controller

import (
	"bookstore/database"
	"bookstore/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

// User Controller
func GetUser(ctx *fiber.Ctx) error {
	user_id := ctx.GetRespHeader("user")
	getStmt := `SELECT "id", "name", "email", "phone", "created_at", "updated_at" FROM "users" WHERE "id"=$1`
	result := database.DB.QueryRow(getStmt, user_id)
	if result.Err() != nil {
		ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Server Error",
		})
	}

	var user model.User
	err := result.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Failed to Get User Details",
		})
	}

	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "User Details Found",
		"data": &fiber.Map{
			"id":         user.Id,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

func UpdateUser(ctx *fiber.Ctx) error {
	user_id := ctx.GetRespHeader("user")
	var err error

	body := new(model.User)
	err = ctx.BodyParser(&body)
	if err != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Required Data is missing!",
		})
	}

	// check email exists
	checkExistingStmt := `SELECT EXISTS(SELECT 1 from "users" WHERE "email" = $1 AND "id" != $2)`
	var isEmailTaken bool
	err = database.DB.QueryRow(checkExistingStmt, body.Email, user_id).Scan(&isEmailTaken)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Something Went Wrong",
		})
	}
	if isEmailTaken {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Email Address is already used",
		})
	}

	// check phone exists
	phoneExistingStmt := `SELECT EXISTS(SELECT 1 from "users" WHERE "phone" = $1 AND "id" != $2)`
	var isPhoneTaken bool
	err = database.DB.QueryRow(phoneExistingStmt, body.Phone, user_id).Scan(&isPhoneTaken)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Something Went Wrong",
		})
	}
	if isPhoneTaken {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Phone Number is already used",
		})
	}

	updateStmt := `UPDATE "users" SET "name"=$1, "email"=$2, "phone"=$3, "updated_at"=$4 WHERE "id"=$5`
	_, updateErr := database.DB.Exec(updateStmt, body.Name, body.Email, body.Phone, time.Now(), user_id)
	if updateErr != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Failed to update user",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  true,
		"message": "User Profile Updated",
	})
}

func DeleteUser(ctx *fiber.Ctx) error {
	user_id := ctx.GetRespHeader("user")

	_, err := database.DB.Exec(`DELETE FROM "users" WHERE "id" = $1`, user_id)
	if err != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Failed to Delete Profile",
		})
	}

	// delete the cookie by expiring it.
	cookie := fiber.Cookie{
		Name:     "Auth",
		HTTPOnly: true,
		Expires:  time.Now(),
		Secure:   false,
		Path:     "/",
		SameSite: fiber.CookieSameSiteStrictMode,
	}
	ctx.Cookie(&cookie)

	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "Delete user",
	})
}
