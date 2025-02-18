package twitchbot

import (
	"context"
	"fmt"
	"log"
	"os"

	"example.com/twitchbot/pkg/twitch"
	"github.com/urfave/cli/v3"
)

var Root = cli.Command{
	Commands: []*cli.Command{
		{
			Name: "serve",
			Action: func(ctx context.Context, _ *cli.Command) error {
				return serve(ctx)
			},
		},
		{
			Name: "util",
			Commands: []*cli.Command{
				{
					Name: "lookup-broadcaster-id",
					Arguments: []cli.Argument{
						&cli.StringArg{
							Name: "broadcaster_id",
							Max:  -1,
						},
					},
					Action: lookupBroadcasterByLoginName,
				},
			},
		},
	},
}

func Run(ctx context.Context) error {
	return Root.Run(ctx, os.Args)
}

func lookupBroadcasterByLoginName(ctx context.Context, c *cli.Command) error {
	cfg := initTwitchConfig([]string{"user:read:chat"})
	token, err := fetchTokenFromKeyring()
	if err != nil {
		token, err = fetchTokenFromTwitch(ctx, cfg)
		if err != nil {
			return err
		}
	}

	client := twitch.New(cfg, token)
	users, err := client.Users(ctx, &twitch.UsersRequest{
		Login: c.Args().Slice(),
	})

	if err != nil {
		return fmt.Errorf("fetch users: %w", err)
	}

	for _, user := range users.Data {
		fmt.Fprintf(c.Writer, "%s\n", user.ID)
	}

	if err := saveTokenInKeyring(token); err != nil {
		// Print a warning but don't die
		log.Printf("could not save token to OS keyring: %s", err)
	}
	return nil
}
