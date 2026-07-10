package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Dir    string
	Output string
}

func Parse() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get cwd: %w", err)
	}

	var cfg Config

	flag.StringVar(&cfg.Dir, "dir", cwd, "directory to build the website from")
	flag.StringVar(&cfg.Output, "output", filepath.Join(cwd, "dist"), "output directory for the website")
	flag.Parse()

	return &cfg, nil
}
