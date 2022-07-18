package config

import "time"

type Config struct {
	Owner     string
	Repo      string
	Version   string
	CreatedAt time.Time
}
