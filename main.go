package main

import (
	"bookstore/route"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	api := fiber.New()
	app := api.Group("/api")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app.Route("/auth", route.AuthRoutes)
	app.Route("/books", route.BookRoutes)
	app.Route("/user", route.UserRoutes)

	panic(api.Listen(fmt.Sprintf(":%v", port)))
}
