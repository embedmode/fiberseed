package books

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(route fiber.Router) {
	route.Get("/books", GetBooks)
	route.Get("/books/:id", GetBook)
	route.Put("/books/:id", UpdateBook)
	route.Post("/books", NewBook)
	route.Delete("/books/:id", DeleteBook)
}
