package handlers

import "github.com/gofiber/fiber/v2"

func MissingPlatform(c *fiber.Ctx) error {
	c.Status(400)
	return c.SendString("error: missing platform - please specify a platform in the URL (e.g. instl.sh/owner/repo/linux)")
}
