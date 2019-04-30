/*
 * Copyright © 2019 Hedzr Yeh.
 */

package more_cmds

import (
	"github.com/hedzr/consul-tags/objects/consul"
)

func tagsList() (err error) {
	err = consul.TagsList()
	return
}

func tagsModify() (err error) {
	err = consul.Tags()
	return
}

func tagsToggle() (err error) {
	err = consul.TagsToggle()
	return
}

func backup() (err error) {
	err = consul.Backup()
	return
}

func restore() (err error) {
	err = consul.Restore()
	return
}
