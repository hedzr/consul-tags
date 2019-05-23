/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/sirupsen/logrus"
)

func Pre(cmd *cmdr.Command, args []string) (err error) {
	if cmdr.GetDebugMode() {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Debug("app starting...")
	return
}

func Post(cmd *cmdr.Command, args []string) {
	logrus.Debug("app stopping...")
	logrus.Debug("app stopped.")
}

