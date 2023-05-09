package apis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", swagger.HandlerDefault)

	group := app.Group("/api")

	// User
	group.Post("/user/login", Login)
	group.Post("/user/register", Register)
	group.Post("/user/reset", Reset)

	// Box
	group.Get("/messageBoxes", ListBoxes)
	group.Get("/messageBox/:id", GetABox)
	group.Post("/messageBox", CreateABox)
	group.Put("/messageBox/:id", ModifyABox)
	group.Delete("/messageBox/:id", DeleteABox)

	// Post
	group.Get("/posts", ListPosts)
	group.Get("/post/:id", GetAPost)
	group.Post("/post", CreateAPost)
	group.Put("/post/:id", ModifyAPost)
	group.Delete("/post/:id", DeleteAPost)

	// Wall
	group.Get("/wall", ListWalls)
}
