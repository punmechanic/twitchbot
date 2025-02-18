package twitchbot

import (
	"context"
	"fmt"
	"os"

	"example.com/twitchbot/pkg/twitch"
	"github.com/urfave/cli/v3"
)

var Root = cli.Command{
	Commands: []*cli.Command{
		{
			Name: "serve",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name: "broadcaster-ids",
				},
			},
			Action: serve,
		},
		{
			Name: "lookup-broadcaster-id",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "broadcaster_id",
					Min:  0,
					// Setting Max to -1 seems to break argument parsing in urfave/cli, and urfave will discard
					// any arguments passed. Commenting out this line fixes that problem but shows us a warning.
					// Line is commented out for now so we can work
					// Max: -1,
				},
			},
			Action: lookupBroadcasterByLoginName,
		},
	},
}

func Run(ctx context.Context) error {
	return Root.Run(ctx, os.Args)
}

func lookupBroadcasterByLoginName(ctx context.Context, c *cli.Command) error {
	cfg := initTwitchConfig([]string{"user:read:chat"})
	token, _, err := fetchTokenWithFallback(ctx, cfg)
	defer saveTokenInKeyring(token)

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

	return nil
}
