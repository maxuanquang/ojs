//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"github.com/google/wire"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/dataaccess"
	"github.com/maxuanquang/ojs/internal/handler"
	"github.com/maxuanquang/ojs/internal/logic"
	"github.com/maxuanquang/ojs/internal/utils"
	"github.com/maxuanquang/ojs/internal/app"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	dataaccess.WireSet,
	handler.WireSet,
	logic.WireSet,
	utils.WireSet,
	app.WireSet,
)

func InitializeAppServer(configFilePath configs.ConfigFilePath, appArguments utils.Arguments) (app.Server, func(), error) {
	wire.Build(WireSet)

	return app.Server{}, nil, nil
}