package database

import (
	"context"
	"errors"

	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrTestCaseNotFound = errors.New("test case not found")
)

type TestCase struct {
	ID           uint64 `gorm:"column:id;primaryKey"`
	ProblemID    uint64 `gorm:"column:of_problem_id"`
	Input        string `gorm:"column:input"`
	Output       string `gorm:"column:output"`
	IsHidden     bool   `gorm:"column:is_hidden"`
}

type TestCaseDataAccessor interface {
	CreateTestCase(ctx context.Context, testCase TestCase) (TestCase, error)
	GetTestCaseByID(ctx context.Context, id uint64) (TestCase, error)
	DeleteTestCase(ctx context.Context, id uint64) error
	GetTestCaseOfProblem(ctx context.Context, problemID uint64) ([]TestCase, error)
	UpdateTestCase(ctx context.Context, testCase TestCase) (TestCase, error)
	WithDatabaseTransaction(database Database) TestCaseDataAccessor
}



func NewTestCaseDataAccessor(database Database, logger *zap.Logger) TestCaseDataAccessor {
	return &testCaseDataAccessor{
		database: database,
		logger:   logger,
	}
}

type testCaseDataAccessor struct {
	database Database
	logger   *zap.Logger
}

// CreateTestCase implements TestCaseDataAccessor.
func (t *testCaseDataAccessor) CreateTestCase(ctx context.Context, testCase TestCase) (TestCase, error) {
	createdTestCase := TestCase{
		ProblemID: testCase.ProblemID,
		Input:     testCase.Input,
		Output:    testCase.Output,
		IsHidden:  testCase.IsHidden,
	}
	result := t.database.Create(&createdTestCase)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("problem_id", testCase.ProblemID))
		logger.Error("error creating test case", zap.Error(result.Error))
		return TestCase{}, result.Error
	}

	return createdTestCase, nil
}

// GetTestCaseByID implements TestCaseDataAccessor.
func (t *testCaseDataAccessor) GetTestCaseByID(ctx context.Context, id uint64) (TestCase, error) {
	var foundTestCase TestCase
	result := t.database.First(&foundTestCase, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return TestCase{}, ErrTestCaseNotFound
		}

		logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("test_case_id", id))
		logger.Error("error getting test case", zap.Error(result.Error))
		return TestCase{}, result.Error
	}

	return foundTestCase, nil
}

// DeleteTestCase implements TestCaseDataAccessor.
func (t *testCaseDataAccessor) DeleteTestCase(ctx context.Context, id uint64) error {
	var foundTestCase TestCase
	result := t.database.First(&foundTestCase, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrTestCaseNotFound
		}

		logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("test_case_id", id))
		logger.Error("error finding test case to delete", zap.Error(result.Error))
		return result.Error
	}

	result = t.database.Delete(&foundTestCase, id)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("test_case_id", id))
		logger.Error("error deleting test case", zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// UpdateTestCase implements TestCaseDataAccessor.
func (t *testCaseDataAccessor) UpdateTestCase(ctx context.Context, testCase TestCase) (TestCase, error) {
	var existingTestCase TestCase
	result := t.database.First(&existingTestCase, testCase.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return TestCase{}, ErrTestCaseNotFound
		}
		logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("test_case_id", testCase.ID))
		logger.Error("error updating test case", zap.Error(result.Error))
		return TestCase{}, result.Error
	}

	existingTestCase.Input = testCase.Input
	existingTestCase.Output = testCase.Output
	existingTestCase.IsHidden = testCase.IsHidden

	result = t.database.Save(&existingTestCase)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("test_case_id", testCase.ID))
		logger.Error("error updating test case", zap.Error(result.Error))
		return TestCase{}, result.Error
	}

	return existingTestCase, nil
}

// GetTestCaseOfProblem implements TestCaseDataAccessor.
func (t *testCaseDataAccessor) GetTestCaseOfProblem(ctx context.Context, problemID uint64) ([]TestCase, error) {
	var testCases []TestCase
	result := t.database.Where("of_problem_id = ?", problemID).Find(&testCases)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, t.logger).With(zap.Uint64("problem_id", problemID))
		logger.Error("error getting test cases of problem", zap.Error(result.Error))
		return nil, result.Error
	}

	return testCases, nil
}

// WithDatabaseTransaction implements TestCaseDataAccessor.
func (t *testCaseDataAccessor) WithDatabaseTransaction(database Database) TestCaseDataAccessor {
	return &testCaseDataAccessor{
		database: database,
		logger:   t.logger,
	}
}
