package cmd

import (
	"fmt"
	"github.com/hedzr/consul-tags/objects/consul"
	"gopkg.in/urfave/cli.v2"
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
	{
		Name:    "logger",
		Aliases: []string{"l"},
		Usage:   "logger helpers",
		// SubFlags: []cli.Flag { {},{},...   },
		//Flags:       consul.Flags,
		Subcommands: consul.LoggerCommands,
		//Before:      consul.Before,
		Action: func(c *cli.Context) error {
			fmt.Println("logger task: ", c.Args().First())
			return nil
		},
	},
}
