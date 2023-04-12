package middleware

import (
	"bookstore/database"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("Auth")

	// cookie is not available
	if cookie == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"status":  false,
			"message": "Session Not Available...!",
		})
	}

	// validate the cookie
	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header)
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"status":  false,
			"message": "Invalid Session...!",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  false,
				"message": "Token Expried",
			})
		}

		id := claims["sub"]
		if id == "" {
			return ctx.JSON(&fiber.Map{
				"status":  false,
				"message": "User not found",
			})
		}

		var user_id string
		result := database.DB.QueryRow(`SELECT id FROM "users" where "id" = $1`, id).Scan(&user_id)

		if result != nil {
			return ctx.JSON(fiber.Map{
				"status":  false,
				"message": "Failed to Login",
			})
		}

		ctx.Set("user", user_id)
		return ctx.Next()
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
		"status":  false,
		"message": "Session Not Found",
	})
}
