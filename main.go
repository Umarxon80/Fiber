package main

import (
    "log"
    "github.com/gofiber/fiber/v3"
	"github.com/Umarxon80/Fiber.git/db"
)

func main() {
    app := fiber.New(fiber.Config{
		AppName: "Fiber",
	})

	// Get
    app.Get("/", func (c fiber.Ctx) error {
        return c.JSON(db.Get())
    })

	// Post
	app.Post("/", func (c fiber.Ctx) error {
		var p db.Product
		if err:=c.Bind().JSON(&p);err!=nil {
			return err
		}

        return c.JSON(db.Add(p))
    })

	// Patch
	app.Patch("/:id", func (c fiber.Ctx) error {
		id:=fiber.Params[int](c,"id")
		var p db.Product
		if err:=c.Bind().JSON(&p);err!=nil {
			return err
		}
        return c.JSON(db.Putch(p,id))
    })

	// Delete
	app.Delete("/:id", func (c fiber.Ctx) error {
		id:=fiber.Params[int](c,"id")
		return c.JSON(db.Delete(int(id)))
    })

    log.Fatal(app.Listen(":8080"))
}