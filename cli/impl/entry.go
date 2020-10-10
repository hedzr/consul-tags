/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/consul-tags"
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

	kvCmd := root.NewSubCommand().
		Titles("kvstore", "kv").
		Description("consul kv store operations...", ``)

	attachConsulConnectFlags(kvCmd)

	kvBackupCmd := kvCmd.NewSubCommand().
		Titles("backup", "b", "bk", "bf", "bkp").
		Description("Dump Consul's KV database to a JSON/YAML file", ``).
		Action(kvBackup)
	kvBackupCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("output", "o").
		Description("Write output to a file (*.json / *.yml)", ``).
		DefaultValue("consul-backup.json", "FILE")

	kvRestoreCmd := kvCmd.NewSubCommand().
		Titles("restore", "r").
		Description("restore to Consul's KV store, from a a JSON/YAML backup file", ``).
		Action(kvRestore)
	kvRestoreCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("input", "i").
		Description("Read the input file (*.json / *.yml)", ``).
		DefaultValue("consul-backup.json", "FILE")

	// ms

	msCmd := root.NewSubCommand().
		Titles("microservice", "ms", "micro-service").
		Description("micro-service operations...", ``)

	msCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("name", "n").
		Description("name of the service", ``).
		DefaultValue("", "NAME")
	msCmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("id", "i", "ID").
		Description("unique id of the service", ``).
		DefaultValue("", "ID")
	msCmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("all", "a").
		Description("all services", ``).
		DefaultValue(false, "")

	// ms list

	/* msListCmd := */
	msCmd.NewSubCommand().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list services.", ``).
		Action(msList)

	// ms tags

	msTagsCmd := kvCmd.NewSubCommand().
		Titles("tags", "t").
		Description("tags operations of a services.", ``).
		Action(msList)

	attachConsulConnectFlags(msTagsCmd)

	msTagsCmd.NewSubCommand().
		Titles("list", "ls", "l", "lst", "dir").
		Description("list tags of a service.", ``).
		Action(msTagsList)

	// ms tags modify

	msTagsModifyCmd := msTagsCmd.NewSubCommand().
		Titles("modify", "m", "mod", "modi", "update", "change").
		Description("modify tags of a service.", ``).
		Action(msTagsModify)

	attachModifyFlags(msTagsModifyCmd)

	msTagsModifyCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("add", "a", "add-list").
		Description("a comma list to be added.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	msTagsModifyCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("remove", "r", "rm-list", "rm", "del", "delete").
		Description("a comma list to be removed.", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")

	// ms tags toggle

	msTagsToggleCmd := msTagsCmd.NewSubCommand().
		Titles("toggle", "t", "tog", "switch").
		Description("toggle tags of a service.", ``).
		Action(msTagsToggle)

	attachModifyFlags(msTagsToggleCmd)

	msTagsToggleCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("set", "s").
		Description("set to `tag` which service specified by --address", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	msTagsToggleCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("unset", "u", "reset").
		Description("and reset the others service nodes to `tag`", ``).
		DefaultValue([]string{}, "LIST").
		Group("List")
	msTagsToggleCmd.NewFlag(cmdr.OptFlagTypeStringSlice).
		Titles("address", "a", "addr").
		Description("the address of the service (by id or name)", ``).
		DefaultValue("", "HOST:PORT")

	return
}

func attachModifyFlags(cmd cmdr.OptCmd) {
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("delim", "d").
		Description("delimitor char in `non-plain` mode.", ``).
		DefaultValue("=", "")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("clear", "c").
		Description("clear all tags.", ``).
		DefaultValue(false, "").
		Group("Operate")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("string", "g", "string-mode").
		Description("In 'String Mode', default will be disabled: default, a tag string will be split by comma(,), and treated as a string list.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("meta", "m", "meta-mode").
		Description("In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false).", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("both", "2", "both-mode").
		Description("In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("plain", "p", "plain-mode").
		Description("In 'Plain Mode', a tag be NOT treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		DefaultValue(false, "").
		Group("Mode")

	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("tag", "t", "tag-mode").
		Description("In 'Tag Mode', a tag be treated as `key=value` or `key:value`, and modify with the `key`.", ``).
		DefaultValue(true, "").
		Group("Mode")

}

func attachConsulConnectFlags(cmd cmdr.OptCmd) {
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("addr", "a").
		Description("Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')", ``).
		DefaultValue(consul.DEFAULT_CONSUL_HOST, "HOST[:PORT]").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeInt).
		Titles("port", "p").
		Description("Consul port", ``).
		DefaultValue(consul.DEFAULT_CONSUL_PORT, "PORT").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeBool).
		Titles("insecure", "K").
		Description("Skip TLS host verification", ``).
		DefaultValue(true, "").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("prefix", "px").
		Description("Root key prefix", ``).
		DefaultValue("/", "ROOT").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "cacert").
		Description("Consul Client CA cert)", ``).
		DefaultValue("", "FILE").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "cert").
		Description("Consul Client cert", ``).
		DefaultValue("", "FILE").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("", "scheme").
		Description("Consul connection protocol", ``).
		DefaultValue(consul.DEFAULT_CONSUL_SCHEME, "SCHEME").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("username", "u", "user", "usr", "uid").
		Description("HTTP Basic auth user", ``).
		DefaultValue("", "USERNAME").
		Group("Consul")
	cmd.NewFlag(cmdr.OptFlagTypeString).
		Titles("password", "pw", "passwd", "pass", "pwd").
		Description("HTTP Basic auth password", ``).
		DefaultValue("", "PASSWORD").
		Group("Consul").
		ExternalTool(cmdr.ExternalToolPasswordInput)

}
