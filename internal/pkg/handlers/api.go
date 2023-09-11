package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.etcd.io/bbolt"
	"strconv"
)

func RepoStatsAPI(c *fiber.Ctx) error {
	owner, repo := getOwnerAndRepo(c)
	linux, windows, macos, err := getInstallationCountPerPlatform(owner, repo)
	if err != nil {
		return err
	}
	return c.JSON(map[string]any{
		"windows": windows,
		"linux":   linux,
		"macos":   macos,
		"total":   linux + windows + macos,
	})
}

func AllStatsAPI(c *fiber.Ctx) error {
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("installations"))
		if b == nil {
			return nil
		}
		stats := map[string]int{}
		err := b.ForEach(func(k, v []byte) error {
			i, err := strconv.Atoi(string(v))
			if err != nil {
				return err
			}
			stats[string(k)] = i
			return nil
		})
		if err != nil {
			return err
		}
		return c.JSON(stats)
	})
}

func AllStatsTotalAPI(c *fiber.Ctx) error {
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

		return c.SendString(strconv.Itoa(total))
	})
}
