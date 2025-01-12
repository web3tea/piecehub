package main

import (
	"fmt"

	"github.com/strahe/piecehub/config"
	"github.com/urfave/cli/v2"
)

var dirCmd = &cli.Command{
	Name:      "dir",
	Usage:     "start a piecehub server with disk storage",
	ArgsUsage: "[path1] [path2] ...",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "direct-io",
			Aliases: []string{"d"},
			Usage:   "enable direct I/O for disk operations",
			Value:   true,
			EnvVars: []string{"PIECEHUB_DIRECT_IO"},
		},
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
				Name:     path,
				RootDir:  path,
				DirectIO: c.Bool("direct-io"),
			})
		}
		return runServer(&cfg)
	},
}
