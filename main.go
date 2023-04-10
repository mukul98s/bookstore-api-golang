package main

import (
	"bookstore/helper"
	"bookstore/route"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	helper.CheckError(err, "Failed to load ENV file")
}

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
