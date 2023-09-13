package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.etcd.io/bbolt"
	"strings"
)

func Sitemap(c *fiber.Ctx) error {
	var pages []string

	pages = append(pages, "https://instl.sh/")
	pages = append(pages, "https://instl.sh/stats")
	pages = append(pages, "https://instl.sh/privacy")
	pages = append(pages, "https://instl.sh/terms")

	var statsMap = map[string]struct{}{}
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("installations"))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var repo string
			parts := strings.Split(string(k), "/")
			repo = strings.Join(parts[:2], "/")
			statsMap[repo] = struct{}{}
			return nil
		})
	})
	if err != nil {
		return err
	}

	for repo := range statsMap {
		pages = append(pages, fmt.Sprintf("https://instl.sh/%s", repo))
	}

	c.Set("Content-Type", "application/xml")
	return c.SendString(generateSitemap(pages))
}

func generateSitemap(urls []string) string {
	var sitemap strings.Builder

	sitemap.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sitemap.WriteString("\n")
	sitemap.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	sitemap.WriteString("\n")
	for _, url := range urls {
		sitemap.WriteString("  <url>")
		sitemap.WriteString("\n")
		sitemap.WriteString(fmt.Sprintf("    <loc>%s</loc>", url))
		sitemap.WriteString("\n")
		sitemap.WriteString("  </url>")
		sitemap.WriteString("\n")
	}
	sitemap.WriteString(`</urlset>`)

	return sitemap.String()
}
