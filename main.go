package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/installer/instl/internal/pkg/global"
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
			Owner:        *owner,
			Repo:         *repo,
			Version:      "latest",
			CreatedAt:    time.Now(),
			Verbose:      *verbose,
			InstlVersion: global.Version,
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

	// Legal
	app.Get("/privacy", handlers.PrivacyPolicy)
	app.Get("/terms", handlers.TermsOfService)
	app.Get("/tos", handlers.Redirect("/terms"))

	// Sitemap
	app.Get("/sitemap.xml", handlers.Sitemap)

	// API
	apiV1 := app.Group("/api/v1")
	apiV1.Get("/stats", handlers.AllStatsAPI)
	apiV1.Get("/stats/:user/:repo", handlers.RepoStatsAPI)
	apiV1.Get("/stats/total", handlers.AllStatsTotalAPI)
	// - Badge
	badgeAPIV1 := apiV1.Group("/badge")
	shieldsIOAPIV1 := badgeAPIV1.Group("/shields.io")
	shieldsIOAPIV1.Get("/stats/:user/:repo", handlers.RepoStatsShieldsIOBadge)
	shieldsIOAPIV1.Get("/stats/total", handlers.AllStatsTotalBadge)

	// Stats
	app.Get("/stats", handlers.AllStatsPage)
	app.Get("/stats/:any", func(ctx *fiber.Ctx) error { return ctx.SendStatus(404) })
	app.Get("/stats/:user/:repo", handlers.RedirectLegacyStats) // Legacy redirect
	app.Get("/:user/:repo", handlers.RepoStatsPage)
	app.Get("/:user/:repo/maintainer", handlers.RepoMaintainerPage)

	// Installation script generator
	app.Get("/:user/:repo/:os", handlers.Installation)
	app.Get("/:user/:repo/:os/verbose", handlers.InstallationVerbose)

	log.Fatal(app.Listen(":80"))
}
