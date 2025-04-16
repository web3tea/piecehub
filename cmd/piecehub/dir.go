package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/web3tea/piecehub/config"
)

var dirCmd = &cli.Command{
	Name:      "dir",
	Usage:     "start a piecehub server with disk storage",
	ArgsUsage: "[path1] [path2] ...",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "listen",
			Usage: "server listen address",
		},
	},
	Action: func(c *cli.Context) error {
		paths := c.Args().Slice()
		if len(paths) == 0 {
			return fmt.Errorf("no paths provided")
		}
		cfg := config.DefaultConfig
		if c.IsSet("listen") {
			cfg.Server.Address = c.String("listen")
		}
		for _, path := range paths {
			cfg.Disks = append(cfg.Disks, config.DiskConfig{
				Name:    path,
				RootDir: path,
			})
		}

		tokens := c.StringSlice("tokens")
		if len(tokens) > 0 {
			cfg.Server.Tokens = append(cfg.Server.Tokens, tokens...)
		}

		return runServer(&cfg)
	},
}
