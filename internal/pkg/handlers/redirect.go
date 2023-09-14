package handlers

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func RedirectToDocs(c *fiber.Ctx) error {
	return c.Redirect("https://docs.instl.sh")
}

func RedirectLegacyStats(c *fiber.Ctx) error {
	// Legacy stats link: https://instl.sh/stats/user/repo
	// New stats link: https://instl.sh/user/repo

	newUrl := strings.ReplaceAll(c.OriginalURL(), "/stats", "")
	return c.Redirect(newUrl)
}

func Redirect(url string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Redirect(url)
	}
}
