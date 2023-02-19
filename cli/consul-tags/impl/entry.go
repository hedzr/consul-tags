/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	consul_tags "github.com/hedzr/consul-tags"
	"github.com/hedzr/consul-tags/objects/consul"
	"github.com/hedzr/logex/logx/logrus"
)

func Entry() {

	if err := cmdr.Exec(buildRootCommand(),
		cmdr.WithEnvPrefix("CT"),
		cmdr.WithLogx(logrus.New("debug", false, true)),
		cmdr.WithWatchMainConfigFileToo(true),
	); err != nil {
		cmdr.Logger.Errorf("Error: %v", err)
	}

}

func buildRootCommand() (rootCmd *cmdr.RootCommand) {
	root := cmdr.Root(consul_tags.AppName, consul_tags.Version).
		Copyright("consul-tags is an effective devops tool", "Hedzr Yeh <hedzrz@gmail.com>").
		PreAction(Pre).PostAction(Post)
	rootCmd = root.RootCommand()

	// kv

	kvCmd := cmdr.NewSubCmd().
		Titles("kvstore", "kv").
		Description("consul kv store operations...", ``).
		AttachTo(root)

	attachConsulConnectFlags(kvCmd)

	kvBackupCmd := cmdr.NewSubCmd().
		Titles("backup", "b", "bk", "bf", "bkp").
		Description("Dump Consul's KV database to a JSON/YAML file", ``).
		Action(kvBackup).
		AttachTo(kvCmd)
	cmdr.NewString("consul-backup.json").Placeholder("FILE").
		Titles("output", "o").
		Description("Write output to a file (*.json / *.yml)", ``).
		AttachTo(kvBackupCmd)

	kvRestoreCmd := cmdr.NewSubCmd().
		Titles("restore", "r").
		Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
		Action(kvRestore).
		AttachTo(kvCmd)
	cmdr.NewString("consul-backup.json").Placeholder("FILE").
		Titles("input", "i").
		Description("Read the input file (*.json / *.yml)", ``).
		AttachTo(kvRestoreCmd)

	// ms

	msCmd := cmdr.NewSubCmd().
		Titles("microservice", "ms", "micro-service").
		Description("micro-service operations...", ``).
		AttachTo(root)

	cmdr.NewString().Placeholder("NAME").
		Titles("name", "n").
		Description("name of the service", ``).
		DefaultValue("", "NAME").
		AttachTo(msCmd)
	cmdr.NewString().Placeholder("ID").
		Titles("id", "i", "ID").
		Description("unique id of the service", ``).
		AttachTo(msCmd)
	cmdr.NewBool().
		Titles("all", "a").
		Description("all services", ``).
		AttachTo(msCmd)

	// ms list

	/* msListCmd := */
	cmdr.NewSubCmd().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list services.", ``).
		Action(msList).
		AttachTo(msCmd)

	// ms tags

	msTagsCmd := cmdr.NewSubCmd().
		Titles("tags", "t").
		Description("tags operations of a services.", ``).
		Action(msList).AttachTo(kvCmd)

	attachConsulConnectFlags(msTagsCmd)

	cmdr.NewSubCmd().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list tags of a service.", ``).
		Action(msTagsList).
		AttachTo(msTagsCmd)

	// ms tags modify

	msTagsModifyCmd := cmdr.NewSubCmd().
		Titles("modify", "m", "mod", "modi", "update", "change").
		Description("modify tags of a service.", ``).
		Action(msTagsModify).
		AttachTo(msTagsCmd)

	attachModifyFlags(msTagsModifyCmd)

	cmdr.NewStringSlice().Placeholder("LIST").
		Titles("add", "a", "add-list").
		Description("a comma list to be added.", ``).
		Group("List").
		AttachTo(msTagsModifyCmd)
	cmdr.NewStringSlice().Placeholder("LIST").
		Titles("remove", "r", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		Group("List").
		AttachTo(msTagsModifyCmd)

	// ms tags toggle

	msTagsToggleCmd := cmdr.NewSubCmd().
		Titles("toggle", "t", "tog", "switch").
		Description("toggle tags of a service.", ``).
		Action(msTagsToggle).
		AttachTo(msTagsCmd)

	attachModifyFlags(msTagsToggleCmd)

	cmdr.NewStringSlice().Placeholder("LIST").
		Titles("set", "s").
		Description("set to `tag` which service specified by --address", ``).
		Group("List").
		AttachTo(msTagsToggleCmd)
	cmdr.NewStringSlice().Placeholder("LIST").
		Titles("unset", "u", "reset").
		Description("and reset the others service nodes to `tag`", ``).
		Group("List").
		AttachTo(msTagsToggleCmd)
	cmdr.NewString().Placeholder("HOST:PORT").
		Titles("address", "a", "addr").
		Description("the address of the service (by id or name)", ``).
		AttachTo(msTagsToggleCmd)

	return
}

func attachModifyFlags(cmd cmdr.OptCmd) {
	cmdr.NewString("=").
		Titles("delim", "d").
		Description("delimitor char in `non-plain` mode.", ``).
		DefaultValue("=", "").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("clear", "c").
		Description("clear all tags.", ``).
		Group("Operate").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("string", "g", "string-mode").
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.", ``).
		DefaultValue(false, "").
		ToggleGroup("Mode").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("meta", "m", "meta-mode").
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).", ``).
		ToggleGroup("Mode").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("both", "2", "both-mode").
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.", ``).
		ToggleGroup("Mode").
		AttachTo(cmd)

	cmdr.NewBool().
		Titles("plain", "p", "plain-mode").
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		ToggleGroup("Mode").
		AttachTo(cmd)

	cmdr.NewBool(true).
		Titles("tag", "t", "tag-mode").
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		ToggleGroup("Mode").
		AttachTo(cmd)

}

func attachConsulConnectFlags(cmd cmdr.OptCmd) {
	cmdr.NewString(consul.DEFAULT_CONSUL_HOST).Placeholder("HOST:PORT").
		Titles("addr", "a").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewInt(consul.DEFAULT_CONSUL_PORT).Placeholder("PORT").
		Titles("port", "p").
		Description("Consul port", ``).
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewBool(true).
		Titles("insecure", "K").
		Description("Skip TLS host verification", ``).
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString("/").Placeholder("ROOT").
		Titles("prefix", "px").
		Description("Root key prefix", ``).
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString().Placeholder("FILE").
		Titles("", "cacert").
		Description("Consul Client CA cert)", ``).
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString().Placeholder("FILE").
		Titles("", "cert").
		Description("Consul Client cert", ``).
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString(consul.DEFAULT_CONSUL_SCHEME).Placeholder("SCHEME").
		Titles("", "scheme").
		Description("Consul connection protocol", ``).
		Group("Consul").
		AttachTo(cmd)
	cmdr.NewString().Placeholder("USERNAME").
		Titles("username", "u", "user", "usr", "uid").
		Description("HTTP Basic auth user", ``).
		Group("Consul").AttachTo(cmd)
	cmdr.NewString().Placeholder("PASSWORD").
		Titles("password", "pw", "passwd", "pass", "pwd").
		Description("HTTP Basic auth password", ``).
		Group("Consul").
		ExternalTool(cmdr.ExternalToolPasswordInput).
		AttachTo(cmd)
}
