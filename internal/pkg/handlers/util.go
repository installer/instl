package handlers

import (
	"go.etcd.io/bbolt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func getOwnerAndRepo(c *fiber.Ctx) (owner, repo string) {
	owner = c.Params("user")
	owner = strings.ToLower(owner)
	repo = c.Params("repo")
	repo = strings.ToLower(repo)
	return owner, repo
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
