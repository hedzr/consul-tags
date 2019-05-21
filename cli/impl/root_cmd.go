/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/consul-tags"
)

var (
	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "austr",
				Flags: []*cmdr.Flag{
					// global options here.
				},
			},
			PreAction:  Pre,
			PostAction: Post,
			SubCommands: []*cmdr.Command{
				// dnsCommands,
				// playCommand,
				// generatorCommands,
				// serverCommands,
				msCommands,
			},
		},

		AppName:    "consul-tags",
		Version:    consul_tags.Version,
		VersionInt: consul_tags.VersionInt,
		Copyright:  "austr is an effective devops tool",
		Author:     "Hedzr Yeh <hedzrz@gmail.com>",
	}

	serverCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Name:        "server",
			Short:       "s",
			Full:        "server",
			Aliases:     []string{"serve", "svr",},
			Description: "server ops: for linux service/daemon.",
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "s",
					Full:        "start",
					Aliases:     []string{"run", "startup",},
					Description: "startup this system service/daemon.",
					// Action:impl.ServerStart,
				},
				// PreAction: impl.ServerStartPre,
				// PostAction: impl.ServerStartPost,
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "stop",
					Aliases:     []string{"stp", "halt", "pause",},
					Description: "stop this system service/daemon.",
					// Action:impl.ServerStop,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "restart",
					Aliases:     []string{"reload",},
					Description: "restart this system service/daemon.",
					// Action:impl.ServerRestart,
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Full:        "status",
					Aliases:     []string{"st",},
					Description: "display its running status as a system service/daemon.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "i",
					Full:        "install",
					Aliases:     []string{"setup",},
					Description: "install as a system service/daemon.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "u",
					Full:        "uninstall",
					Aliases:     []string{"remove",},
					Description: "remove from a system service/daemon.",
				},
			},
		},
	}

	msCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:        "microservices",
			Full:        "ms",
			Aliases:     []string{"microservice", "micro-service",},
			Description: "micro-service operations...",
			Flags: []*cmdr.Flag{
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "n",
						Full:                    "name",
						Description:             "name of the service",
						DefaultValuePlaceholder: "NAME",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "i",
						Full:                    "id",
						Description:             "unique id of the service",
						DefaultValuePlaceholder: "ID",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:       "a",
						Full:        "all",
						Description: "all services",
					},
					DefaultValue: false,
				},
			},
		},
		SubCommands: []*cmdr.Command{
			tagsCommands,
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "l",
					Full:        "list",
					Aliases:     []string{"ls", "lst", "dir",},
					Description: "list services.",
				},
			},
		},
	}

	tagsCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			// Short:       "t",
			Full:        "tags",
			Aliases:     []string{},
			Description: "tags op.",
			Flags: []*cmdr.Flag{
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "a",
						Full:                    "addr",
						Description:             "Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')",
						DefaultValuePlaceholder: "HOST",
					},
					DefaultValue: "consul.ops.local",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "p",
						Full:                    "port",
						Description:             "Consul port",
						DefaultValuePlaceholder: "PORT",
					},
					DefaultValue: 8500,
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "K",
						Full:                    "insecure",
						Description:             "Skip TLS host verification",
						DefaultValuePlaceholder: "PORT",
					},
					DefaultValue: true,
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "",
						Full:                    "prefix",
						Description:             "Root key prefix",
						DefaultValuePlaceholder: "ROOT",
					},
					DefaultValue: "/",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "",
						Full:                    "cacert",
						Description:             "Client CA cert",
						DefaultValuePlaceholder: "FILE",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "",
						Full:                    "cert",
						Description:             "Client cert",
						DefaultValuePlaceholder: "FILE",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "",
						Full:                    "scheme",
						Description:             "Consul connection scheme (HTTP or HTTPS)",
						DefaultValuePlaceholder: "SCHEME",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "u",
						Full:                    "username",
						Description:             "HTTP Basic auth user",
						DefaultValuePlaceholder: "USERNAME",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "pw",
						Full:                    "password",
						Aliases:                 []string{"passwd", "pwd",},
						Description:             "HTTP Basic auth password",
						DefaultValuePlaceholder: "PASSWORD",
					},
					DefaultValue: "",
				},
				{
					BaseOpt: cmdr.BaseOpt{
						Short:                   "p",
						Full:                    "prefix",
						Description:             "Root key prefix",
						DefaultValuePlaceholder: "ROOT",
					},
					DefaultValue: "/",
				},
			},
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "ls",
					Full:        "list",
					Aliases:     []string{"l", "lst", "dir",},
					Description: "list tags.",
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "a",
					Full:        "add",
					Aliases:     []string{"create", "new",},
					Description: "add tags.",
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "ls",
								Full:                    "list",
								Aliases:                 []string{"l", "lst", "dir",},
								Description:             "a comma list to be added",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "rm",
					Aliases:     []string{"remove", "erase", "delete", "del",},
					Description: "remove tags.",
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "ls",
								Full:                    "list",
								Aliases:                 []string{"l", "lst", "dir",},
								Description:             "a comma list to be added.",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "m",
					Full:        "modify",
					Aliases:     []string{"mod", "update", "change",},
					Description: "modify tags.",
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "a",
								Full:                    "add",
								Description:             "a comma list to be added.",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "r",
								Full:                    "rm",
								Aliases:                 []string{"remove", "erase", "del",},
								Description:             "a comma list to be removed.",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "toggle",
					Aliases:     []string{"tog", "switch",},
					Description: "toggle tags for ms.",
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "s",
								Full:                    "set",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "u",
								Full:                    "unset",
								DefaultValuePlaceholder: "LIST",
							},
							DefaultValue: []string{},
						},
					},
				},
			},
		},
	}
)
