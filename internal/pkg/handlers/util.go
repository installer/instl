package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func getOwnerAndRepo(c *fiber.Ctx) (owner, repo string) {
	owner = c.Params("user")
	owner = strings.ToLower(owner)
	repo = c.Params("repo")
	repo = strings.ToLower(repo)
	return owner, repo
}
