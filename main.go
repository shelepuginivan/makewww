package main

import (
	"log"

	"github.com/shelepuginivan/makewww/pkg/builder"
	"github.com/shelepuginivan/makewww/pkg/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	builder := builder.New(cfg)
	if err := builder.Build(); err != nil {
		log.Fatal(err)
	}
}
