package main

import (
	"fmt"

	"github.com/strahe/piecehub/config"
	"github.com/urfave/cli/v2"
)

var s3Cmd = &cli.Command{
	Name:      "s3",
	Usage:     "start a piecehub server with s3 storage",
	ArgsUsage: "[bucket1] [bucket2] ...",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "endpoint",
			Usage:    "s3 endpoint",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "ak",
			Usage: "s3 access key",
		},
		&cli.StringFlag{
			Name:  "sk",
			Usage: "s3 secret key",
		},
		&cli.StringFlag{
			Name:  "region",
			Usage: "s3 region",
		},
		&cli.BoolFlag{
			Name:  "ssl",
			Usage: "use ssl for s3",
		},
		&cli.StringFlag{
			Name:  "listen",
			Usage: "server listen address",
		},
	},
	Action: func(c *cli.Context) error {
		buckets := c.Args().Slice()
		if len(buckets) == 0 {
			return fmt.Errorf("no buckets provided")
		}
		cfg := config.DefaultConfig
		if c.IsSet("listen") {
			cfg.Server.Address = c.String("listen")
		}
		for _, bucket := range buckets {
			cfg.S3s = append(cfg.S3s, config.S3Config{
				Name:      bucket,
				Endpoint:  c.String("endpoint"),
				Region:    c.String("region"),
				Bucket:    bucket,
				AccessKey: c.String("ak"),
				SecretKey: c.String("sk"),
				UseSSL:    c.Bool("ssl"),
			})
		}
		return runServer(&cfg)
	},
}
