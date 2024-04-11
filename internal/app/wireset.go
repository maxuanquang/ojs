package app

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewStandaloneServer,
	NewHTTPServer,
	NewWorker,
	NewCron,
)
