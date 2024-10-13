package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/raghav1030/go-ws-chat/cmd/internal/chat"
	"github.com/raghav1030/go-ws-chat/cmd/internal/handlers"
)

func main() {
	app := fiber.New()
	go chat.Manager.Start()

	app.Get("/api/v1/ws/register/:nick", websocket.New(handlers.RegisterHandler))

	app.Get("/api/v1", sampleHandler)

	app.Listen(":8080")
}

func sampleHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"ping": "pong",
	})
}
