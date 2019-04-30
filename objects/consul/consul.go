/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2" // imports as package "cli"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_CONSUL_HOST      = "consul.ops.local"
	DEFAULT_CONSUL_LOCALHOST = "localhost"
	DEFAULT_CONSUL_PORT      = 8500
	DEFAULT_CONSUL_SCHEME    = "http"
)

type valueEnc struct {
	Encoding string `json:"encoding,omitempty"`
	Str      string `json:"value"`
}

type kvJSON struct {
	BackupDate time.Time           `json:"date"`
	Connection map[string]string   `json:"connection_info"`
	Values     map[string]valueEnc `json:"values"`
}

type Glyph struct {
	Name     string
	Position int
	ID       int // ID must > 0
	ParentID int // if ParentID equal 0, the Glyph has no parents
}

type InputZ interface {
	Int(name string) (int, error)
	CoolPrint(str string) (int, error)
}

type MyInputZ struct {
	Ok string
}

func (s MyInputZ) Int(name string) (int, error) {
	return 0, nil
}

//
//
// func (s *InputZ) Int(name string) (int, error) {
// 	//
// }

// var consul_addr = ""
// var gingerCrouton = false

var Flags = []cli.Flag{
	// altsrc.NewIntFlag(cli.IntFlag{Name: "test"}),

	&cli.StringFlag{
		Name:    "addr",
		Aliases: []string{"a"},
		Value:   DEFAULT_CONSUL_HOST,
		// Destination: &consul_addr,
		Usage:   "Consul address and port: `HOST[:PORT]` (No leading 'http(s)://')",
		EnvVars: []string{"CONSUL_ADDR"},
	},
	&cli.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   DEFAULT_CONSUL_PORT,
		Usage:   "Consul port",
		EnvVars: []string{"CONSUL_PORT"},
	},
	&cli.StringFlag{
		Name:    "prefix",
		Value:   "/",
		Usage:   "Root key prefix",
		EnvVars: []string{"CONSUL_PREFIX"},
	},
	&cli.StringFlag{
		Name:    "cacert",
		Aliases: []string{"r"},
		Usage:   "Client CA cert",
		EnvVars: []string{"CONSUL_CA_CERT"},
	},
	&cli.StringFlag{
		Name:    "cert",
		Aliases: []string{"t"},
		Usage:   "Client cert",
		EnvVars: []string{"CONSUL_CERT"},
	},
	&cli.StringFlag{
		Name:    "scheme",
		Aliases: []string{"s"},
		Value:   "http",
		Usage:   "Consul connection scheme (HTTP or HTTPS)",
		EnvVars: []string{"CONSUL_SCHEME"},
	},
	&cli.BoolFlag{
		Name:    "insecure",
		Aliases: []string{"K"},
		Usage:   "Skip TLS host verification",
		EnvVars: []string{"CONSUL_INSECURE"},
	},
	&cli.StringFlag{
		Name:    "username", // short name can be 2 chars, best 1 char
		Aliases: []string{"U"},
		Usage:   "HTTP Basic auth user",
		EnvVars: []string{"CONSUL_USER"},
	},
	&cli.StringFlag{
		Name:    "password", // short name can be 2 chars, best 1 char
		Aliases: []string{"P"},
		Usage:   "HTTP Basic auth password",
		EnvVars: []string{"CONSUL_PASS"},
	},
	&cli.StringFlag{
		Name:    "key",
		Aliases: []string{"Y"},
		Usage:   "Client key",
		EnvVars: []string{"CONSUL_KEY"},
	},
}

func Before(c *cli.Context) error {
	// logrus.Printf("**** consul.Before()\n")
	// logrus.Info("**** consul.Before()")
	logrus.Trace("**** consul.Before()")
	addr := c.String("addr")

	if addr == "" {
		addr = DEFAULT_CONSUL_HOST + ":" + strconv.Itoa(DEFAULT_CONSUL_PORT)
	}
	if strings.Index(addr, ":") < 0 {
		addr = addr + ":" + c.String("port")
	}
	c.Set("addr", addr)
	// logrus.Info("**** Formal --addr='%s'", addr)
	logrus.Tracef("**** Formal --addr='%s'", addr)
	return nil
}

var Commands = []*cli.Command{
	{
		Name:  "kv",
		Usage: "K/V pair Operations, ...",
		Subcommands: []*cli.Command{
			{
				Name:  "backup",
				Usage: "Dump Consul's KV database to a JSON/YAML file",
				// Action: Backup,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "outfile,o",
						Usage:   "Write output to a file (*.json / *.yml)",
						EnvVars: []string{"CONSUL_OUTPUT"},
					},
				},
			},
			{
				Name:  "restore",
				Usage: "restore a JSON/YAML backup of Consul's KV store",
				// Action: Restore,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("K/V pair opeations: ", c.Args().First())
			return nil
		},
	},
	{
		Name:    "ms",
		Aliases: []string{"service"},
		Usage:   "Microservice Operations, ...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "",
				Usage:   "specify the service `Name`",
			},
			&cli.StringFlag{
				Name:    "id",
				Aliases: []string{"i"},
				Value:   "",
				Usage:   "specify the service `ID`",
			},
		},
		Subcommands: []*cli.Command{
			// {
			// 	Name:  "test",
			// 	Usage: "make a test (consulapi, consul)",
			// 	Action: func(c *cli.Context) error {
			// 		Test(c)
			// 		return nil
			// 	},
			// },
			// {
			// 	Name:  "test1",
			// 	Usage: "make a test 1 (redis tags modify)",
			// 	Action: func(c *cli.Context) error {
			// 		Test1(c)
			// 		return nil
			// 	},
			// },
			{
				Name:  "ls",
				Usage: "list service definitions by its name or id",
				// Action: TagsList,
				Description: `Normal Usages:

    $ devops consul ms --name serive-name ls
      list the tags of service nodes by name
    $ devops consul ms --id serive-id ls
      list the tags of service nodes by id
`,
			},
			{
				Name:      "tags",
				Usage:     "maintains a service's tags",
				UsageText: "doo - does the dooing",
				Description: `'tags' was designed for maintaining the tags removing and
    appending in one line.
    notice that '--clear' has first prior, so you can:
    $ devops consul ms --name serive-name tags --clear --add x,xx,xx --rm y
    $ devops consul ms --id serive-id tags --clear --add x,xx,xx --rm y

    Normal Usages:

    --string string mode or list mode:
        string mode: '--add x=1,xx' will be treated as a tag string: ["x=1,xx"]
        list   mode: '--add x=1,xx' will be treated as a tag list: ["x=1", "xx"]
    --plain  plain mode or k/v mode:
      string mode:
        plain mode:  '--add x=1,xx' => ["x=1,xx"]
        k/v   mode:  '--add x=1,xx' => ["x=1,xx"] => [ {"x":"1,xx"}
      list mode:
        plain mode:  '--add x=1,xx' => ["x=1","xx"]
        k/v   mode:  '--add x=1,xx' => ["x=1","xx"] => [ {"x":"1"}, {"xx":""} ]

    $ devops consul ms --name serive-name tags -r x,xx,xxx --plain --string
      remove the tag 'x,xx,xxx'
    $ devops consul ms --name serive-name tags -r x,xx,xxx --plain
      remove the tags 'x', 'xx' and 'xxx'
    $ devops consul ms --name serive-name tags -a role=master,xx,xxx --string
      find all tags with prefix 'role=', replace with 'role=master,xx,xxx'
    $ devops consul ms --name serive-name tags -r role=master,xx,xxx
      find all tags with prefix 'role=', replace with 'role=master', and, append tags 'xx' and 'xxx' if not exists.
    $ devops consul ms --name serive-name tags ls
      list the tags of service nodes by name
    $ devops consul ms --id serive-id tags ls
      list the tags of service nodes by id

    DevOps Usages:

    $ devops consul ms --name test-service tags toggle --address 'ip:port' --set "role=master" --reset "role=slave"
`,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "clear",
						Aliases: []string{"c"},
						Value:   false,
						Usage:   "clear all tags",
					},
					&cli.StringSliceFlag{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "add `tag,...`",
					},
					&cli.StringSliceFlag{
						Name:    "rm",
						Aliases: []string{"r"},
						Usage:   "remove `tag,...`",
					},
					&cli.BoolFlag{
						Name:    "plain",
						Aliases: []string{"p"},
						Usage:   "In plain mode a tag do NOT treated as `key=value` or `key:value`, and modify with the `key`",
					},
					&cli.StringFlag{
						Name:    "delim",
						Aliases: []string{"d"},
						Value:   "=",
						Usage:   "delimitor char in `non-plain` mode",
					},
					&cli.BoolFlag{
						Name:    "string",
						Aliases: []string{"s"},
						Usage:   "Default a tag string will be split by comma(,), and treated as a string list; but string mode disable this.",
					},
					&cli.BoolFlag{
						Name:    "meta",
						Aliases: []string{"m"},
						Usage:   "In 'Meta Mode', service 'NodeMeta' field will be updated instead of 'Tags'. (--plain assumed false)",
					},
					&cli.BoolFlag{
						Name:    "both",
						Aliases: []string{"2"},
						Usage:   "In 'Both Mode', both of 'NodeMeta' and 'Tags' field will be updated.",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "ls",
						Usage: "list service tags by its name or id",
						// Action: TagsList,
						Description: `Normal Usages:

    $ devops consul ms --name serive-name tags ls
      list the tags of service nodes by name
    $ devops consul ms --id serive-id tags ls
      list the tags of service nodes by id
`,
					},
					{
						Name:  "toggle",
						Usage: "Toggle service tags, master node to something and slaves nodes to others",
						// Action: TagsToggle,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "address",
								Aliases: []string{"x"},
								Value:   "",
								Usage:   "specify the service `ID`",
							},
							&cli.StringSliceFlag{
								Name:    "set",
								Aliases: []string{"s"},
								Usage:   "set to `tag` which service specified by --address",
								EnvVars: []string{"CONSUL_OUTPUT"},
							},
							&cli.StringSliceFlag{
								Name:    "reset",
								Aliases: []string{"r"},
								Usage:   "and reset the others service nodes to `tag`",
								EnvVars: []string{"CONSUL_OUTPUT"},
							},
						},
						Description: `DevOps Usages:

    $ devops consul ms --name test-service tags toggle --address 'ip:port' --set "role=master" --reset "role=slave"
`,
					},
				},
				Action: func(c *cli.Context) error {
					// Tags(c)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Microservice opeations: ", c.Args().First())
			return nil
		},
	},

	{
		Name:    "logger",
		Aliases: []string{"l"},
		Usage:   "logger helpers",
		// SubFlags: []cli.Flag { {},{},...   },
		// Flags:       consul.Flags,
		// Subcommands: LoggerCommands,
		// Before:      consul.Before,
		Action: func(c *cli.Context) error {
			fmt.Println("logger task: ", c.Args().First())
			return nil
		},
	},
}

func init() {
}
