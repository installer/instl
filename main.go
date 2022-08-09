package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"go.uber.org/zap"

	"github.com/installer/instl/internal/pkg/config"
	"github.com/installer/instl/internal/pkg/platforms"
	"github.com/installer/instl/scripts"

	"github.com/installer/instl/internal/pkg/handlers"
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
		pterm.Fatal.PrintOnError(err)
		script, err := scripts.ParseTemplateForPlatform(platform, config.Config{
			Owner:     *owner,
			Repo:      *repo,
			Version:   "latest",
			CreatedAt: time.Now(),
			Verbose:   *verbose,
		})
		pterm.Fatal.PrintOnError(err)
		fmt.Println(script)
		os.Exit(0)
	}

	pterm.DefaultBigText.WithLetters(putils.LettersFromString("  INSTL")).Render()

	logger, _ := zap.NewProduction(zap.WithCaller(false))
	defer logger.Sync()

	engine := html.New("./html", ".html")
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	app.Get("/", handlers.RedirectToDocs)
	app.Get("/docs", handlers.RedirectToDocs)
	app.Get("/documentation", handlers.RedirectToDocs)

	// API
	app.Get("/api/stats", handlers.AllStatsAPI)
	app.Get("/api/stats/:user/:repo", handlers.RepoStatsAPI)
	app.Get("/api/stats/total", handlers.AllStatsTotalAPI)

	// Stats
	// - Stats landing page
	app.Static("/stats", "./html/stats.html")
	app.Get("/stats/:any", func(ctx *fiber.Ctx) error { return ctx.SendStatus(404) })
	// - Total stats
	app.Get("/stats/total/badge/shields.io", handlers.AllStatsTotalBadge)
	// - Repo stats
	app.Get("/stats/:user/:repo", handlers.RepoStats)
	app.Get("/stats/:user/:repo/badge/shields.io", handlers.RepoStatsBadge)

	// Installation script generator
	app.Get("/:user/:repo", handlers.MissingPlatform)
	app.Get("/:user/:repo/:os", handlers.Installation)
	app.Get("/:user/:repo/:os/verbose", handlers.InstallationVerbose)

	pterm.Fatal.PrintOnError(app.Listen(":80"))
}
