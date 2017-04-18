package main

import (
	"os"

	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"

	log "github.com/cihub/seelog"

	"github.com/hedzr/consul-tags/cmd"
	"github.com/hedzr/consul-tags/objects/consul"
)

const (
	DEFAULT_DEBUG = false

	Ln2       = 0.693147180559945309417232121458176568075500134360255254120680009
	Log2E     = 1 / Ln2 // this is a precise reciprocal
	Billion   = 1e9     // float constant
	hardEight = (1 << 100) >> 97

	Sunday                     = iota
	Monday, Tuesday, Wednesday = 1, 2, 3
	Thursday, Friday, Saturday = 4, 5, 6
)

//var consul_addr = ""

func init() {
	consul.InitLogger()
}

func main() {

	// setup the flags/options within a .ops-config.yml
	var yamlConfig = cmd.NewYamlSourceFromFlagFunc("config")

	app := &cli.App{
		Name:    "consul-tags",
		Usage:   "The DevOps Toolset!",
		Version: "0.3.2",
		Authors: []*cli.Author{
			{
				Name: "Hedzr", Email: "hedzrz@gmail.com",
			},
		},
		Copyright: "Copyright (c) by Hedzr Studio, 2017. All Rights Reserved.",

		EnableShellCompletion: true,
		ShellComplete: func(c *cli.Context) {
			cli.DefaultAppComplete(c)
		},

		Flags:    consul.Flags, // flags 无法分组，当前。
		Commands: consul.Commands,

		Before: func(c *cli.Context) error {
			err := altsrc.InitInputSourceWithContext(cmd.Flags, yamlConfig)(c)
			if err != nil {
				log.Criticalf("Error: %v", err)
				return err
			}
			cmd.LoadConfigFile(c)
			return consul.Before(c)
		},
		After:  cmd.AfterFunc,
		Action: cmd.DefaultAction,
	}

	//sort.Sort(cli.FlagsByName(app.Flags))
	//sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
