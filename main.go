package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"go.uber.org/zap"

	"github.com/installer/installer/internal/pkg/config"
	"github.com/installer/installer/internal/pkg/platforms"
	"github.com/installer/installer/scripts"

	"github.com/installer/installer/internal/pkg/handlers"
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

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(recover.New())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	app.Get("/", handlers.RedirectToDocs)
	app.Get("/:user/:repo", handlers.MissingPlatform)
	app.Get("/:user/:repo/:os", handlers.Installation)
	app.Get("/:user/:repo/:os/verbose", handlers.InstallationVerbose)

	pterm.Fatal.PrintOnError(app.Listen(":80"))
}
