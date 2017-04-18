package consul

import (
	"os"
	"path"
	log "github.com/cihub/seelog"

	"fmt"
	"gopkg.in/urfave/cli.v2"
)

const (
	DEFAULT_SEELOG_CONFIG = `<seelog type="asynctimer" asyncinterval="5000000" minlevel="info" maxlevel="critical">
    <outputs formatid="main">
        <!-- filter formatid="main-console" levels="debug,info,warn,error,critical">
            <console />
        </filter>
        <filter levels="debug,info,warn,error,critical">
            <file path="/tmp/devops.log"/>
        </filter -->
        <filter levels="trace">
            <console formatid="colored_trace"/>
        </filter>
        <filter levels="debug">
            <console formatid="colored_debug_c"/>
        </filter>
        <filter levels="info">
            <console formatid="colored_info_c"/>
        </filter>
        <filter levels="warn">
            <console formatid="colored_warn_c"/>
            <file path="/tmp/devops.log"/>
        </filter>
        <filter levels="error">
            <console formatid="colored_error_c"/>
            <file path="/tmp/devops.log"/>
        </filter>
        <filter levels="critical">
            <console formatid="colored_critical"/>
            <file path="/tmp/devops.log"/>
        </filter>
    </outputs>

    <formats>
        <format id="main-console" format="[%LEV] %Msg%n"/>
        <format id="main" format="%Date/%Time [%LEV] %Msg%n"/>
        <format id="colored_trace_1" format="%EscM(0)[%Date(2/Jan/2006 15:04:05)] %EscM(36;1)[%l]%EscM(0)%EscM(36) %Msg%n%EscM(0)"/>
        <format id="colored_trace" format="%EscM(0)[%Date(2006-01-01 15:04:05)] %EscM(36;1)[%l]%EscM(0)%EscM(36) %Msg%n%EscM(0)"/>
        <format id="colored_debug" format="%EscM(0)[%Date(2006-01-01 15:04:05.000000)] %EscM(34)[%l]%EscM(0)%EscM(34;1) %Msg%n%EscM(0)"/>
        <format id="colored_info" format="%EscM(0)[%Date(2006-01-01 15:04:05.000000)] %EscM(32;1)[%l]%EscM(0)%EscM(32) %Msg%n%EscM(0)"/>
        <format id="colored_warn" format="%EscM(0)[%Date(2006-01-01 15:04:05.000000)] %EscM(33;1)[%l]%EscM(0)%EscM(33) %Msg%n%EscM(0)"/>
        <format id="colored_error" format="%EscM(0)[%Date(2006-01-01 15:04:05.000000)] %EscM(35;1)[%l]%EscM(0)%EscM(35) %Msg%n%EscM(0)"/>
        <format id="colored_critical" format="%EscM(0)[%Date(2006-01-01 15:04:05.000000)] %EscM(31;1)[%l]%EscM(0)%EscM(31) %Msg%n%EscM(0)"/>
        <format id="colored_debug_c" format="%EscM(0)[%Date(15:04:05.000)] %EscM(34)[%l]%EscM(0)%EscM(34;1) %Msg%n%EscM(0)"/>
        <format id="colored_info_c" format="%EscM(0)[%Date(15:04:05.000)] %EscM(32;1)[%l]%EscM(0)%EscM(32) %Msg%n%EscM(0)"/>
        <format id="colored_warn_c" format="%EscM(0)[%Date(15:04:05.000)] %EscM(33;1)[%l]%EscM(0)%EscM(33) %Msg%n%EscM(0)"/>
        <format id="colored_error_c" format="%EscM(0)[%Date(15:04:05.000)] %EscM(35;1)[%l]%EscM(0)%EscM(35) %Msg%n%EscM(0)"/>
        <format id="colored_critical_c" format="%EscM(0)[%Date(15:04:05.000)] %EscM(31;1)[%l]%EscM(0)%EscM(31) %Msg%n%EscM(0)"/>
    </formats>
</seelog>`
)

func InitLogger() {
	var logger log.LoggerInterface
	var err error
	homeDir := os.Getenv("HOME")
	seelogPath := path.Join(homeDir, ".devops.seelog.xml")
	if _, err := os.Stat(seelogPath); err == nil {
		logger, err = log.LoggerFromConfigAsFile(seelogPath)
	}else {
		basicLogConfig := DEFAULT_SEELOG_CONFIG
		logger, err = log.LoggerFromConfigAsString(basicLogConfig)
	}
	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)


	//log.Trace("TRACE")
	//log.Debug("DEBUG")
	//log.Info("INFO")
	//log.Warn("WARN")
	//log.Error("ERROR")
	//log.Critical("CRITICAL")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteDefaultLogger (c *cli.Context) {
	homeDir := os.Getenv("HOME")
	seelogPath := path.Join(homeDir, ".devops.seelog.xml")

	//err := ioutil.WriteFile(seelogPath, DEFAULT_SEELOG_CONFIG, 0644)
	//check(err)
	f, err := os.Create(seelogPath)
	check(err)
	defer f.Close()
	n3, err := f.WriteString(DEFAULT_SEELOG_CONFIG)
	check(err)
	f.Sync()
	fmt.Printf("wrote %d bytes: %s\n", n3, seelogPath)
}

var LoggerCommands = []*cli.Command{
	{
		Name:  "write",
		Usage: "write default logger config to ~/.devops.seelog.xml",
		//Subcommands: []*cli.Command{
		//	{
		//		Name:   "backup",
		//		Usage:  "Dump Consul's KV database to a JSON/YAML file",
		//		Action: Backup,
		//		Flags: []cli.Flag{
		//			&cli.StringFlag{
		//				Name:    "outfile,o",
		//				Usage:   "Write output to a file (*.json / *.yml)",
		//				EnvVars: []string{"CONSUL_OUTPUT"},
		//			},
		//		},
		//	},
		//	{
		//		Name:   "restore",
		//		Usage:  "restore a JSON/YAML backup of Consul's KV store",
		//		Action: Restore,
		//	},
		//},
		Action: func(c *cli.Context) error {
			//fmt.Println("K/V pair opeations: ", c.Args().First())
			WriteDefaultLogger(c)
			return nil
		},
	},
}
