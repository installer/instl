package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/installer/instl/templates"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/installer/instl/internal/pkg/config"
	"github.com/installer/instl/internal/pkg/handlers"
	"github.com/installer/instl/internal/pkg/platforms"
	"github.com/installer/instl/scripts"
)

func main() {
	test := flag.Bool("test", false, "enable test mode; don't start server; print script to stdout")
	verbose := flag.Bool("verbose", false, "verbose output for test mode")
	owner := flag.String("owner", "installer", "repo owner for test mode")
	repo := flag.String("repo", "test-repo", "repo name for test mode")
	flag.Parse()

	// Check if test flag is set
	if *test {
		platform, err := platforms.Parse(runtime.GOOS)
		if err != nil {
			log.Fatal(err)
		}

		script, err := scripts.ParseTemplateForPlatform(platform, config.Config{
			Owner:     *owner,
			Repo:      *repo,
			Version:   "latest",
			CreatedAt: time.Now(),
			Verbose:   *verbose,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(script)
		os.Exit(0)
	}

	fmt.Println(`        ██ ███    ██ ███████ ████████ ██      
        ██ ████   ██ ██         ██    ██      
        ██ ██ ██  ██ ███████    ██    ██      
        ██ ██  ██ ██      ██    ██    ██      
        ██ ██   ████ ███████    ██    ███████ 
`)

	engine := templates.New()
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/", handlers.HomePage)
	app.Static("/", "./static")
	app.Get("/docs", handlers.RedirectToDocs)
	app.Get("/documentation", handlers.RedirectToDocs)

	// Sitemap
	app.Get("/sitemap.xml", handlers.Sitemap)

	// API
	app.Get("/api/stats", handlers.AllStatsAPI)
	app.Get("/api/stats/:user/:repo", handlers.RepoStatsAPI)
	app.Get("/api/stats/total", handlers.AllStatsTotalAPI)

	// Stats
	// - Stats landing page
	app.Get("/stats", handlers.AllStatsPage)
	app.Get("/stats/:any", func(ctx *fiber.Ctx) error { return ctx.SendStatus(404) })
	// - Total stats
	app.Get("/stats/total/badge/shields.io", handlers.AllStatsTotalBadge)
	// - Legacy redirect
	app.Get("/:user/:repo", handlers.RepoStatsPage)
	app.Get("/stats/:user/:repo", handlers.RedirectLegacyStats)
	// - Badge
	app.Get("/stats/:user/:repo/badge/shields.io", handlers.RepoStatsBadge)
	app.Get("/:user/:repo/badge/shields.io", handlers.RepoStatsBadge)

	// Installation script generator
	app.Get("/:user/:repo/:os", handlers.Installation)
	app.Get("/:user/:repo/:os/verbose", handlers.InstallationVerbose)

	log.Fatal(app.Listen(":80"))
}
