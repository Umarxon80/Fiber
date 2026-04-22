package main

import (
	"time"

	"github.com/Umarxon80/Fiber.git/db"
	"github.com/Umarxon80/Fiber.git/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

func main() {
	log.SetOutput(logger.Logger())             // logger set up
	cacheMiddleware := cache.New(cache.Config{ // caching
		Expiration: 10 * time.Second,
	})
	// db connection
	db.ConnectDb()
	db.Migrate()

	app := fiber.New(fiber.Config{AppName: "Fiber"})

	app.Get("/", compress.New(), cacheMiddleware, db.Get)
	app.Get("/:id", db.GetOne)
	app.Post("/", db.Add)
	app.Patch("/:id", db.Patch)
	app.Delete("/:id", db.Delete)
	log.Fatal(app.Listen(":8080"))
}
