package main

import (
	"context"

	"git.andinfinity.de/gitea-release-drafter/src"
	"git.andinfinity.de/gitea-release-drafter/src/config"
	githubactions "github.com/sethvargo/go-githubactions"
)

func run() error {
	ctx := context.Background()
	ghaction := githubactions.New()

	cfg, err := config.NewFromInputs(ghaction)
	if err != nil {
		return err
	}

	action, err := src.NewAction(&ctx, cfg)
	if err != nil {
		return err
	}

	return action.Run()
}

func main() {
	err := run()
	if err != nil {
		githubactions.Fatalf("%v", err)
	}
}
