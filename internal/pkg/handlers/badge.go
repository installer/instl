package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.etcd.io/bbolt"
	"strconv"
)

func RepoStatsShieldsIOBadge(c *fiber.Ctx) error {
	owner, repo := getOwnerAndRepo(c)
	linux, windows, macos, err := getInstallationCountPerPlatform(owner, repo)
	if err != nil {
		return err
	}

	return c.JSON(map[string]any{
		"schemaVersion": 1,
		"label":         "installations",
		"message":       strconv.Itoa(linux + windows + macos),
		"color":         "orange",
	})
}

func AllStatsTotalBadge(c *fiber.Ctx) error {
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("installations"))
		if b == nil {
			return nil
		}
		var total int
		err := b.ForEach(func(k, v []byte) error {
			tmp, err := strconv.Atoi(string(v))
			if err != nil {
				return err
			}
			total += tmp
			return nil
		})
		if err != nil {
			return err
		}

		return c.JSON(map[string]any{
			"schemaVersion": 1,
			"label":         "handled installations",
			"message":       strconv.Itoa(total),
			"color":         "orange",
		})
	})
}
