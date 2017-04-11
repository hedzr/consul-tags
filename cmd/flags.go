package cmd

import (
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

var GingerCrouton = false

var tasks = []string{"cook", "clean", "laundry", "eat", "sleep", "code"}
var language string

// https://github.com/urfave/cli

var Flags = []cli.Flag{
	altsrc.NewIntFlag(&cli.IntFlag{Name: "test"}),

	&cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   ".ops-config.yml",
		Usage:   "Load configuration from `FILE`",
	},

	&cli.StringFlag{
		Name:        "lang",
		Aliases:     []string{"l"},
		Value:       "english",
		Usage:       "language for the greeting",
		Destination: &language,
		//EnvVar: "APP_LANG", // value from the environment
		EnvVars: []string{"LEGACY_COMPAT_LANG", "APP_LANG", "LANG"},
	},

	//&cli.StringFlag{
	//	Name:  "addr, a",
	//	Value: "consul.ops.local:8500",
	//	Usage: "consul center host `HOST:IP`",
	//	Destination: &consul_addr,
	//},

	&cli.BoolFlag{
		Name:  "4",
		Usage: "IPV4 prior",
	},

	&cli.BoolFlag{
		Name:  "6",
		Usage: "IPV6 prior",
	},

	&cli.BoolFlag{
		Name:    "V",
		Aliases: []string{"verbose"},
		Usage:   "verbose debug info?",
	},

	&cli.BoolFlag{
		Name:  "vv",
		Usage: "verbose debug info?",
	},

	&cli.BoolFlag{
		Name:  "vvv",
		Usage: "verbose debug info?",
	},

	//&cli.BoolFlag{
	//	Name:  "ginger-crouton",
	//	Destination: &GingerCrouton,
	//	Usage: "is it in the soup?",
	//},
}

func AfterFunc(c *cli.Context) error {
	if c.Bool("vv") {
		c.Set("V", "true")
	}

	if c.Bool("vvv") {
		c.Set("vv", "true")
		c.Set("V", "true")
	}

	if c.Bool("6") {
		c.Set("4", "false")
	} else {
		c.Set("4", "true")
	}

	return nil
}
