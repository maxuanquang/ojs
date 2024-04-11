package handler

import (
	"github.com/google/wire"
	"github.com/maxuanquang/ojs/internal/handler/consumer"
	"github.com/maxuanquang/ojs/internal/handler/grpc"
	"github.com/maxuanquang/ojs/internal/handler/http"
	"github.com/maxuanquang/ojs/internal/handler/jobs"
)

var WireSet = wire.NewSet(
	grpc.WireSet,
	http.WireSet,
	consumer.WireSet,
	jobs.WireSet,
)
