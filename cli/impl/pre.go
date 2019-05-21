/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/sirupsen/logrus"
)

func Pre(cmd *cmdr.Command, args []string) (err error) {
	logrus.Info("app starting...")
	return
}

func Post(cmd *cmdr.Command, args []string) {
	logrus.Info("app stopping...")
	logrus.Info("app stopped.")
}

