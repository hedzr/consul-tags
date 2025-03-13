/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"fmt"

	cmdrv2 "github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

func ServiceList() (err error) {
	registrar := getRegistrarV2()
	err = listServices(registrar)
	return
}

func listServices(registrar *Registrar) error {
	// cmdr.Logger.Debugf("List the services at '%s'...", cmdr.GetStringP(TAGS_PREFIX, "addr"))

	cs2 := cmdrv2.Store().WithPrefix(MS_PREFIX_V2)
	logz.Print("List the services...", "addr", cs2.MustString("addr"))

	WaitForResult(func() (bool, error) {
		vm, qm, err := registrar.FirstClient.Catalog().Services(nil)
		if err != nil {
			// cmdr.Logger.Errorf("qm: %v, err: %v", qm, err)
			logz.Error("retrieve service list failed", "qm", qm, "err", err)
			return false, err
		}

		for k, v := range vm {
			fmt.Printf("%31s: %v\n", k, v)
		}
		return true, nil
	}, func(err error) {
		// cmdr.Logger.Fatalf("err: %v", err)
		logz.Fatal("retrieve service list failed", "err", err)
	})
	return nil
}
