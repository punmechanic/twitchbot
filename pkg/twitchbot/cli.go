package twitchbot

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

var Root = cli.Command{
	Commands: []*cli.Command{
		{
			Name: "local",
			Action: func(ctx context.Context, _ *cli.Command) error {
				return runLocal(ctx)
			},
		},
	},
}

func Run(ctx context.Context) error {
	return Root.Run(ctx, os.Args)
}
