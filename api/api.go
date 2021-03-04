package api

import (
	"fiberseed/api/books"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	v1 := app.Group("/api/v1")
	books.Routes(v1)
}
