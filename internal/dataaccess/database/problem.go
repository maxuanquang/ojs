package database

import (
	"context"
	"errors"

	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrProblemNotFound      = errors.New("problem not found")
	ErrProblemAlreadyExists = errors.New("problem already exists")
)

type Problem struct {
	ID          uint64 `gorm:"column:id;primaryKey"`
	DisplayName string `gorm:"column:display_name"`
	AuthorID    uint64 `gorm:"column:author_id"`
	Description string `gorm:"column:description"`
	TimeLimit   uint64 `gorm:"column:time_limit"`
	MemoryLimit uint64 `gorm:"column:memory_limit"`
}

type ProblemDataAccessor interface {
	CreateProblem(ctx context.Context, problem Problem) (Problem, error)
	GetProblemByID(ctx context.Context, id uint64) (Problem, error)
	GetProblemByName(ctx context.Context, name string) (Problem, error)
	GetProblemList(ctx context.Context, offset uint64, limit uint64) ([]Problem, error)
	GetProblemCount(ctx context.Context) (uint64, error)
	UpdateProblem(ctx context.Context, id uint64, problem Problem) (Problem, error)
	DeleteProblem(ctx context.Context, id uint64) error
	WithDatabaseTransaction(database Database) ProblemDataAccessor
}

func NewProblemDataAccessor(database Database, logger *zap.Logger) ProblemDataAccessor {
	return &problemDataAccessor{
		database: database,
		logger:   logger,
	}
}

type problemDataAccessor struct {
	database Database
	logger   *zap.Logger
}

// GetProblemCount implements ProblemDataAccessor.
func (p *problemDataAccessor) GetProblemCount(ctx context.Context) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, p.logger).With(zap.String("method", "GetProblemList"))

	var problems []Problem
	result := p.database.Find(&problems)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		logger.Error("error getting problem list", zap.Error(result.Error))
		return 0, result.Error
	}

	return uint64(len(problems)), nil
}

// GetProblemList implements ProblemDataAccessor.
func (p *problemDataAccessor) GetProblemList(ctx context.Context, offset uint64, limit uint64) ([]Problem, error) {
	logger := utils.LoggerWithContext(ctx, p.logger).With(zap.String("method", "GetProblemList"))

	var problems []Problem
	result := p.database.Offset(int(offset)).Limit(int(limit)).Find(&problems)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.Error("error getting problem list", zap.Error(result.Error))
		return nil, result.Error
	}

	return problems, nil
}

// CreateProblem implements ProblemDataAccessor.
func (p *problemDataAccessor) CreateProblem(ctx context.Context, problem Problem) (Problem, error) {
	createdProblem := Problem{
		DisplayName: problem.DisplayName,
		AuthorID:    problem.AuthorID,
		Description: problem.Description,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
	}
	result := p.database.Create(&createdProblem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return Problem{}, ErrProblemAlreadyExists
		}

		logger := utils.LoggerWithContext(ctx, p.logger).With(zap.String("display_name", problem.DisplayName))
		logger.Error("error creating problem", zap.Error(result.Error))
		return Problem{}, result.Error
	}

	return createdProblem, nil
}

// GetProblemByID implements ProblemDataAccessor.
func (p *problemDataAccessor) GetProblemByID(ctx context.Context, id uint64) (Problem, error) {
	var foundProblem Problem
	result := p.database.First(&foundProblem, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Problem{}, ErrProblemNotFound
		}

		logger := utils.LoggerWithContext(ctx, p.logger).With(zap.Uint64("problem_id", id))
		logger.Error("error getting problem", zap.Error(result.Error))
		return Problem{}, result.Error
	}

	return foundProblem, nil
}

// GetProblemByName implements ProblemDataAccessor.
func (p *problemDataAccessor) GetProblemByName(ctx context.Context, name string) (Problem, error) {
	var foundProblem Problem
	result := p.database.Where("display_name = ?", name).First(&foundProblem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Problem{}, ErrProblemNotFound
		}

		logger := utils.LoggerWithContext(ctx, p.logger).With(zap.String("display_name", name))
		logger.Error("error getting problem by name", zap.Error(result.Error))
		return Problem{}, result.Error
	}

	return foundProblem, nil
}

// UpdateProblem implements ProblemDataAccessor.
func (p *problemDataAccessor) UpdateProblem(ctx context.Context, id uint64, newProblem Problem) (Problem, error) {
	var foundProblem Problem
	result := p.database.First(&foundProblem, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Problem{}, ErrProblemNotFound
		}

		logger := utils.LoggerWithContext(ctx, p.logger).With(zap.Uint64("problem_id", id))
		logger.Error("error finding problem to update", zap.Error(result.Error))
		return Problem{}, result.Error
	}

	if newProblem.DisplayName != "" {
		foundProblem.DisplayName = newProblem.DisplayName
	}
	if newProblem.Description != "" {
		foundProblem.Description = newProblem.Description
	}
	if newProblem.TimeLimit != 0 {
		foundProblem.TimeLimit = newProblem.TimeLimit
	}
	if newProblem.MemoryLimit != 0 {
		foundProblem.MemoryLimit = newProblem.MemoryLimit
	}

	result = p.database.Save(&foundProblem)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, p.logger).With(zap.Uint64("problem_id", id))
		logger.Error("error updating problem", zap.Error(result.Error))
		return Problem{}, result.Error
	}

	return foundProblem, nil
}

// DeleteProblem implements ProblemDataAccessor.
func (p *problemDataAccessor) DeleteProblem(ctx context.Context, id uint64) error {
	var foundProblem Problem
	result := p.database.First(&foundProblem, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrProblemNotFound
		}

		logger := utils.LoggerWithContext(ctx, p.logger).With(zap.Uint64("problem_id", id))
		logger.Error("error finding problem to delete", zap.Error(result.Error))
		return result.Error
	}

	result = p.database.Delete(&foundProblem, id)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, p.logger).With(zap.Uint64("problem_id", id))
		logger.Error("error deleting problem", zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// WithDatabaseTransaction implements ProblemDataAccessor.
func (p *problemDataAccessor) WithDatabaseTransaction(database Database) ProblemDataAccessor {
	return &problemDataAccessor{
		database: database,
		logger:   p.logger,
	}
}
