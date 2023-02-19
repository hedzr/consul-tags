/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
)

func Pre(cmd *cmdr.Command, args []string) (err error) {
	//if cmdr.GetDebugMode() {
	//	cmdr.Logger.SetLevel(logrus.DebugLevel)
	//}
	cmdr.Logger.Debugf("app starting...")
	return
}

func Post(cmd *cmdr.Command, args []string) {
	cmdr.Logger.Debugf("app stopping...")
	cmdr.Logger.Debugf("app stopped.")
}
