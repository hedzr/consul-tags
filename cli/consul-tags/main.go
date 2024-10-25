/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr-loaders/local"
	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/store"

	"github.com/hedzr/consul-tags/cli/consul-tags/impl"
)

func main() {
	// dbglog.Debug("log/slog/dbglog: app starting...", "args", os.Args)
	// defer func() { dbglog.Debug("log/slog/dbglog: app terminated.", "args", os.Args) }()

	//

	//

	//

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := impl.PrepareApp(
		cmdr.WithStore(store.New()), // create a standard Store instead of internal dummyStore
		cmdr.WithExternalLoaders(
			local.NewConfigFileLoader(),
			local.NewEnvVarLoader(),
		),

		// cmdr.WithTasksBeforeRun(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
		// 	dbglog.DebugContext(ctx, "command running...", "cmd", cmd, "runner", runner, "extras", extras)
		// 	// dbglog.WarnContext(ctx, "warn is verified")
		// 	// dbglog.ErrorContext(ctx, "error is verified")
		// 	return
		// }),

		// true for debug in developing time, it'll disable onAction on each Cmd.
		// for productive mode, comment this line.
		// cmdr.WithForceDefaultAction(true),

		// cmdr.WithSortInHelpScreen(true),       // default it's false
		// cmdr.WithDontGroupInHelpScreen(false), // default it's false

		// cmdr.WithTasksSetupPeripherals(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
		// 	obj := new(Obj)
		// 	basics.RegisterPeripheral(obj.Init(ctx))
		// 	return
		// }),
	)

	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err)
		os.Exit(app.SuggestRetCode())
	}
}

// func main1() {
// 	impl.Entry()
// }
