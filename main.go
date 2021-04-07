package main

import (
	"log"
	"os"

	"github.com/NikSchaefer/go-fiber/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Set Env variables for
	// PORT=3000
	// DATABASE_URL="host=localhost port=5432 user=postgres password= dbname= sslmode=disable"
	// If deploying with Heroku these will be
	// automatically set provided you have created a db add on

	godotenv.Load()
	router := fiber.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*", // comma format e.g. "localhost, nikschaefer.tech"
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST",
	}))

	router.Use(func(c *fiber.Ctx) error {
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-DNS-Prefetch-Control", "off")
		return c.Next()
	})

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello, World!")
	})

	api.Initalize(router, db)

	router.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(router.Listen(":" + os.Getenv("PORT")))
}