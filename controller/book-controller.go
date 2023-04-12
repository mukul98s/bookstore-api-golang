package controller

import (
	"bookstore/database"
	"bookstore/model"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetBooks(ctx *fiber.Ctx) error {
	user_id := ctx.GetRespHeader("user")
	books := make([]model.Book, 0)

	result, err := database.DB.Query(`SELECT * FROM "books" WHERE "user_id" = $1`, user_id)
	fmt.Println("User ID is", user_id)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(&fiber.Map{
			"status":  false,
			"message": "Failed to load Books",
			"error":   err.Error(),
		})
	}

	// create a book slice and append values after scanning each result
	for result.Next() {
		var book model.Book
		err := result.Scan(&book.Id, &book.Name, &book.CreatedAt, &book.AuthorName, &book.UpdatedAt, &book.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"status":  false,
				"message": "Something Went Wrong",
			})
		}

		books = append(books, book)
	}

	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "Get Books",
		"books":   books,
	})
}

func AddBook(ctx *fiber.Ctx) error {
	user_id := ctx.GetRespHeader("user")
	book := new(model.Book)
	var err error

	err = ctx.BodyParser(&book)
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(&fiber.Map{
			"status":   false,
			"messgage": "Required Data is missing",
		})
	}

	book.UserID = user_id

	// TODO: Add Validation

	// add the book to the database
	var id string
	insertSmt := `INSERT INTO "books" ("user_id", "name", "author_name") VALUES ($1, $2, $3) RETURNING id`
	err = database.DB.QueryRow(insertSmt, book.UserID, book.Name, book.AuthorName).Scan(&id)
	if err != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Something went wrong while saving the book",
		})
	}

	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "Add Books",
		"data": &fiber.Map{
			"id":          id,
			"name":        book.Name,
			"author_name": book.AuthorName,
		},
	})
}

func UpdateBook(ctx *fiber.Ctx) error {
	user_id := ctx.GetRespHeader("user")
	book_id := ctx.Params("book_id")

	var err error

	// Check the existance of the book
	var bookExists bool
	err = database.DB.QueryRow(`SELECT EXISTS(SELECT id FROM "books" WHERE "id" = $1 and "user_id"= $2)`, book_id, user_id).Scan(&bookExists)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Server Error",
		})
	}
	if !bookExists {
		return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"status":  false,
			"message": "This book does not exist",
		})
	}

	// get the body
	body := new(model.Book)
	err = ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(&fiber.Map{
			"status":  false,
			"message": "Required data is missing",
		})
	}

	// update the details
	updateStatement := `UPDATE "books" SET "name"=$1, "author_name"=$2, "updated_at"=$3 WHERE "id" = $4 AND "user_id"=$5`
	_, updateErr := database.DB.Exec(updateStatement, body.Name, body.AuthorName, time.Now(), book_id, user_id)
	if updateErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Something Went Wrong. Update Failed",
		})
	}

	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "Book Updated Successfully",
	})
}

func DeleteBook(ctx *fiber.Ctx) error {
	book_id := ctx.Params("book_id")
	user_id := ctx.GetRespHeader("user")

	var err error

	// check book with this id exists
	var bookExists bool
	err = database.DB.QueryRow(`SELECT EXISTS(SELECT id FROM "books" WHERE "id" = $1 and "user_id"= $2)`, book_id, user_id).Scan(&bookExists)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  false,
			"message": "Server Error",
		})
	}
	if !bookExists {
		return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"status":  false,
			"message": "This book does not exist",
		})
	}

	// delete the book
	_, deleteErr := database.DB.Exec(`DELETE FROM "books" WHERE "id"=$1 AND "user_id"=$2`, book_id, user_id)
	if deleteErr != nil {
		return ctx.JSON(&fiber.Map{
			"status":  false,
			"message": "Failed to Delete Book. ",
		})
	}

	return ctx.JSON(&fiber.Map{
		"status":  true,
		"message": "Delete Books",
		"data": &fiber.Map{
			"book_id": book_id,
		},
	})
}
