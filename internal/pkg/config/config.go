package config

import "time"

type Config struct {
	Owner        string
	Repo         string
	Version      string
	InstlVersion string
	CreatedAt    time.Time
	Verbose      bool
}

// Combine merges two configs together.
// If a field is set in both configs, the value from the second config will be used.
func Combine(configs ...Config) Config {
	var cfg Config
	for _, c := range configs {
		if c.Owner != "" {
			cfg.Owner = c.Owner
		}
		if c.Repo != "" {
			cfg.Repo = c.Repo
		}
		if c.Version != "" {
			cfg.Version = c.Version
		}
		if c.CreatedAt != (time.Time{}) {
			cfg.CreatedAt = c.CreatedAt
		}
		if c.Verbose {
			cfg.Verbose = true
		}
		if c.InstlVersion != "" {
			cfg.InstlVersion = c.InstlVersion
		} else {
			cfg.InstlVersion = "dev"
		}
	}
	return cfg
}
