/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/consul-tags"
	"github.com/hedzr/consul-tags/objects/consul"
)

var (
	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "consul-tags",
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
				kvCommands,
			},
		},

		AppName:    "consul-tags",
		Version:    consul_tags.Version,
		VersionInt: consul_tags.VersionInt,
		Copyright:  "consul-tags is an effective devops tool",
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

	kvCommands = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:        "kvstore",
			Full:        "kv",
			Aliases:     []string{"kvstore",},
			Description: "consul kv store operations...",
			Flags:       *cmdr.Clone(&consulConnectFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag),
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "b",
					Full:        "backup",
					Aliases:     []string{"bk", "bf", "bkp",},
					Description: "Dump Consul's KV database to a JSON/YAML file",
					Action:      kvBackup,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "o",
								Full:                    "output",
								Description:             "Write output to a file (*.json / *.yml)",
								DefaultValuePlaceholder: "FILE",
							},
							DefaultValue: "consul-backup.json",
						},
					},
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "r",
					Full:        "restore",
					Description: "restore to Consul's KV store, from a a JSON/YAML backup file",
					Action:      kvRestore,
					Flags: []*cmdr.Flag{
						{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "i",
								Full:                    "input",
								Description:             "read the input file (*.json / *.yml)",
								DefaultValuePlaceholder: "FILE",
							},
							DefaultValue: "consul-backup.json",
						},
					},
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
					Action:      msList,
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
			Flags:       *cmdr.Clone(&consulConnectFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag),
		},
		SubCommands: []*cmdr.Command{
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "ls",
					Full:        "list",
					Aliases:     []string{"l", "lst", "dir",},
					Description: "list tags.",
					Action:      msTagsList,
				},
			},
			// {
			// 	BaseOpt: cmdr.BaseOpt{
			// 		Short:       "a",
			// 		Full:        "add",
			// 		Aliases:     []string{"create", "new",},
			// 		Description: "add tags.",
			// 		Action:      msTagsAdd,
			// 		Flags: []*cmdr.Flag{
			// 			{
			// 				BaseOpt: cmdr.BaseOpt{
			// 					Short:                   "ls",
			// 					Full:                    "list",
			// 					Aliases:                 []string{"l", "lst", "dir",},
			// 					Description:             "a comma list to be added",
			// 					DefaultValuePlaceholder: "LIST",
			// 				},
			// 				DefaultValue: []string{},
			// 			},
			// 		},
			// 	},
			// },
			// {
			// 	BaseOpt: cmdr.BaseOpt{
			// 		Short:       "r",
			// 		Full:        "rm",
			// 		Aliases:     []string{"remove", "erase", "delete", "del",},
			// 		Description: "remove tags.",
			// 		Action:      msTagsRemove,
			// 		Flags: []*cmdr.Flag{
			// 			{
			// 				BaseOpt: cmdr.BaseOpt{
			// 					Short:                   "ls",
			// 					Full:                    "list",
			// 					Aliases:                 []string{"l", "lst", "dir",},
			// 					Description:             "a comma list to be added.",
			// 					DefaultValuePlaceholder: "LIST",
			// 				},
			// 				DefaultValue: []string{},
			// 			},
			// 		},
			// 	},
			// },
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "m",
					Full:        "modify",
					Aliases:     []string{"mod", "update", "change",},
					Description: "modify tags.",
					Action:      msTagsModify,
					Flags: append(*cmdr.Clone(&modifyFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag),
						&cmdr.Flag{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "a",
								Full:                    "add",
								Aliases:                 []string{"add-list"},
								Description:             "a comma list to be added.",
								DefaultValuePlaceholder: "LIST",
								Group:                   "List",
							},
							DefaultValue: []string{},
						},
						&cmdr.Flag{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "r",
								Full:                    "rm",
								Aliases:                 []string{"remove", "rm-list"},
								Description:             "a comma list to be removed.",
								DefaultValuePlaceholder: "LIST",
								Group:                   "List",
							},
							DefaultValue: []string{},
						}),
				},
			},
			{
				BaseOpt: cmdr.BaseOpt{
					Short:       "t",
					Full:        "toggle",
					Aliases:     []string{"tog", "switch",},
					Description: "toggle tags for ms.",
					Action:      msTagsToggle,
					Flags: append(*cmdr.Clone(&modifyFlags, &[]*cmdr.Flag{}).(*[]*cmdr.Flag),
						&cmdr.Flag{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "x",
								Full:                    "address",
								Description:             "the address of the service (by id or name)",
								DefaultValuePlaceholder: "HOST:PORT",
							},
							DefaultValue: "",
						},
						&cmdr.Flag{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "s",
								Full:                    "set",
								Description:             "set to `tag` which service specified by --address",
								DefaultValuePlaceholder: "LIST",
								Group:                   "List",
							},
							DefaultValue: []string{},
						},
						&cmdr.Flag{
							BaseOpt: cmdr.BaseOpt{
								Short:                   "u",
								Full:                    "unset",
								Aliases:                 []string{"reset",},
								Description:             "and reset the others service nodes to `tag`",
								DefaultValuePlaceholder: "LIST",
								Group:                   "List",
							},
							DefaultValue: []string{},
						},
					),
				},
			},
		},
	}

	modifyFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "d",
				Full:        "delim",
				Description: "delimitor char in `non-plain` mode.",
			},
			DefaultValue: "=",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "c",
				Full:        "clear",
				Description: "clear all tags.",
				Group:       "Operate",
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "g",
				Full:        "string",
				Aliases:     []string{"string-mode"},
				Description: "In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.",
				Group:       "Mode",
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "m",
				Full:        "meta",
				Aliases:     []string{"meta-mode"},
				Description: "In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false)",
				Group:       "Mode",
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "2",
				Full:        "both",
				Aliases:     []string{"both-mode"},
				Description: "In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.",
				Group:       "Mode",
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "p",
				Full:        "plain",
				Aliases:     []string{"plain-mode"},
				Description: "In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.",
				Group:       "Mode",
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "t",
				Full:        "tag",
				Aliases:     []string{"tag-mode"},
				Description: "In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.",
				Group:       "Mode",
			},
			DefaultValue: true,
		},
	}

	consulConnectFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "a",
				Full:                    "addr",
				Description:             "Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')",
				DefaultValuePlaceholder: "HOST[:PORT]",
				Group:                   "Consul",
			},
			DefaultValue: consul.DEFAULT_CONSUL_HOST,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "p",
				Full:                    "port",
				Description:             "Consul port",
				DefaultValuePlaceholder: "PORT",
				Group:                   "Consul",
			},
			DefaultValue: consul.DEFAULT_CONSUL_PORT,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "K",
				Full:                    "insecure",
				Description:             "Skip TLS host verification",
				DefaultValuePlaceholder: "PORT",
				Group:                   "Consul",
			},
			DefaultValue: true,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "",
				Full:                    "prefix",
				Description:             "Root key prefix",
				DefaultValuePlaceholder: "ROOT",
				Group:                   "Consul",
			},
			DefaultValue: "/",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "",
				Full:                    "cacert",
				Description:             "Client CA cert",
				DefaultValuePlaceholder: "FILE",
				Group:                   "Consul",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "",
				Full:                    "cert",
				Description:             "Client cert",
				DefaultValuePlaceholder: "FILE",
				Group:                   "Consul",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "",
				Full:                    "scheme",
				Description:             "Consul connection scheme (HTTP or HTTPS)",
				DefaultValuePlaceholder: "SCHEME",
				Group:                   "Consul",
			},
			DefaultValue: consul.DEFAULT_CONSUL_SCHEME,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "u",
				Full:                    "username",
				Description:             "HTTP Basic auth user",
				DefaultValuePlaceholder: "USERNAME",
				Group:                   "Consul",
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
				Group:                   "Consul",
			},
			DefaultValue: "",
		},
	}
)
