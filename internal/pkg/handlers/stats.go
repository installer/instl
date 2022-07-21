package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.etcd.io/bbolt"
)

func RepoStats(c *fiber.Ctx) error {
	user := c.Params("user")
	repo := c.Params("repo")
	linux, windows, macos, err := getInstallationCountPerPlatform(user, repo)
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

func getInstallationCountPerPlatform(user, repo string) (linux, windows, macos int, err error) {
	keyBase := user + "/" + repo + "/"
	keyWindows := keyBase + "windows"
	keyLinux := keyBase + "linux"
	keyMacOS := keyBase + "macos"
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("installations"))
		if b == nil {
			return nil
		}
		windowsInstallationsRaw := b.Get([]byte(keyWindows))
		linuxInstallationsRaw := b.Get([]byte(keyLinux))
		macosInstallationsRaw := b.Get([]byte(keyMacOS))

		windowsInstallations, _ := strconv.Atoi(string(windowsInstallationsRaw))
		windows = windowsInstallations
		linuxInstallations, _ := strconv.Atoi(string(linuxInstallationsRaw))
		linux = linuxInstallations
		macosInstallations, _ := strconv.Atoi(string(macosInstallationsRaw))
		macos = macosInstallations
		return nil
	})
	if err != nil {
		return 0, 0, 0, err
	}

	return linux, windows, macos, nil
}

func AllStats(c *fiber.Ctx) error {
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

func AllStatsTotal(c *fiber.Ctx) error {
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
