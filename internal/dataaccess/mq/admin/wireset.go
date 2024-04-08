package admin

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAdmin,
)