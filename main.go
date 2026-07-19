package main

import (
	"log"

	"github.com/shelepuginivan/makewww/pkg/builder"
	"github.com/shelepuginivan/makewww/pkg/config"
	"github.com/shelepuginivan/makewww/pkg/output"
	"github.com/shelepuginivan/makewww/pkg/source"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	src, err := source.NewFS(cfg.Dir)
	if err != nil {
		log.Fatal(err)
	}

	out, err := output.NewFS(cfg.Output)
	if err != nil {
		log.Fatal(err)
	}

	builder, err := builder.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = builder.Build(src, out)
	if err != nil {
		log.Fatal(err)
	}
}
