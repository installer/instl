package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/installer/installer/internal/pkg/config"
	"github.com/installer/installer/internal/pkg/platforms"
	"github.com/installer/installer/scripts"
)

func installationWithConfig(cfg config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		owner := c.Params("user")
		repo := c.Params("repo")
		platform, err := platforms.Parse(c.Params("os"))
		if err != nil {
			return err
		}

		script, err := scripts.ParseTemplateForPlatform(platform, config.Combine(
			cfg,
			config.Config{
				Owner:     owner,
				Repo:      repo,
				CreatedAt: time.Now(),
				Version:   "latest",
			}),
		)
		if err != nil {
			return err
		}

		return c.SendString(script)
	}
}

func Installation(c *fiber.Ctx) error {
	return installationWithConfig(config.Config{})(c)
}

func InstallationVerbose(c *fiber.Ctx) error {
	return installationWithConfig(config.Config{Verbose: true})(c)
}
