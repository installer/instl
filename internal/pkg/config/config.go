package config

import "time"

type Config struct {
	Owner     string
	Repo      string
	Version   string
	CreatedAt time.Time
	Verbose   bool
}

func Combine(configs ...Config) Config {
	var config Config
	for _, c := range configs {
		config.Owner = c.Owner
		config.Repo = c.Repo
		config.Version = c.Version
		config.CreatedAt = c.CreatedAt
		config.Verbose = c.Verbose
	}
	return config
}
