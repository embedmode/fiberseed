package books

import (
	"errors"
	"fiberseed/database"
	"fiberseed/pkg"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Book struct {
	database.DefaultModel
	Title  string `json:"title"`
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
	err := db.First(&book, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return pkg.EntityNotFound("No book found")
	} else if err != nil {
		return pkg.Unexpected(err.Error())
	}

	return c.JSON(book)
}

func NewBook(c *fiber.Ctx) error {
	db := database.DB
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return pkg.BadRequest("Invalid params")
	}
	db.Create(&book)
	return c.JSON(book)
}

func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var book Book
	err := db.First(&book, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return pkg.EntityNotFound("No book found")
	} else if err != nil {
		return pkg.Unexpected(err.Error())
	}

	updatedBook := new(Book)

	if err := c.BodyParser(updatedBook); err != nil {
		return pkg.BadRequest("Invalid params")
	}

	updatedBook = &Book{Title: updatedBook.Title, Author: updatedBook.Author, Rating: updatedBook.Rating}

	if err = db.Model(&book).Updates(updatedBook).Error; err != nil {
		return pkg.Unexpected(err.Error())
	}

	return c.SendStatus(204)
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var book Book
	err := db.First(&book, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return pkg.EntityNotFound("No book found")
	} else if err != nil {
		return pkg.Unexpected(err.Error())
	}

	db.Delete(&book)
	return c.SendStatus(204)
}
