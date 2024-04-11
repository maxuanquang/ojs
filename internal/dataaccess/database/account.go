package database

import (
	"context"
	"errors"

	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
)

type Account struct {
	ID   uint64 `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
	Role int8   `gotm:"column:role"`
}

type AccountDataAccessor interface {
	CreateAccount(ctx context.Context, account Account) (Account, error)
	GetAccountByID(ctx context.Context, id uint64) (Account, error)
	GetAccountByName(ctx context.Context, name string) (Account, error)
	WithDatabaseTransaction(database Database) AccountDataAccessor
}

func NewAccountDataAccessor(database Database, logger *zap.Logger) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   logger,
	}
}

type accountDataAccessor struct {
	database Database
	logger   *zap.Logger
}

// CreateAccount implements AccountDataAccessor.
func (a *accountDataAccessor) CreateAccount(ctx context.Context, account Account) (Account, error) {
	createdAccount := Account{
		Name: account.Name,
		Role: account.Role,
	}
	result := a.database.Create(&createdAccount)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return Account{}, ErrAccountAlreadyExists
		}

		logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("name", account.Name))
		logger.Error("error creating account", zap.Error(result.Error))
		return Account{}, result.Error
	}

	return createdAccount, nil
}

// GetAccountByID implements AccountDataAccessor.
func (a *accountDataAccessor) GetAccountByID(ctx context.Context, id uint64) (Account, error) {
	var foundAccount Account
	result := a.database.First(&foundAccount, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Account{}, nil
		}

		logger := utils.LoggerWithContext(ctx, a.logger).With(zap.Uint64("account_id", id))
		logger.Error("error getting account", zap.Error(result.Error))
		return Account{}, result.Error
	}

	return foundAccount, nil
}

// GetAccountByName implements AccountDataAccessor.
func (a *accountDataAccessor) GetAccountByName(ctx context.Context, name string) (Account, error) {
	var foundAccount Account
	result := a.database.Where("name = ?", name).First(&foundAccount)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Account{}, nil
		}

		logger := utils.LoggerWithContext(ctx, a.logger).With(zap.String("name", name))
		logger.Error("error getting account", zap.Error(result.Error))
		return Account{}, result.Error
	}

	return foundAccount, nil
}

// WithDatabaseTransaction implements AccountDataAccessor.
func (a *accountDataAccessor) WithDatabaseTransaction(database Database) AccountDataAccessor {
	return &accountDataAccessor{
		database: database,
		logger:   a.logger,
	}
}
