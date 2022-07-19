package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.etcd.io/bbolt"
)

func Stats(c *fiber.Ctx) error {
	user := c.Params("user")
	repo := c.Params("repo")
	key := user + "/" + repo
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("installations"))
		if b == nil {
			return nil
		}
		v := b.Get([]byte(key))
		if v == nil {
			return nil
		}
		return c.SendString(string(v))
	})
}

func AllStats(c *fiber.Ctx) error {
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("installations"))
		if b == nil {
			return nil
		}
		var lines []string
		err := b.ForEach(func(k, v []byte) error {
			lines = append(lines, string(k)+": "+string(v))
			return nil
		})
		if err != nil {
			return err
		}
		return c.SendString(strings.Join(lines, "\n"))
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
			"label":         "handeld installations",
			"message":       strconv.Itoa(total),
			"color":         "orange",
		})
	})
}
