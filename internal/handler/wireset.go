package handler

import (
	"github.com/google/wire"
	"github.com/maxuanquang/ojs/internal/handler/grpc"
	"github.com/maxuanquang/ojs/internal/handler/http"
)

var WireSet = wire.NewSet(
	grpc.WireSet,
	http.WireSet,
)
