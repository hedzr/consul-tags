/*
 * Copyright © 2019 Hedzr Yeh.
 */

package more_cmds

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Enable(rootCmd *cobra.Command) {
	rootCmd.AddCommand(kvCmd, msCmd)

	msCmd.AddCommand(msListCmd, msTagsCmd)
	msTagsCmd.AddCommand(msTagsListCmd, msTagsToggleCmd, msTagsModifyCmd)

	kvCmd.AddCommand(kvBackupCmd, kvRestoreCmd)

	flags := kvBackupCmd.PersistentFlags()

	flags.StringP("output", "o", "consul-backup.json", "Write output to a file (*.json / *.yml)")
	_ = viper.BindPFlag("app.kv.output", flags.Lookup("output"))

	flags = kvRestoreCmd.PersistentFlags()

	flags.StringP("input", "i", "consul-backup.json", "read the input file (*.json / *.yml)")
	_ = viper.BindPFlag("app.kv.input", flags.Lookup("input"))

	//
	//
	//

	flags = msCmd.PersistentFlags()

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bashCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bashCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// flags.String("addr", "localhost:5002", "Address of the service")

	// flags.BoolP("debug", "D", false, "Turn on debugging.")
	// _ = viper.BindPFlag("debug", flags.Lookup("debug"))

	flags.StringP("addr", "", "consul.ops.local", "Consul ip/host and port: HOST[:PORT] (No leading 'http(s)://')")
	flags.IntP("port", "p", 8500, "Consul port")
	flags.BoolP("insecure", "K", true, "Skip TLS host verification")
	flags.StringP("prefix", "", "/", "Root key prefix")
	flags.StringP("cacert", "", "", "Client CA cert")
	flags.StringP("cert", "", "", "Client cert")
	flags.StringP("scheme", "", "", "Consul connection scheme (HTTP or HTTPS)")
	flags.StringP("username", "U", "", "HTTP Basic auth user")
	flags.StringP("password", "P", "", "HTTP Basic auth password")
	flags.StringP("key", "", "", "Client key")

	// $CT_APP_MS_ADDR
	_ = viper.BindPFlag("app.ms.addr", flags.Lookup("addr"))
	// $CT_APP_MS_PORT or yml['app.ms.port']
	_ = viper.BindPFlag("app.ms.port", flags.Lookup("port"))
	_ = viper.BindPFlag("app.ms.insecure", flags.Lookup("insecure"))
	_ = viper.BindPFlag("app.ms.prefix", flags.Lookup("prefix"))
	_ = viper.BindPFlag("app.ms.cacert", flags.Lookup("cacert"))
	_ = viper.BindPFlag("app.ms.cert", flags.Lookup("cert"))
	_ = viper.BindPFlag("app.ms.scheme", flags.Lookup("scheme"))
	_ = viper.BindPFlag("app.ms.username", flags.Lookup("username"))
	_ = viper.BindPFlag("app.ms.password", flags.Lookup("password"))
	_ = viper.BindPFlag("app.ms.key", flags.Lookup("key"))

	flags.StringP("name", "n", "", "specify the service `Name`")
	_ = viper.BindPFlag("app.ms.name", flags.Lookup("name"))

	flags.StringP("id", "", "", "specify the service `ID`")
	_ = viper.BindPFlag("app.ms.id", flags.Lookup("id"))

	//

	// 'tags' was designed for maintaining the tags removing and
	// appending in one line.
	// 	notice that '--clear' has first prior, so you can:
	// $ devops consul ms --name serive-name tags --clear --add x,xx,xx --rm y
	// $ devops consul ms --id serive-id tags --clear --add x,xx,xx --rm y
	//
	// Normal Usages:
	//
	// --string string mode or list mode:
	// string mode: '--add x=1,xx' will be treated as a tag string: ["x=1,xx"]
	// list   mode: '--add x=1,xx' will be treated as a tag list: ["x=1", "xx"]
	// --plain  plain mode or k/v mode:
	// string mode:
	// plain mode:  '--add x=1,xx' => ["x=1,xx"]
	// k/v   mode:  '--add x=1,xx' => ["x=1,xx"] => [ {"x":"1,xx"}
	// list mode:
	// plain mode:  '--add x=1,xx' => ["x=1","xx"]
	// k/v   mode:  '--add x=1,xx' => ["x=1","xx"] => [ {"x":"1"}, {"xx":""} ]
	//
	// $ devops consul ms --name serive-name tags -r x,xx,xxx --plain --string
	// remove the tag 'x,xx,xxx'
	// $ devops consul ms --name serive-name tags -r x,xx,xxx --plain
	// remove the tags 'x', 'xx' and 'xxx'
	// $ devops consul ms --name serive-name tags -a role=master,xx,xxx --string
	// find all tags with prefix 'role=', replace with 'role=master,xx,xxx'
	// $ devops consul ms --name serive-name tags -r role=master,xx,xxx
	// find all tags with prefix 'role=', replace with 'role=master', and, append tags 'xx' and 'xxx' if not exists.
	// $ devops consul ms --name serive-name tags ls
	// list the tags of service nodes by name
	// $ devops consul ms --id serive-id tags ls
	// list the tags of service nodes by id
	//
	// DevOps Usages:
	//
	// $ devops consul ms --name test-service tags toggle --address 'ip:port' --set "role=master" --reset "role=slave"

	// $ devops consul ms --name serive-name tags ls
	// list the tags of service nodes by name
	// $ devops consul ms --id serive-id tags ls
	// list the tags of service nodes by id
	flags = msTagsModifyCmd.PersistentFlags()

	flags.BoolP("clear", "c", false, "Clear All Tags")
	flags.BoolP("plain", "", false, "Plain Mode: in this mode a tag do NOT treated as `key=value` or `key:value`, and it could be modified with the `key`")
	flags.BoolP("string", "", false, "String Mode: a tag string will be split by comma(,) generally, and treated as a string list; but string mode disable this.")
	flags.BoolP("meta", "", false, "Meta Mode: in this mode, service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false)")
	flags.BoolP("both", "", false, "Both Mode: in this mode, both of 'NodeMeta' and 'Tags' field will be updated")
	flags.StringSliceP("add", "a", []string{}, "add `tag,..`")
	flags.StringSliceP("rm", "r", []string{}, "remove `tag,..`")
	flags.StringP("delim", "d", "", "delimitor char in `non-plain` mode")

	_ = viper.BindPFlag("app.ms.tags.clear", flags.Lookup("clear"))
	_ = viper.BindPFlag("app.ms.tags.plain-mode", flags.Lookup("plain"))
	_ = viper.BindPFlag("app.ms.tags.string-mode", flags.Lookup("string"))
	_ = viper.BindPFlag("app.ms.tags.meta-mode", flags.Lookup("meta"))
	_ = viper.BindPFlag("app.ms.tags.both-mode", flags.Lookup("both"))
	_ = viper.BindPFlag("app.ms.tags.add-list", flags.Lookup("add"))
	_ = viper.BindPFlag("app.ms.tags.rm-list", flags.Lookup("rm"))
	_ = viper.BindPFlag("app.ms.tags.delim", flags.Lookup("delim"))

	// $ devops consul ms --name test-service tags toggle --address 'ip:port' --set "role=master" --reset "role=slave"
	flags = msTagsToggleCmd.PersistentFlags()

	flags.StringP("address", "x", "", "the address of the service (by id or name)")
	flags.StringSliceP("set", "t", []string{}, "set to `tag` which service specified by --address")
	flags.StringSliceP("reset", "e", []string{}, "and reset the others service nodes to `tag`")

	_ = viper.BindPFlag("app.ms.tags.toggle.address", flags.Lookup("address"))
	_ = viper.BindPFlag("app.ms.tags.toggle.set-list", flags.Lookup("set"))
	_ = viper.BindPFlag("app.ms.tags.toggle.reset-list", flags.Lookup("reset"))
}

var kvCmd = &cobra.Command{
	Use:     "kv",
	Aliases: []string{"KV"},
	Short:   "K/V operations: get|put...",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var kvBackupCmd = &cobra.Command{
	Use:     "backup",
	Aliases: []string{"bf", "bk", "bkp", "b"},
	Short:   "Dump Consul's KV database to a JSON/YAML file",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := backup(); err != nil {
			logrus.Errorf("Error: %v", err)
		}
	},
}

var kvRestoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"rest", "r"},
	Short:   "restore to Consul's KV store, from a a JSON/YAML backup file",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := restore(); err != nil {
			logrus.Errorf("Error: %v", err)
		}
	},
}

var msCmd = &cobra.Command{
	Use:     "ms",
	Aliases: []string{"service", "m"},
	Short:   "Microservices operations: tags|...",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		// cobra.Z

		fmt.Printf(`
		ADDR: %v
		PORT: %v
		`, viper.GetString("app.ms.addr"), viper.GetInt("app.ms.port"))
	},
}

var msListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List microservices",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tagsList(); err != nil {
			logrus.Errorf("Error: %v", err)
		}
	},
}

var msTagsCmd = &cobra.Command{
	Use:     "tags",
	Aliases: []string{"tag"},
	Short:   "Tags operations: get|put...",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var msTagsListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "lst", "l"},
	Short:   "Tags operations: get|put...",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tagsList(); err != nil {
			logrus.Errorf("Error: %v", err)
		}
	},
}

var msTagsModifyCmd = &cobra.Command{
	Use:     "modify",
	Aliases: []string{"update", "change", "mod", "m"},
	Short:   "Tags operations: Add/Remove...",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tagsModify(); err != nil {
			logrus.Errorf("Error: %v", err)
		}
	},
}

var msTagsToggleCmd = &cobra.Command{
	Use:     "toggle",
	Aliases: []string{"t", "tog"},
	Short:   "Tags operations: get|put...",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tagsToggle(); err != nil {
			logrus.Errorf("Error: %v", err)
		}
	},
}

var tmplCmd = &cobra.Command{
	Use:     "tmpl",
	Aliases: []string{"tmp"},
	Short:   "Template operations: get|put...",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
