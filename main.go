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

	app := fiber.New(fiber.Config{
		AppName: "Fiber",
	})

	// Get
	app.Get("/", compress.New(), cacheMiddleware, func(c fiber.Ctx) error {
		log.Info("Products are shown")
		return c.JSON(db.Get())
	})

	// Post
	app.Post("/", func(c fiber.Ctx) error {
		var p db.Product
		if err := c.Bind().JSON(&p); err != nil {
			log.Errorf("Wrong input: %v", err)
			return err
		}

		log.Info("Product is added")
		return c.JSON(db.Add(p))
	})

	// Patch
	app.Patch("/:id", func(c fiber.Ctx) error {
		id := fiber.Params[int](c, "id")
		var p db.Product
		if err := c.Bind().JSON(&p); err != nil {
			log.Errorf("Wrong input: %v", err)
			return err
		}

		pr, err := db.Patch(p, id)
		if err != nil {
			log.Errorf("Wrong input %v ", err)
			return err
		}
		return c.JSON(pr)
	})

	// Delete
	app.Delete("/:id", func(c fiber.Ctx) error {
		id := fiber.Params[int](c, "id")

		if err := db.Delete(id); err != nil {
			log.Errorf("Wrong input %v ", err)
			return err
		}
		return c.SendString("Product is deleted")
	})

	log.Fatal(app.Listen(":8080"))
}
