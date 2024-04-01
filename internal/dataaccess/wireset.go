package dataaccess

import (
	"github.com/google/wire"
	"github.com/maxuanquang/ojs/internal/dataaccess/cache"
	"github.com/maxuanquang/ojs/internal/dataaccess/database"
	"github.com/maxuanquang/ojs/internal/dataaccess/mq"
)

var WireSet = wire.NewSet(
	database.WireSet,
	cache.WireSet,
	mq.WireSet,
)
