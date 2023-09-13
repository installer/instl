package handlers

import "github.com/gofiber/fiber/v2"

func PrivacyPolicy(c *fiber.Ctx) error {
	return c.Render("privacy.gohtml", fiber.Map{})
}

func TermsOfService(c *fiber.Ctx) error {
	return c.Render("tos.gohtml", fiber.Map{})
}
