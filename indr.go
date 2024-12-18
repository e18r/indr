package main

import (
	"strings"
	"regexp"
	"slices"
	"context"
	"log"
	"strconv"
	"os"
	"net"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

func dumpEnv() {
	for _, envVar := range os.Environ() {
		components := strings.Split(envVar, "=")
		if len(components[1]) > 0 {
			log.Println(components[0])
		}
	}
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

func getIP(IP string, IPs []string) string {
	if len(IPs) == 0 {
		return IP
	} else {
		for i := 0; i < len(IPs); i++ {
			if !net.ParseIP(IPs[i]).IsPrivate() {
				return IPs[i]
			}
		}
		return IP
	}
}

func store(
	c *fiber.Ctx,
	dbType string,
	url string,
	palindrome *Palindrome,
	norm string,
	IP string,
) error {
	conn, error := pgx.Connect(context.Background(), url)
	if error != nil {
		log.Println(error)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer conn.Close(context.Background())
	var normID int
	error = conn.QueryRow(context.Background(),
		"INSERT INTO norm (norm, attempts) " +
			"VALUES ($1, 0) " +
			"ON CONFLICT ON CONSTRAINT norm_norm_key " +
			"DO UPDATE SET attempts = norm.attempts + 1 " +
			"RETURNING id",
		norm).Scan(&normID)
	if error != nil {
		log.Println(error)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	var textID int
	error = conn.QueryRow(context.Background(),
		"INSERT INTO text (" +
			"text, origin, norm_id, created, attempts, seen, edited" +
			") VALUES ($1, $2, $3, CURRENT_TIMESTAMP, 0, 0, 0) " +
			"ON CONFLICT ON CONSTRAINT text_text_key " +
			"DO UPDATE SET attempts = text.attempts + 1 " +
			"RETURNING id",
		palindrome.Text, IP, normID).Scan(&textID)
	if error != nil {
		log.Println(error)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	log.Printf("[indr] publish [%s] (%s) \"%s\"\n",
		dbType,
		IP,
		palindrome.Text,
	)
	return c.SendString(strconv.Itoa(textID))
}

func storeAsync(
	c *fiber.Ctx,
	dbType string,
	palindrome *Palindrome,
	norm string,
	IP string,
) {
	url := os.Getenv("DATABASE_URL_2")
	if url == "" { return }
	conn, error := pgx.Connect(context.Background(), url)
	if error != nil {
		log.Println(error)
		return
	}
	defer conn.Close(context.Background())
	var normID int
	error = conn.QueryRow(context.Background(),
		"INSERT INTO norm (norm, attempts) " +
			"VALUES ($1, 0) " +
			"ON CONFLICT ON CONSTRAINT norm_norm_key " +
			"DO UPDATE SET attempts = norm.attempts + 1 " +
			"RETURNING id",
		norm).Scan(&normID)
	if error != nil {
		log.Println(error)
		return
	}
	var textID int
	error = conn.QueryRow(context.Background(),
		"INSERT INTO text (" +
			"text, origin, norm_id, created, attempts, seen, edited" +
			") VALUES ($1, $2, $3, CURRENT_TIMESTAMP, 0, 0, 0) " +
			"ON CONFLICT ON CONSTRAINT text_text_key " +
			"DO UPDATE SET attempts = text.attempts + 1 " +
			"RETURNING id",
		palindrome.Text, IP, normID).Scan(&textID)
	if error != nil {
		log.Println(error)
		return
	}
	log.Printf("[indr] publish [%s] (%s) \"%s\"\n",
		dbType,
		IP,
		palindrome.Text,
	)
	return
}

func main() {

	app := fiber.New()
	app.Use(cors.New())

	log.SetFlags(0)
	dumpEnv()

	app.Get("/", func(c *fiber.Ctx) error {
		return nil
	});

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
		if !isPalindrome(norm) || norm == "" {
			return c.Status(fiber.StatusBadRequest).
				SendString("Not a palindrome")
		}
		url := os.Getenv("DATABASE_URL")
		if url == "" {
			log.Println("env var DATABASE_URL is empty")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		IP := getIP(c.IP(), c.IPs())
		error = store(c, "primary", url, palindrome, norm, IP)
		if error != nil {
			log.Println(error)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		go storeAsync(c, "secondary", palindrome, norm, IP)
		return nil
	})

	app.Get("/list", func(c *fiber.Ctx) error {
		conn, err := pgx.Connect(context.Background(),
			os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer conn.Close(context.Background())
		after := c.QueryInt("after")
		rows, _ := conn.Query(context.Background(),
			"SELECT id, text FROM text " +
				"WHERE id > $1 " +
				"ORDER BY created ASC", after)
		var ID int
		var text string
		list := make([]map[string]string, 0, 50)
		_, err = pgx.ForEachRow(rows, []any{&ID, &text}, func() error {
			row := map[string]string{"id": strconv.Itoa(ID), "text": text}
			list = append(list, row)
			return nil
		})
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		jsonList, err := json.Marshal(list)
		if (err != nil) {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Set("Content-Type", "application/json")
		return c.SendString(string(jsonList))
	})

	app.Listen(":" + os.Getenv("PORT"))

}
