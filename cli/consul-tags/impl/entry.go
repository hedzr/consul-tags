/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"context"
	"io"
	"os"

	"github.com/hedzr/cmdr/v2/conf"

	"github.com/hedzr/consul-tags/objects/consul"

	cmdrv2 "github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	errorsv3 "gopkg.in/hedzr/errors.v3"
)

func PrepareApp(opts ...cli.Opt) (app cli.App) {
	app = cmdrv2.New(opts...).
		Info("consul-tags", "0.8.0").
		Copyright("An effective devops tool").
		Author("CLI Authors")
	// a cli tool to operate tags in consul service

	// TODO app.WithEnvPrefix("CT")

	// cmd.More1(app)

	app.Flg("dry-run", "n").
		Default(false).
		Description("run all but without committing").
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		Build() // no matter even if you're adding the duplicated one.

	// another way to disable `cmdr.WithForceDefaultAction(true)` is using
	// env-var FORCE_RUN=1 (builtin already).
	app.Flg("no-default").
		Description("disable force default action").
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			if b, ok := hitState.Value.(bool); ok {
				f.Set().Set("app.force-default-action", b) // disable/enable the final state about 'force default action'
			}
			return
		}).
		Build()

	// cmd.More2(app)

	app.Cmd("wrong").
		Description("a wrong command to return error for testing only").
		Hidden(true, true).
		// cmdline `FORCE_RUN=1 go run ./tiny wrong -d 8s` to verify this command to see the returned application error.
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			dur := cmd.Store().MustDuration("wrong.duration")
			println("the duration is:", dur.String(), conf.AppName, conf.Version)

			ec := errorsv3.New()
			defer ec.Defer(&err) // store the collected errors in native err and return it
			ec.Attach(io.ErrClosedPipe, errorsv3.New("something's wrong"), os.ErrPermission)
			// see the application error by running `go run ./tiny/tiny/main.go wrong`.
			return
		}).
		With(func(b cli.CommandBuilder) {
			b.Flg("duration", "d").
				Default("5s").
				Description("a duration var").
				Build()
		})

	// // redirect root commands into "server", so `app start|stop` -> `app server start|stop`.
	// app.WithRootCommand(func(root *cli.RootCommand) {
	// 	root.SetRedirectTo("server")
	// })

	// add 'server' command
	attachKVCommand(app.Cmd("kv", "kv", "kvstore", "kv"))

	// add more...
	attachMsCommand(app.Cmd("ms", "ms", "microservice", "ms"))
	return
}

func attachKVCommand(b cli.CommandBuilder) {
	b.Description("consul kv store operations...", ``).
		With(func(b cli.CommandBuilder) {
			attachConsulConnectFlagsV2(b)

			b.Cmd("backup", "b", "bk", "bf", "bkp").
				Description("Dump Consul's KV database to a JSON/YAML file", ``).
				OnAction(kvBackupV2).
				With(func(b cli.CommandBuilder) {
					b.Flg("output", "o").
						Default("consul-backup.json").
						PlaceHolder("FILE").
						Description("Write output to a file (*.json / *.yml)", ``).
						Build()
				})

			b.Cmd("restore", "r").
				Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
				OnAction(kvRestoreV2).
				With(func(b cli.CommandBuilder) {
					b.Flg("input", "i").
						Default("consul-backup.json").
						PlaceHolder("FILE").
						Description("Read the input file (*.json / *.yml)", ``).
						Build()
				})
		})
}

func attachMsCommand(b cli.CommandBuilder) {
	b.Description("consul micro-service operations...", ``).
		With(func(b cli.CommandBuilder) {
			b.Flg("name", "n").
				Default("").
				PlaceHolder("NAME").
				EnvVars("NAME").
				Description("name of the service", ``).
				Build()
			b.Flg("id", "i", "ID").
				Default("").
				PlaceHolder("ID").
				EnvVars("ID").
				Description("unique id of the service", ``).
				Build()
			b.Flg("all", "a").
				Default(false).
				Description("all services", ``).
				Build()

			attachConsulConnectFlagsV2(b)

			/* msListCmd := */
			b.Cmd("list", "ls", "l", "lst", "dir").
				Description("list services.", ``).
				OnAction(msListV2).
				Build()

			/* msTagsCmd := */
			b.Cmd("tags", "t").
				Description("tags operations of a services.", ``).
				With(func(b cli.CommandBuilder) {
					// ms tags list
					b.Cmd("list", "ls", "l", "lst", "dir").
						Description("list tags of a service.", ``).
						OnAction(msTagsListV2).
						Build()

					// ms tags modify
					b.Cmd("modify", "m", "mod", "modi", "update", "change").
						Description("modify tags of a service.", ``).
						OnAction(msTagsModifyV2).
						With(func(b cli.CommandBuilder) {
							attachModifyFlagsV2(b)

							b.Flg("add", "a", "add-list").
								Default([]string{}).
								PlaceHolder("LIST").
								Description("a comma list to be added.", ``).
								Group("List").
								Build()
							b.Flg("remove", "r", "rm-list", "rm", "del", "delete").
								Default([]string{}).
								PlaceHolder("LIST").
								Description("a comma list to be removed.", ``).
								Group("List").
								Build()
						})

					// ms tags toggle
					b.Cmd("toggle", "t", "tog", "switch").
						Description("toggle tags of a service.", ``).
						OnAction(msTagsToggleV2).
						With(func(b cli.CommandBuilder) {
							attachModifyFlagsV2(b)

							b.Flg("set", "s").
								Description("set to `tag` which service specified by --address", ``).
								Default([]string{}).
								PlaceHolder("LIST").
								Group("List").
								Build()
							b.Flg("unset", "u", "reset").
								Description("and reset the others service nodes to `tag`", ``).
								Default([]string{}).
								PlaceHolder("LIST").
								Group("List").
								Build()
							b.Flg("service-addr", "sa", "service-address").
								Default("").
								PlaceHolder("HOST:PORT").
								Description("address of the modifying service", ``).
								EnvVars("SERVICE_ADDR").
								Build()
						})
				})
		})
}

const ConsulGroup = "Consul"

func attachModifyFlagsV2(b cli.CommandBuilder) {
	b.Flg("delim", "d").
		Default("=").
		Description("delimitor char in `non-plain` mode.", ``).
		Build()

	const OperateGroup = "Operate"

	b.Flg("clear", "c").
		Default(false).
		Description("clear all tags.", ``).
		Group(OperateGroup).
		Build()

	const ModeGroup = "Mode"

	b.Flg("string", "g", "string-mode").
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.", ``).
		Default(false).
		ToggleGroup(ModeGroup).
		Build()

	b.Flg("meta", "m", "meta-mode").
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).", ``).
		Default(false).
		ToggleGroup(ModeGroup).
		Build()

	b.Flg("both", "2", "both-mode").
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.", ``).
		Default(false).
		ToggleGroup(ModeGroup).
		Build()

	b.Flg("plain", "p", "plain-mode").
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		Default(false).
		ToggleGroup(ModeGroup).
		Build()

	b.Flg("tag", "t", "tag-mode").
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		Default(true).
		ToggleGroup(ModeGroup).
		Build()
}

func attachConsulConnectFlagsV2(b cli.CommandBuilder) {
	b.Flg("addr", "a").
		Default(consul.DEFAULT_CONSUL_HOST).
		PlaceHolder("HOST:PORT").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		Group(ConsulGroup).
		EnvVars("ADDR").
		Build()
	b.Flg("port", "p").
		Default(consul.DEFAULT_CONSUL_PORT).
		PlaceHolder("PORT").
		Description("Consul port", ``).
		Group(ConsulGroup).
		EnvVars("PORT").
		Build()
	b.Flg("insecure", "K").
		Default(true).
		Description("Skip TLS host verification", ``).
		Group(ConsulGroup).
		EnvVars("PORT").
		Build()
	b.Flg("prefix", "px").
		Default("/").
		PlaceHolder("ROOT").
		Description("Root key prefix", ``).
		Group(ConsulGroup).
		EnvVars("PREFIX").
		Build()
	b.Flg("cacert").
		Default("").
		PlaceHolder("FILE").
		Description("Consul Client CA cert)", ``).
		Group(ConsulGroup).
		EnvVars("CACERT").
		Build()
	b.Flg("cert").
		Default("").
		PlaceHolder("FILE").
		Description("Consul Client cert", ``).
		Group(ConsulGroup).
		EnvVars("CERT").
		Build()
	b.Flg("scheme").
		Default(consul.DEFAULT_CONSUL_SCHEME).
		PlaceHolder("SCHEME").
		Description("Consul connection protocol", ``).
		Group(ConsulGroup).
		EnvVars("SCHEME").
		Build()
	b.Flg("username", "u", "user", "usr", "uid").
		Default("").
		PlaceHolder("USER").
		Description("HTTP Basic auth user", ``).
		Group(ConsulGroup).
		EnvVars("USER").
		Build()
	b.Flg("password", "pw", "passwd", "pass", "pwd").
		Default("").
		PlaceHolder("PASSWORD").
		Description("HTTP Basic auth password", ``).
		ExternalEditor(ExternalToolPasswordInput).
		Group(ConsulGroup).
		EnvVars("PASS").
		Build()
}

const ExternalToolPasswordInput = "PASSWD"

// func Entry() {
// 	if err := cmdr.Exec(buildRootCommand(),
// 		cmdr.WithEnvPrefix("CT"),
// 		cmdr.WithLogx(logrus.New("debug", false, true)),
// 		cmdr.WithWatchMainConfigFileToo(true),
// 	); err != nil {
// 		cmdr.Logger.Errorf("Error: %v", err)
// 	}
// }
//
// func buildRootCommand() (rootCmd *cmdr.RootCommand) {
// 	root := cmdr.Root(consul_tags.AppName, consul_tags.Version).
// 		Copyright("consul-tags is an effective devops tool", "Hedzr Yeh <hedzrz@gmail.com>").
// 		PreAction(Pre).PostAction(Post)
// 	rootCmd = root.RootCommand()
//
// 	// kv
//
// 	kvCmd := cmdr.NewSubCmd().
// 		Titles("kvstore", "kv").
// 		Description("consul kv store operations...", ``).
// 		AttachTo(root)
//
// 	attachConsulConnectFlags(kvCmd)
//
// 	kvBackupCmd := cmdr.NewSubCmd().
// 		Titles("backup", "b", "bk", "bf", "bkp").
// 		Description("Dump Consul's KV database to a JSON/YAML file", ``).
// 		Action(kvBackup).
// 		AttachTo(kvCmd)
// 	cmdr.NewString("consul-backup.json").Placeholder("FILE").
// 		Titles("output", "o").
// 		Description("Write output to a file (*.json / *.yml)", ``).
// 		AttachTo(kvBackupCmd)
//
// 	kvRestoreCmd := cmdr.NewSubCmd().
// 		Titles("restore", "r").
// 		Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
// 		Action(kvRestore).
// 		AttachTo(kvCmd)
// 	cmdr.NewString("consul-backup.json").Placeholder("FILE").
// 		Titles("input", "i").
// 		Description("Read the input file (*.json / *.yml)", ``).
// 		AttachTo(kvRestoreCmd)
//
// 	// ms
//
// 	msCmd := cmdr.NewSubCmd().
// 		Titles("microservice", "ms", "micro-service").
// 		Description("micro-service operations...", ``).
// 		AttachTo(root)
//
// 	cmdr.NewString().Placeholder("NAME").
// 		Titles("name", "n").
// 		Description("name of the service", ``).
// 		DefaultValue("", "NAME").
// 		AttachTo(msCmd)
// 	cmdr.NewString().Placeholder("ID").
// 		Titles("id", "i", "ID").
// 		Description("unique id of the service", ``).
// 		AttachTo(msCmd)
// 	cmdr.NewBool().
// 		Titles("all", "a").
// 		Description("all services", ``).
// 		AttachTo(msCmd)
//
// 	// ms list
//
// 	/* msListCmd := */
// 	cmdr.NewSubCmd().
// 		Titles("list", "ls", "l", "lst", "dir").
// 		Description("list services.", ``).
// 		Action(msList).
// 		AttachTo(msCmd)
//
// 	// ms tags
//
// 	msTagsCmd := cmdr.NewSubCmd().
// 		Titles("tags", "t").
// 		Description("tags operations of a services.", ``).
// 		Action(msList).
// 		AttachTo(msCmd)
//
// 	attachConsulConnectFlags(msTagsCmd)
//
// 	cmdr.NewSubCmd().
// 		Titles("list", "ls", "l", "lst", "dir").
// 		Description("list tags of a service.", ``).
// 		Action(msTagsList).
// 		AttachTo(msTagsCmd)
//
// 	// ms tags modify
//
// 	msTagsModifyCmd := cmdr.NewSubCmd().
// 		Titles("modify", "m", "mod", "modi", "update", "change").
// 		Description("modify tags of a service.", ``).
// 		Action(msTagsModify).
// 		AttachTo(msTagsCmd)
//
// 	attachModifyFlags(msTagsModifyCmd)
//
// 	cmdr.NewStringSlice().Placeholder("LIST").
// 		Titles("add", "a", "add-list").
// 		Description("a comma list to be added.", ``).
// 		Group("List").
// 		AttachTo(msTagsModifyCmd)
// 	cmdr.NewStringSlice().Placeholder("LIST").
// 		Titles("remove", "r", "rm-list", "rm", "del", "delete").
// 		Description("a comma list to be removed.", ``).
// 		Group("List").
// 		AttachTo(msTagsModifyCmd)
//
// 	// ms tags toggle
//
// 	msTagsToggleCmd := cmdr.NewSubCmd().
// 		Titles("toggle", "t", "tog", "switch").
// 		Description("toggle tags of a service.", ``).
// 		Action(msTagsToggle).
// 		AttachTo(msTagsCmd)
//
// 	attachModifyFlags(msTagsToggleCmd)
//
// 	cmdr.NewStringSlice().Placeholder("LIST").
// 		Titles("set", "s").
// 		Description("set to `tag` which service specified by --address", ``).
// 		Group("List").
// 		AttachTo(msTagsToggleCmd)
// 	cmdr.NewStringSlice().Placeholder("LIST").
// 		Titles("unset", "u", "reset").
// 		Description("and reset the others service nodes to `tag`", ``).
// 		Group("List").
// 		AttachTo(msTagsToggleCmd)
// 	cmdr.NewString().Placeholder("HOST:PORT").
// 		Titles("address", "a", "addr").
// 		Description("the address of the service (by id or name)", ``).
// 		AttachTo(msTagsToggleCmd)
//
// 	return
// }
//
// func attachModifyFlags(cmd cmdr.OptCmd) {
// 	cmdr.NewString("=").
// 		Titles("delim", "d").
// 		Description("delimitor char in `non-plain` mode.", ``).
// 		DefaultValue("=", "").
// 		AttachTo(cmd)
//
// 	cmdr.NewBool().
// 		Titles("clear", "c").
// 		Description("clear all tags.", ``).
// 		Group("Operate").
// 		AttachTo(cmd)
//
// 	cmdr.NewBool().
// 		Titles("string", "g", "string-mode").
// 		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.", ``).
// 		DefaultValue(false, "").
// 		ToggleGroup("Mode").
// 		AttachTo(cmd)
//
// 	cmdr.NewBool().
// 		Titles("meta", "m", "meta-mode").
// 		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).", ``).
// 		ToggleGroup("Mode").
// 		AttachTo(cmd)
//
// 	cmdr.NewBool().
// 		Titles("both", "2", "both-mode").
// 		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.", ``).
// 		ToggleGroup("Mode").
// 		AttachTo(cmd)
//
// 	cmdr.NewBool().
// 		Titles("plain", "p", "plain-mode").
// 		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.", ``).
// 		ToggleGroup("Mode").
// 		AttachTo(cmd)
//
// 	cmdr.NewBool(true).
// 		Titles("tag", "t", "tag-mode").
// 		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.", ``).
// 		ToggleGroup("Mode").
// 		AttachTo(cmd)
//
// }
//
// func attachConsulConnectFlags(cmd cmdr.OptCmd) {
// 	cmdr.NewString(consul.DEFAULT_CONSUL_HOST).Placeholder("HOST:PORT").
// 		Titles("addr", "a").
// 		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
// 		Group("Consul").
// 		EnvKeys("ADDR").
// 		AttachTo(cmd)
// 	cmdr.NewInt(consul.DEFAULT_CONSUL_PORT).Placeholder("PORT").
// 		Titles("port", "p").
// 		Description("Consul port", ``).
// 		Group("Consul").
// 		AttachTo(cmd)
// 	cmdr.NewBool(true).
// 		Titles("insecure", "K").
// 		Description("Skip TLS host verification", ``).
// 		Group("Consul").
// 		AttachTo(cmd)
// 	cmdr.NewString("/").Placeholder("ROOT").
// 		Titles("prefix", "px").
// 		Description("Root key prefix", ``).
// 		Group("Consul").
// 		AttachTo(cmd)
// 	cmdr.NewString().Placeholder("FILE").
// 		Titles("", "cacert").
// 		Description("Consul Client CA cert)", ``).
// 		Group("Consul").
// 		AttachTo(cmd)
// 	cmdr.NewString().Placeholder("FILE").
// 		Titles("", "cert").
// 		Description("Consul Client cert", ``).
// 		Group("Consul").
// 		AttachTo(cmd)
// 	cmdr.NewString(consul.DEFAULT_CONSUL_SCHEME).Placeholder("SCHEME").
// 		Titles("", "scheme").
// 		Description("Consul connection protocol", ``).
// 		Group("Consul").
// 		AttachTo(cmd)
// 	cmdr.NewString().Placeholder("USERNAME").
// 		Titles("username", "u", "user", "usr", "uid").
// 		Description("HTTP Basic auth user", ``).
// 		Group("Consul").AttachTo(cmd)
// 	cmdr.NewString().Placeholder("PASSWORD").
// 		Titles("password", "pw", "passwd", "pass", "pwd").
// 		Description("HTTP Basic auth password", ``).
// 		Group("Consul").
// 		ExternalTool(cmdr.ExternalToolPasswordInput).
// 		AttachTo(cmd)
// }
