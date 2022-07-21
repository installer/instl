package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pterm/pterm"
	"go.etcd.io/bbolt"

	"github.com/installer/installer/internal/pkg/config"
	"github.com/installer/installer/internal/pkg/platforms"
	"github.com/installer/installer/scripts"
)

var db *bbolt.DB

func init() {
	dbTmp, err := bbolt.Open("./metrics.db", 0666, nil)
	pterm.Fatal.PrintOnError(err)
	db = dbTmp
}

func installationWithConfig(cfg config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		owner, repo := getOwnerAndRepo(c)
		platform, err := platforms.Parse(c.Params("os"))
		if err != nil {
			return err
		}

		script, err := scripts.ParseTemplateForPlatform(platform, config.Combine(
			config.Config{
				Owner:     owner,
				Repo:      repo,
				CreatedAt: time.Now(),
				Version:   "latest",
			},
			cfg,
		))
		if err != nil {
			return err
		}

		err = c.SendString(script)

		key := owner + "/" + repo + "/" + platform.String()
		if owner != "status-health-check" {
			// Dont return error, user is not affected by errors that happen in stats collection
			db.Update(func(tx *bbolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("installations"))
				if err != nil {
					return err
				}
				var installCountString string
				if v := b.Get([]byte(key)); v != nil {
					installCountString = string(v)
				}
				var installCount int
				if installCountString != "" {
					installCount, err = strconv.Atoi(installCountString)
					if err != nil {
						return err
					}
				}
				installCount++
				return b.Put([]byte(key), []byte(strconv.Itoa(installCount)))
			})
		}

		return err
	}
}

func Installation(c *fiber.Ctx) error {
	return installationWithConfig(config.Config{})(c)
}

func InstallationVerbose(c *fiber.Ctx) error {
	return installationWithConfig(config.Config{Verbose: true})(c)
}
