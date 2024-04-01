// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/maxuanquang/ojs/internal/app"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/dataaccess"
	"github.com/maxuanquang/ojs/internal/dataaccess/cache"
	"github.com/maxuanquang/ojs/internal/dataaccess/database"
	consumer2 "github.com/maxuanquang/ojs/internal/dataaccess/mq/consumer"
	"github.com/maxuanquang/ojs/internal/dataaccess/mq/producer"
	"github.com/maxuanquang/ojs/internal/handler"
	"github.com/maxuanquang/ojs/internal/handler/consumer"
	"github.com/maxuanquang/ojs/internal/handler/grpc"
	"github.com/maxuanquang/ojs/internal/handler/http"
	"github.com/maxuanquang/ojs/internal/logic"
	"github.com/maxuanquang/ojs/internal/utils"
)

// Injectors from wire.go:

func InitializeAppServer(configFilePath configs.ConfigFilePath) (app.Server, func(), error) {
	config, err := configs.NewConfig(configFilePath)
	if err != nil {
		return app.Server{}, nil, err
	}
	configsGRPC := config.GRPC
	configsDatabase := config.Database
	databaseDatabase, cleanup, err := database.InitializeDB(configsDatabase)
	if err != nil {
		return app.Server{}, nil, err
	}
	log := config.Log
	logger, cleanup2, err := utils.InitializeLogger(log)
	if err != nil {
		cleanup()
		return app.Server{}, nil, err
	}
	accountDataAccessor := database.NewAccountDataAccessor(databaseDatabase, logger)
	accountPasswordDataAccessor := database.NewAccountPasswordDataAccessor(databaseDatabase, logger)
	auth := config.Auth
	hashLogic := logic.NewHashLogic(auth)
	tokenPublicKeyDataAccessor, err := database.NewTokenPublicKeyDataAccessor(databaseDatabase, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	configsCache := config.Cache
	client, err := cache.NewClient(configsCache, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	tokenPublicKey, err := cache.NewTokenPublicKey(client)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	tokenLogic, err := logic.NewTokenLogic(accountDataAccessor, tokenPublicKeyDataAccessor, logger, auth, tokenPublicKey)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	takenAccountName, err := cache.NewTakenAccountName(client)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	accountLogic := logic.NewAccountLogic(databaseDatabase, accountDataAccessor, accountPasswordDataAccessor, hashLogic, tokenLogic, takenAccountName, logger)
	problemDataAccessor := database.NewProblemDataAccessor(databaseDatabase, logger)
	submissionDataAccessor := database.NewSubmissionDataAccessor(databaseDatabase, logger)
	testCaseDataAccessor := database.NewTestCaseDataAccessor(databaseDatabase, logger)
	problemLogic := logic.NewProblemLogic(logger, accountDataAccessor, problemDataAccessor, submissionDataAccessor, testCaseDataAccessor, tokenLogic)
	mq := config.MQ
	producerClient, err := producer.NewClient(mq, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	submissionCreatedProducer, err := producer.NewSubmissionCreatedProducer(producerClient, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	submissionLogic := logic.NewSubmissionLogic(logger, accountDataAccessor, problemDataAccessor, submissionDataAccessor, testCaseDataAccessor, tokenLogic, submissionCreatedProducer, databaseDatabase)
	ojsServiceServer := grpc.NewHandler(accountLogic, problemLogic, submissionLogic)
	server := grpc.NewServer(configsGRPC, ojsServiceServer)
	configsHTTP := config.HTTP
	httpServer := http.NewServer(configsHTTP, configsGRPC, auth, logger)
	submissionCreatedHandler, err := consumer.NewSubmissionCreatedHandler(submissionLogic, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	consumerConsumer, err := consumer2.NewConsumer(mq, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	rootConsumer := consumer.NewRootConsumer(submissionCreatedHandler, consumerConsumer, logger)
	appServer, err := app.NewServer(server, httpServer, rootConsumer, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	return appServer, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var WireSet = wire.NewSet(configs.WireSet, dataaccess.WireSet, handler.WireSet, logic.WireSet, utils.WireSet, app.WireSet)
