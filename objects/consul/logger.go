/*
 * Copyright © 2019 Hedzr Yeh.
 */

package consul

func check(e error) {
	if e != nil {
		panic(e)
	}
}
