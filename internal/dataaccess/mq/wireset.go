package mq

import (
	"github.com/google/wire"
	"github.com/maxuanquang/ojs/internal/dataaccess/mq/admin"
	"github.com/maxuanquang/ojs/internal/dataaccess/mq/consumer"
	"github.com/maxuanquang/ojs/internal/dataaccess/mq/producer"
)

var WireSet = wire.NewSet(
	producer.WireSet,
	consumer.WireSet,
	admin.WireSet,
)
