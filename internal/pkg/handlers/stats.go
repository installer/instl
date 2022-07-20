package handlers

import (
	"strconv"

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
