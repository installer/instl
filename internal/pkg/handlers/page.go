package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.etcd.io/bbolt"
	"slices"
	"strconv"
	"strings"
)

func HomePage(c *fiber.Ctx) error {
	return c.Render("home.gohtml", nil)
}

func RepoStatsPage(c *fiber.Ctx) error {
	owner, repo := getOwnerAndRepo(c)
	linux, windows, macos, err := getInstallationCountPerPlatform(owner, repo)
	if err != nil {
		return err
	}
	return c.Render("repo.gohtml", map[string]any{
		"Windows": windows,
		"Linux":   linux,
		"MacOS":   macos,
		"Total":   linux + windows + macos,
		"Owner":   owner,
		"Repo":    repo,
	})
}

type Stat struct {
	Owner string
	Repo  string
	Count int
}

func AllStatsPage(c *fiber.Ctx) error {
	statsMap := map[string]int{}

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("installations"))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			count, err := strconv.Atoi(string(v))
			if err != nil {
				return err
			}

			var repo string
			parts := strings.Split(string(k), "/")
			repo = strings.Join(parts[:2], "/")
			statsMap[repo] += count
			return nil
		})
	})
	if err != nil {
		return err
	}

	var stats []Stat
	var total int

	for k, v := range statsMap {
		parts := strings.Split(k, "/")
		stats = append(stats, Stat{
			Owner: parts[0],
			Repo:  parts[1],
			Count: v,
		})
		total += v
	}

	slices.SortFunc(stats, func(a, b Stat) int {
		return b.Count - a.Count
	})

	return c.Render("stats.gohtml", map[string]any{
		"Stats": stats,
		"Total": total,
	})
}
