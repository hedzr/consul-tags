/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/consul-tags/objects/consul"
)

//
//
//

//
//
//

func kvBackup(cmd *cmdr.Command, args []string) (err error) {
	err = consul.Backup()
	return
}

func kvRestore(cmd *cmdr.Command, args []string) (err error) {
	err = consul.Restore()
	return
}

//
//
//

func msList(cmd *cmdr.Command, args []string) (err error) {
	err = consul.ServiceList()
	return
}

//
//
//

func msTagsList(cmd *cmdr.Command, args []string) (err error) {
	err = consul.TagsList()
	return
}

func msTagsAdd(cmd *cmdr.Command, args []string) (err error) {
	err = consul.Tags()
	return
}

func msTagsRemove(cmd *cmdr.Command, args []string) (err error) {
	err = consul.Tags()
	return
}

func msTagsModify(cmd *cmdr.Command, args []string) (err error) {
	err = consul.Tags()
	return
}

func msTagsToggle(cmd *cmdr.Command, args []string) (err error) {
	err = consul.TagsToggle()
	return
}
