package main

import (
	"github.com/gofiber/fiber/v2"
)

type Palindrome struct {
	Text string `json:"text"`
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
		return nil
	})

	app.Listen(":3000")

}
