//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"github.com/google/wire"
	"github.com/maxuanquang/ojs/internal/app"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/dataaccess"
	"github.com/maxuanquang/ojs/internal/handler"
	"github.com/maxuanquang/ojs/internal/logic"
	"github.com/maxuanquang/ojs/internal/utils"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	dataaccess.WireSet,
	handler.WireSet,
	logic.WireSet,
	utils.WireSet,
	app.WireSet,
)

func InitializeStandaloneServer(configFilePath configs.ConfigFilePath, appArguments utils.Arguments) (app.StandaloneServer, func(), error) {
	wire.Build(WireSet)

	return app.StandaloneServer{}, nil, nil
}

func InitializeHTTPServer(configFilePath configs.ConfigFilePath, appArguments utils.Arguments) (app.HTTPServer, func(), error) {
	wire.Build(WireSet)

	return app.HTTPServer{}, nil, nil
}

func InitializeWorker(configFilePath configs.ConfigFilePath, appArguments utils.Arguments) (app.Worker, func(), error) {
	wire.Build(WireSet)

	return app.Worker{}, nil, nil
}

func InitializeCron(configFilePath configs.ConfigFilePath, appArguments utils.Arguments) (app.Cron, func(), error) {
	wire.Build(WireSet)

	return app.Cron{}, nil, nil
}
