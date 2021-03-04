package books

import (
	"fiberseed/database"
	"fiberseed/pkg"

	"github.com/gofiber/fiber/v2"
)

type Book struct {
	database.DefaultModel
	Title  string `json:"name"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetBooks(c *fiber.Ctx) error {
	db := database.DB
	var books []Book
	db.Find(&books)
	return c.JSON(books)
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var book Book
	db.Find(&book, id)
	return c.JSON(book)
}

func NewBook(c *fiber.Ctx) error {
	db := database.DB
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return pkg.BadRequest(err.Error())
	}
	db.Create(&book)
	return c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var book Book
	if err := db.First(&book, id); err != nil {
		return pkg.EntityNotFound("No book found")
	}
	db.Delete(&book)
	return c.SendStatus(204)
}
