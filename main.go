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
	test := flag.Bool("test", false, "test")
	flag.Parse()
	// Check if test flag is set
	if *test {
		platform, err := platforms.Parse(runtime.GOOS)
		pterm.Fatal.PrintOnError(err)
		script, err := scripts.ParseTemplateForPlatform(platform, config.Config{
			Owner:     "installer",
			Repo:      "test-repo",
			Version:   "latest",
			CreatedAt: time.Now(),
			Verbose:   true,
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
	app.Get("/:user/:repo/:os", handlers.Installation)

	pterm.Fatal.PrintOnError(app.Listen(":80"))
}
