package logic

import (
	"context"

	"github.com/maxuanquang/ojs/internal/dataaccess/database"
	"go.uber.org/zap"
)

type TestCaseLogic interface {
	CreateTestCase(ctx context.Context, in CreateTestCaseInput) (CreateTestCaseOutput, error)
	GetTestCase(ctx context.Context, in GetTestCaseInput) (GetTestCaseOutput, error)
	GetProblemTestCaseList(ctx context.Context, in GetProblemTestCaseListInput) (GetProblemTestCaseListOutput, error)
	UpdateTestCase(ctx context.Context, in UpdateTestCaseInput) (UpdateTestCaseOutput, error)
	DeleteTestCase(ctx context.Context, in DeleteTestCaseInput) error
}

func NewTestCaseLogic(
	logger *zap.Logger,
	accountDataAccessor database.AccountDataAccessor,
	problemDataAccessor database.ProblemDataAccessor,
	submissionDataAccessor database.SubmissionDataAccessor,
	testCaseDataAccessor database.TestCaseDataAccessor,
	tokenLogic TokenLogic,
) TestCaseLogic {
	return &testCaseLogic{
		logger:                 logger,
		accountDataAccessor:    accountDataAccessor,
		problemDataAccessor:    problemDataAccessor,
		submissionDataAccessor: submissionDataAccessor,
		testCaseDataAccessor:   testCaseDataAccessor,
		tokenLogic:             tokenLogic,
	}
}

type testCaseLogic struct {
	logger                 *zap.Logger
	accountDataAccessor    database.AccountDataAccessor
	problemDataAccessor    database.ProblemDataAccessor
	submissionDataAccessor database.SubmissionDataAccessor
	testCaseDataAccessor   database.TestCaseDataAccessor
	tokenLogic             TokenLogic
}

func (t *testCaseLogic) CreateTestCase(ctx context.Context, in CreateTestCaseInput) (CreateTestCaseOutput, error) {
	logger := t.logger.With(zap.Any("create_test_case_input", in))

	// Create the test case in the database
	createdTestCase, err := t.testCaseDataAccessor.CreateTestCase(ctx, database.TestCase{
		OfProblemID: in.OfProblemID,
		Input:       in.Input,
		Output:      in.Output,
		IsHidden:    in.IsHidden,
	})
	if err != nil {
		logger.Error("failed to create test case", zap.Error(err))
		return CreateTestCaseOutput{}, err
	}

	return CreateTestCaseOutput{
		TestCase: TestCase{
			ID:          createdTestCase.ID,
			OfProblemID: createdTestCase.OfProblemID,
			Input:       createdTestCase.Input,
			Output:      createdTestCase.Output,
			IsHidden:    createdTestCase.IsHidden,
		},
	}, nil
}

func (t *testCaseLogic) GetTestCase(ctx context.Context, in GetTestCaseInput) (GetTestCaseOutput, error) {
	logger := t.logger.With(zap.Any("get_test_case_input", in))

	// Retrieve the test case from the database
	testCase, err := t.testCaseDataAccessor.GetTestCaseByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get test case", zap.Error(err))
		return GetTestCaseOutput{}, err
	}
	if testCase.ID == 0 {
		err := ErrTestCaseNotFound
		logger.Error("test case not found", zap.Error(err))
		return GetTestCaseOutput{}, err
	}

	return GetTestCaseOutput{
		TestCase: TestCase{
			ID:          testCase.ID,
			OfProblemID: testCase.OfProblemID,
			Input:       testCase.Input,
			Output:      testCase.Output,
			IsHidden:    testCase.IsHidden,
		},
	}, nil
}

func (t *testCaseLogic) GetProblemTestCaseList(ctx context.Context, in GetProblemTestCaseListInput) (GetProblemTestCaseListOutput, error) {
	logger := t.logger.With(zap.String("method", "GetProblemTestCaseList"))

	// Retrieve the list of test cases for a problem from the database
	testCases, err := t.testCaseDataAccessor.GetProblemTestCaseList(ctx, in.OfProblemID, in.Offset, in.Limit)
	if err != nil {
		logger.Error("failed to get test case list", zap.Error(err))
		return GetProblemTestCaseListOutput{}, err
	}

	var testCaseList []TestCase
	for _, tc := range testCases {
		testCaseList = append(testCaseList, TestCase{
			ID:          tc.ID,
			OfProblemID: tc.OfProblemID,
			Input:       tc.Input,
			Output:      tc.Output,
		})
	}

	totalTestCasesCount, err := t.testCaseDataAccessor.GetProblemTestCaseCount(ctx, in.OfProblemID)
	if err != nil {
		logger.Error("failed to get test case count", zap.Error(err))
		return GetProblemTestCaseListOutput{}, err
	}

	return GetProblemTestCaseListOutput{
		TestCases:           testCaseList,
		TotalTestCasesCount: totalTestCasesCount,
	}, nil
}

func (t *testCaseLogic) UpdateTestCase(ctx context.Context, in UpdateTestCaseInput) (UpdateTestCaseOutput, error) {
	logger := t.logger.With(zap.String("method", "UpdateTestCase"))

	// Check if the test case exists
	testCase, err := t.testCaseDataAccessor.GetTestCaseByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get test case", zap.Error(err))
		return UpdateTestCaseOutput{}, err
	}
	if testCase.ID == 0 {
		err := ErrTestCaseNotFound
		logger.Error("test case not found", zap.Error(err))
		return UpdateTestCaseOutput{}, err
	}

	// Update the test case in the database
	updatedTestCase, err := t.testCaseDataAccessor.UpdateTestCase(ctx, database.TestCase{
		ID:       in.ID,
		Input:    in.Input,
		Output:   in.Output,
		IsHidden: in.IsHidden,
	})
	if err != nil {
		logger.Error("failed to update test case", zap.Error(err))
		return UpdateTestCaseOutput{}, err
	}

	return UpdateTestCaseOutput{
		TestCase: TestCase{
			ID:          updatedTestCase.ID,
			OfProblemID: updatedTestCase.OfProblemID,
			Input:       updatedTestCase.Input,
			Output:      updatedTestCase.Output,
			IsHidden:    updatedTestCase.IsHidden,
		},
	}, nil
}

func (t *testCaseLogic) DeleteTestCase(ctx context.Context, in DeleteTestCaseInput) error {
	logger := t.logger.With(zap.String("method", "DeleteTestCase"))

	// Check if the test case exists
	testCase, err := t.testCaseDataAccessor.GetTestCaseByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get test case", zap.Error(err))
		return err
	}
	if testCase.ID == 0 {
		err := ErrTestCaseNotFound
		logger.Error("test case not found", zap.Error(err))
		return err
	}

	// Delete the test case from the database
	err = t.testCaseDataAccessor.DeleteTestCase(ctx, in.ID)
	if err != nil {
		logger.Error("failed to delete test case", zap.Error(err))
		return err
	}

	return nil
}

type TestCase struct {
	ID          uint64
	OfProblemID uint64
	Input       string
	Output      string
	IsHidden    bool
}
type CreateTestCaseInput struct {
	OfProblemID uint64
	Input       string
	Output      string
	IsHidden    bool
}

type CreateTestCaseOutput struct {
	TestCase TestCase
}

type GetTestCaseInput struct {
	ID uint64
}

type GetTestCaseOutput struct {
	TestCase TestCase
}

type GetProblemTestCaseListInput struct {
	OfProblemID uint64
	Offset      uint64
	Limit       uint64
}

type GetProblemTestCaseListOutput struct {
	TestCases           []TestCase
	TotalTestCasesCount uint64
}

type UpdateTestCaseInput struct {
	ID       uint64
	Input    string
	Output   string
	IsHidden bool
}

type UpdateTestCaseOutput struct {
	TestCase TestCase
}

type DeleteTestCaseInput struct {
	ID uint64
}

type DeleteTestCaseOutput struct{}
