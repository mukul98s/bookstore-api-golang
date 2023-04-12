package controller

import (
	"bookstore/database"
	"bookstore/helper"
	"bookstore/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Implement validation
func Login(ctx *fiber.Ctx) error {
	var err error
	body := new(model.User)

	// get the body
	err = ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(&fiber.Map{
			"status":  false,
			"message": "Cannot login without email and password",
		})
	}

	var userExists bool
	// check email
	err = database.DB.QueryRow(`SELECT EXISTS(SELECT 1 from "users" where "email"=$1)`, body.Email).Scan(&userExists)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Something Went Wrong",
		})
	}
	if !userExists {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Invalid Email or Password",
		})
	}

	// get User
	user := &model.User{}
	result := database.DB.QueryRow(`SELECT * FROM "users" where "email" = $1`, body.Email)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  false,
			"message": "Failed to Login",
		})
	}

	// the order of variable is similar to the order of table
	err = result.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Something Went Wrong",
		})
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Invalid Email or Password",
		})
	}

	// Set Cookie
	token, err := helper.GetTokens(user.Id)
	if err != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Login Failed....!",
		})
	}

	cookie := fiber.Cookie{
		Name:     "Auth",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   false,
		Path:     "/",
		SameSite: fiber.CookieSameSiteStrictMode,
	}
	ctx.Cookie(&cookie)

	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "User Logged in successfully",
		"data": &fiber.Map{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
		},
	})
}

// TODO: FIX Validation
func Signup(ctx *fiber.Ctx) error {
	var err error
	user := new(model.User)

	// get the body from request
	err = ctx.BodyParser(user)
	if err != nil {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(&fiber.Map{
			"status":  false,
			"message": "Please fill all the required details",
		})
	}

	// validate body with User sturct
	errors := helper.ValidateUser(*user)
	if errors != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid Data",
			"data":    errors,
		})
	}

	// check if email or phone number already existance
	checkExistingStmt := `SELECT EXISTS(SELECT 1 from "users" WHERE "email" = $1 OR "phone" = $2 )`
	var exists bool
	err = database.DB.QueryRow(checkExistingStmt, user.Email, user.Phone).Scan(&exists)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Something Went Wrong",
		})
	}
	if exists {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Email or Phone Number already Exists",
		})
	}

	// generate password hass
	passwordHash, err := helper.GenerateHash(user.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Something Went wrong",
		})
	}

	// create user
	var id string
	insertStmt := `INSERT INTO "users"("name", "email", "phone", "password") VALUES ($1, $2, $3, $4) RETURNING id`
	err = database.DB.QueryRow(insertStmt, user.Name, user.Email, user.Phone, passwordHash).Scan(&id)
	if err != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Failed to signin. Please Try again !",
			"error":   err.Error(),
		})
	}

	// Generate JWT and send cookie
	token, _ := helper.GetTokens(id)
	cookie := fiber.Cookie{
		Name:     "Auth",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   false,
		Path:     "/",
		SameSite: fiber.CookieSameSiteStrictMode,
	}
	ctx.Cookie(&cookie)

	// send back response
	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "User Created successfully",
		"data": &fiber.Map{
			"id": id,
		},
	})
}
