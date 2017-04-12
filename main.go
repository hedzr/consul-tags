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

	basicLogConfig := `
<seelog type="asynctimer" asyncinterval="5000000" minlevel="debug" maxlevel="critical">
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
	logger, err := log.LoggerFromConfigAsString(basicLogConfig)
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

func main() {

	// setup the flags/options within a .ops-config.yml
	var yamlConfig = cmd.NewYamlSourceFromFlagFunc("config")

	app := &cli.App{
		Name:    "consul-tags",
		Usage:   "The DevOps Toolset!",
		Version: "0.3.1",
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
