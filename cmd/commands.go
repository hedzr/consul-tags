package cmd

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"hedzr.com/consul-tags/objects/consul"
)

var Commands = []*cli.Command{
	{
		Name:    "consul",
		Aliases: []string{"c"},
		Usage:   "consul helpers",
		// SubFlags: []cli.Flag { {},{},...   },
		Flags:       consul.Flags,
		Subcommands: consul.Commands,
		Before:      consul.Before,
		Action: func(c *cli.Context) error {
			fmt.Println("consul task: ", c.Args().First())
			return nil
		},
	},
}
