package books

import (
	"bytes"
	"encoding/json"
	"fiberseed/database"
	"fiberseed/server"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app = server.Create()
	database.DB.AutoMigrate(&Book{})
	Routes(app)

	// Cleanup books
	database.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Book{})

	exitVal := m.Run()

	// Cleanup books
	database.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Book{})

	os.Exit(exitVal)
}

func TestGetBooks(t *testing.T) {
	req := httptest.NewRequest("GET", "/books", nil)
	resp, _ := app.Test(req)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode, "status ok")
	assert.Equal(t, string(body), "[]", "empty body")

	book := &Book{Title: "The Name of the Wind: The Kingkiller Chronicle", Author: "Patrick Rothfuss", Rating: 10}
	database.DB.Create(book)

	req = httptest.NewRequest("GET", "/books", nil)
	resp, _ = app.Test(req)
	body, _ = ioutil.ReadAll(resp.Body)

	var books []Book
	err := json.Unmarshal(body, &books)
	assert.Equal(t, err, nil)

	assert.Equal(t, books[0].Title, book.Title)
	assert.Equal(t, books[0].ID, book.ID)
	assert.Equal(t, books[0].Author, book.Author)
	assert.Equal(t, books[0].Rating, book.Rating)
}

func TestGetBook(t *testing.T) {
	req := httptest.NewRequest("GET", "/books/1", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 404, resp.StatusCode, "status ok")

	req = httptest.NewRequest("GET", "/books/foo", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 500, resp.StatusCode, "status ok")

	newBook := &Book{Title: "The Wise Man's Fear", Author: "Patrick Rothfuss", Rating: 10}
	database.DB.Create(newBook)

	req = httptest.NewRequest("GET", fmt.Sprintf("/books/%d", newBook.ID), nil)
	resp, _ = app.Test(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var book Book
	err := json.Unmarshal(body, &book)
	assert.Equal(t, err, nil)

	assert.Equal(t, book.Title, newBook.Title)
	assert.Equal(t, book.ID, newBook.ID)
	assert.Equal(t, book.Author, newBook.Author)
	assert.Equal(t, book.Rating, newBook.Rating)
}

func TestNewBook(t *testing.T) {
	newBook := map[string]interface{}{
		"title":  "The Name of the Wind: The Kingkiller Chronicle",
		"author": "Patrick Rothfuss",
		"rating": 10,
	}
	body, _ := json.Marshal(newBook)
	req := httptest.NewRequest("POST", "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode, "status ok")

	var book Book
	err := json.Unmarshal(body, &book)
	assert.Equal(t, err, nil)

	assert.NotEqual(t, book.ID, nil)
	assert.Equal(t, book.Title, newBook["title"])
	assert.Equal(t, book.Author, newBook["author"])
	assert.Equal(t, book.Rating, newBook["rating"])
}

func TestUpdateBook(t *testing.T) {
	book := &Book{Title: "The Wise Man's Fear", Author: "Patrick Rothfuss", Rating: 10}
	database.DB.Create(book)

	newTitle := "Foo"
	body, _ := json.Marshal(map[string]interface{}{
		"title":  newTitle,
		"author": book.Author,
		"rating": book.Rating,
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/books/%d", book.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, 204, resp.StatusCode, "status ok")

	req = httptest.NewRequest("GET", fmt.Sprintf("/books/%d", book.ID), nil)
	resp, _ = app.Test(req)
	body, _ = ioutil.ReadAll(resp.Body)

	var updatedBook Book
	err := json.Unmarshal(body, &updatedBook)
	assert.Equal(t, err, nil)

	assert.Equal(t, updatedBook.Title, newTitle)
	assert.Equal(t, book.ID, updatedBook.ID)
	assert.Equal(t, book.Author, updatedBook.Author)
	assert.Equal(t, book.Rating, updatedBook.Rating)
}

func TestDeleteBook(t *testing.T) {
	req := httptest.NewRequest("GET", "/books/0", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 404, resp.StatusCode, "status ok")

	req = httptest.NewRequest("GET", "/books/foo", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 500, resp.StatusCode, "status ok")

	newBook := &Book{Title: "The Wise Man's Fear", Author: "Patrick Rothfuss", Rating: 10}
	database.DB.Create(newBook)

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/books/%d", newBook.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 204, resp.StatusCode, "status ok")
}
