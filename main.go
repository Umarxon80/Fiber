package main

import (
	"time"

	"github.com/Umarxon80/Fiber.git/auth"
	"github.com/Umarxon80/Fiber.git/db"
	costomLogger "github.com/Umarxon80/Fiber.git/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/earlydata"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/redirect"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/responsetime"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/timeout"
)

func main() {
	// logger set up
	log.SetOutput(costomLogger.Logger())

	// caching
	cacheMiddleware := cache.New(cache.Config{
		Expiration: 10 * time.Second,
	})
	// db connection
	db.Connect()
	defer db.DbConnection.Close()

	// generating application
	app := fiber.New(fiber.Config{AppName: "Fiber"})
	app.Use(recoverer.New(recoverer.Config{EnableStackTrace: true}))
	app.Use(requestid.New())
	app.Use(responsetime.New())

	// helmet - basic protection
	app.Use(helmet.New())

	// limiter
	app.Use(limiter.New(limiter.Config{
		Max:          20,
		Expiration:   5 * time.Minute,
		KeyGenerator: func(ctx fiber.Ctx) string { return ctx.IP() },
		LimitReached: func(ctx fiber.Ctx) error {
			log.Error("Too many requests user: ", ctx.IP())
			return ctx.Status(429).JSON(fiber.Map{
				"error": "Too many requests try later",
			})
		},
	}))

	// req logs
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} status:${status} - ${method} reqId: ${requestid}, time:${latency}\n",
		Stream: costomLogger.Logger(),
		CustomTags: map[string]logger.LogFunc{
			"requestid": func(output logger.Buffer, c fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString(requestid.FromContext(c))
			},
		},
	}))

	// earlydata
	app.Use(earlydata.New())

	// idempotency
	app.Use(idempotency.New(idempotency.Config{Lifetime: 10 * time.Second}))

	//favicon
	app.Use(favicon.New(favicon.Config{File: "./favicon.ico"}))

	// redirect
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/info": "/",
		},
		StatusCode: fiber.StatusMovedPermanently,
	}))

	// session
	app.Use(session.New(session.Config{
		IdleTimeout:     30 * time.Minute,
		AbsoluteTimeout: 24 * time.Hour,
		CookieSecure:    true,
		CookieHTTPOnly:  true,
		CookieSameSite:  "Lax",
	}))
	app.Post("/login", func(ctx fiber.Ctx) error {
		s := session.FromContext(ctx)
		s.Set("authenticated", true)
		return ctx.SendString("loged in")
	})

	// healthcheck
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	app.Get(healthcheck.StartupEndpoint, healthcheck.New())

	// handlers
	app.Get("/", compress.New(), cacheMiddleware, timeout.New(db.GetProducts, timeout.Config{Timeout: 1 * time.Minute}))
	app.Get("/:id", db.GetOneProduct)
	app.Post("/", auth.RequreAuth, db.CreateProduct)
	app.Patch("/:id", auth.RequreAuth, db.PatchProduct)
	app.Delete("/:id", auth.RequreAuth, db.DeleteProduct)
	log.Fatal(app.Listen(":8080"))
	
}
