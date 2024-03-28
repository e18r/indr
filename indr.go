package main

import (
	"strings"
	"regexp"
	"slices"
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type Palindrome struct {
	Text string `json:"text"`
}

var ascii = map[string]string{
  "á": "a",
  "à": "a",
  "â": "a",
  "ä": "a",
  "é": "e",
  "è": "e",
  "ê": "e",
  "ë": "e",
  "í": "i",
  "ì": "i",
  "î": "i",
  "ï": "i",
  "ó": "o",
  "ò": "o",
  "ô": "o",
  "ö": "o",
  "ú": "u",
  "ù": "u",
  "û": "u",
  "ü": "u",
}

func normalize(text string) string {
	text = strings.ToLower(text)
	for chr, asciiChr := range(ascii) {
		text = strings.ReplaceAll(text, chr, asciiChr)
	}
	palindromic, _ := regexp.Compile("[^a-z0-9ñçß]")
	text = palindromic.ReplaceAllString(text, "")
	return text
}

func isPalindrome(norm string) bool {
	chars := strings.Split(norm, "")
	slices.Reverse(chars)
	reverse := strings.Join(chars, "")
	return norm == reverse
}

func main() {

	app := fiber.New()

	app.Post("/publish", func(c *fiber.Ctx) error {
		palindrome := new(Palindrome)
		error := c.BodyParser(palindrome)
		if error != nil {
			return c.Status(fiber.StatusBadRequest).SendString(error.Error())
		}
		if palindrome.Text == "" {
			return c.Status(fiber.StatusBadRequest).
				SendString("Text param not found in request body")
		}
		norm := normalize(palindrome.Text)
		if !isPalindrome(norm) {
			return c.Status(fiber.StatusBadRequest).
				SendString("Not a palindrome")
		}
		connection, error := pgx.Connect(context.Background(), "postgresql://indr:indr@localhost:5432/palindr")
		if error != nil {
			log.Println(error)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer connection.Close(context.Background())
		var id int
		error = connection.QueryRow(context.Background(),
			"INSERT INTO norm (norm, created, attempts) VALUES ($1, CURRENT_TIMESTAMP, 1) ON CONFLICT ON CONSTRAINT norm_norm_key DO UPDATE SET attempts = norm.attempts + 1 RETURNING id", norm).
			Scan(&id)
		if error != nil {
			log.Println(error)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendString(strconv.Itoa(id))
	})

	app.Listen(":3000")

}
