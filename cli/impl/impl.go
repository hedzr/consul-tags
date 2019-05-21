/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"fmt"
	"github.com/hedzr/common/cli-tpl/cli_common"
	"github.com/hedzr/common/cli-tpl/cmd"
	"github.com/hedzr/common/cli-tpl/cmd/daemonx"
	"github.com/hedzr/consul-tags"
	"github.com/hedzr/consul-tags/cli/more_cmds"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/takama/daemon"
)

func init() {
	// viper_init(consul_tags.APP_NAME, consul_tags.Version)
}

func Main() {
	// main_ws()
	mainCliTpl()
}

func mainCliTpl() {
	cmd.SetAppName(consul_tags.APP_NAME, consul_tags.Version, "CT")
	// cmd.SetConfigReloader(configReloader)
	cmd.SetRealServerStart(realStart, deregister)
	cmd.SetPrintVersion(printVersion)

	daemonx.Enable(cmd.RootCmd, &TheService{})
	more_cmds.Enable(cmd.RootCmd)

	cmd.Execute()
}

type TheService struct {
}

// Install the service into the system
func (s *TheService) Install(args ...string) (msg string, err error) {
	return
}

// Remove the service and all corresponding files from the system
func (s *TheService) Remove() (msg string, err error) {
	return
}

// Start the service
func (s *TheService) Start() (msg string, err error) {
	return
}

// Stop the service
func (s *TheService) Stop() (msg string, err error) {
	return
}

// Status - check the service status
func (s *TheService) Status() (msg string, err error) {
	return
}

// Run - run executable service
func (s *TheService) Run(e daemon.Executable) (msg string, err error) {
	return
}

func realStart() {
	logrus.Infof("%s v%v realStart at %v", cli_common.AppName, cli_common.Version, viper.GetInt("app.port"))
}

func deregister() {
	logrus.Infof("    Stopped...")
	// xs.Deregister() // deregister self from consul/etcd registrar (13.registrar.yml)
}

func printVersion() {
	fmt.Printf(`%15s Version: %s
	        Githash: %s
	       Build at: %s
`,
		cli_common.AppName, cli_common.Version, cli_common.Githash, cli_common.Buildstamp)
}
