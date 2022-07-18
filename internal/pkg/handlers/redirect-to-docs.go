package handlers

import "github.com/gofiber/fiber/v2"

func RedirectToDocs(c *fiber.Ctx) error {
	return c.Redirect("https://docs.instl.sh")
}
