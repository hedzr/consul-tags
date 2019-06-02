/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/sirupsen/logrus"
)

func ServiceList() (err error) {
	registrar := getRegistrar()
	err = listServices(registrar)
	return
}

func listServices(registrar *Registrar) error {
	logrus.Debugf("List the services at '%s'...", cmdr.GetStringP(TAGS_PREFIX, "addr"))

	WaitForResult(func() (bool, error) {
		vm, qm, err := registrar.FirstClient.Catalog().Services(nil)
		if err != nil {
			logrus.Errorf("qm: %v, err: %v", qm, err)
			return false, err
		}

		for k, v := range vm {
			fmt.Printf("%31s: %v\n", k, v)
		}
		return true, nil
	}, func(err error) {
		logrus.Fatalf("err: %v", err)
	})
	return nil
}
