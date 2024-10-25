/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"

	"github.com/hedzr/consul-tags/objects/consul"
)

//
//
//

//
//
//

// func kvBackup(cmd *cmdr.Command, args []string) (err error) {
// 	err = consul.Backup()
// 	return
// }
//
// func kvRestore(cmd *cmdr.Command, args []string) (err error) {
// 	err = consul.Restore()
// 	return
// }

func kvBackupV2(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	err = consul.Backup()
	return
}

func kvRestoreV2(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	err = consul.Restore()
	return
}

//
//
//

// func msList(cmd *cmdr.Command, args []string) (err error) {
// 	err = consul.ServiceList()
// 	return
// }

func msListV2(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	err = consul.ServiceList()
	return
}

//
//
//

// func msTagsList(cmd *cmdr.Command, args []string) (err error) {
// 	err = consul.TagsList()
// 	return
// }
//
// func msTagsAdd(cmd *cmdr.Command, args []string) (err error) {
// 	err = consul.Tags()
// 	return
// }
//
// func msTagsRemove(cmd *cmdr.Command, args []string) (err error) {
// 	err = consul.Tags()
// 	return
// }
//
// func msTagsModify(cmd *cmdr.Command, args []string) (err error) {
// 	err = consul.Tags()
// 	return
// }
//
// func msTagsToggle(cmd *cmdr.Command, args []string) (err error) {
// 	err = consul.TagsToggle()
// 	return
// }

func msTagsListV2(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	err = consul.TagsList()
	return
}

func msTagsAddV2(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	err = consul.Tags()
	return
}

func msTagsRemoveV2(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	err = consul.Tags()
	return
}

func msTagsModifyV2(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	err = consul.Tags()
	return
}

func msTagsToggleV2(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	err = consul.TagsToggle()
	return
}
