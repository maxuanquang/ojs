package logic

import (
	"context"

	"github.com/maxuanquang/ojs/internal/dataaccess/database"
	"github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrProblemNotFound  = status.Error(codes.NotFound, "problem not found")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "permission denied")
	ErrTestCaseNotFound = status.Error(codes.NotFound, "test case not found")
)

type ProblemLogic interface {
	CreateProblem(ctx context.Context, in CreateProblemInput) (CreateProblemOutput, error)
	GetProblem(ctx context.Context, in GetProblemInput) (GetProblemOutput, error)
	GetProblemList(ctx context.Context, in GetProblemListInput) (GetProblemListOutput, error)
	UpdateProblem(ctx context.Context, in UpdateProblemInput) (UpdateProblemOutput, error)
	DeleteProblem(ctx context.Context, in DeleteProblemInput) error

	CreateTestCase(ctx context.Context, in CreateTestCaseInput) (CreateTestCaseOutput, error)
	GetTestCase(ctx context.Context, in GetTestCaseInput) (GetTestCaseOutput, error)
	GetProblemTestCaseList(ctx context.Context, in GetProblemTestCaseListInput) (GetProblemTestCaseListOutput, error)
	UpdateTestCase(ctx context.Context, in UpdateTestCaseInput) (UpdateTestCaseOutput, error)
	DeleteTestCase(ctx context.Context, in DeleteTestCaseInput) error
}

func NewProblemLogic(
	logger *zap.Logger,
	accountDataAccessor database.AccountDataAccessor,
	problemDataAccessor database.ProblemDataAccessor,
	submissionDataAccessor database.SubmissionDataAccessor,
	testCaseDataAccessor database.TestCaseDataAccessor,
	tokenLogic TokenLogic,
) ProblemLogic {
	return &problemLogic{
		logger:                 logger,
		accountDataAccessor:    accountDataAccessor,
		problemDataAccessor:    problemDataAccessor,
		submissionDataAccessor: submissionDataAccessor,
		testCaseDataAccessor:   testCaseDataAccessor,
		tokenLogic:             tokenLogic,
	}
}

type problemLogic struct {
	logger                 *zap.Logger
	accountDataAccessor    database.AccountDataAccessor
	problemDataAccessor    database.ProblemDataAccessor
	submissionDataAccessor database.SubmissionDataAccessor
	testCaseDataAccessor   database.TestCaseDataAccessor
	tokenLogic             TokenLogic
}

func (p *problemLogic) CreateProblem(ctx context.Context, in CreateProblemInput) (CreateProblemOutput, error) {
	logger := p.logger.With(zap.Any("create_problem_input", in))

	account_id, _, _, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return CreateProblemOutput{}, err
	}

	// Create the problem in the database
	createdProblem, err := p.problemDataAccessor.CreateProblem(ctx, database.Problem{
		DisplayName: in.DisplayName,
		Description: in.Description,
		TimeLimit:   in.TimeLimit,
		MemoryLimit: in.MemoryLimit,
		AuthorID:    account_id,
	})
	if err != nil {
		logger.Error("failed to create problem", zap.Error(err))
		return CreateProblemOutput{}, err
	}

	return CreateProblemOutput{
		Problem: Problem{
			ID:          createdProblem.ID,
			DisplayName: createdProblem.DisplayName,
			AuthorId:    createdProblem.AuthorID,
			Description: createdProblem.Description,
			TimeLimit:   createdProblem.TimeLimit,
			MemoryLimit: createdProblem.MemoryLimit,
		},
	}, nil
}

func (p *problemLogic) GetProblem(ctx context.Context, in GetProblemInput) (GetProblemOutput, error) {
	logger := p.logger.With(zap.Any("get_problem_input", in))

	// Retrieve the problem from the database
	problem, err := p.problemDataAccessor.GetProblemByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return GetProblemOutput{}, err
	}
	if problem.ID == 0 {
		err := ErrProblemNotFound
		logger.Error("problem not found", zap.Error(err))
		return GetProblemOutput{}, err
	}

	return GetProblemOutput{
		Problem: Problem{
			ID:          problem.ID,
			DisplayName: problem.DisplayName,
			AuthorId:    problem.AuthorID,
			Description: problem.Description,
			TimeLimit:   problem.TimeLimit,
			MemoryLimit: problem.MemoryLimit,
		},
	}, nil
}

func (p *problemLogic) GetProblemList(ctx context.Context, in GetProblemListInput) (GetProblemListOutput, error) {
	logger := p.logger.With(zap.String("method", "GetProblemList"))

	// Retrieve the list of problems from the database
	problems, err := p.problemDataAccessor.GetProblemList(ctx, in.Offset, in.Limit)
	if err != nil {
		logger.Error("failed to get problem list", zap.Error(err))
		return GetProblemListOutput{}, err
	}
	total, err := p.problemDataAccessor.GetProblemCount(ctx)
	if err != nil {
		logger.Error("failed to get problem count", zap.Error(err))
		return GetProblemListOutput{}, err
	}

	var problemList []Problem
	for _, pb := range problems {
		problemList = append(problemList, Problem{
			ID:          pb.ID,
			DisplayName: pb.DisplayName,
			AuthorId:    pb.AuthorID,
			Description: pb.Description,
			TimeLimit:   pb.TimeLimit,
			MemoryLimit: pb.MemoryLimit,
		})
	}

	return GetProblemListOutput{
		Problems:          problemList,
		TotalProblemCount: total,
	}, nil
}

func (p *problemLogic) UpdateProblem(ctx context.Context, in UpdateProblemInput) (UpdateProblemOutput, error) {
	logger := p.logger.With(zap.String("method", "UpdateProblem"))

	accountID, accountName, accountRole, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return UpdateProblemOutput{}, err
	}
	if ojs.Role(accountRole) != ojs.Role_Admin {
		return UpdateProblemOutput{}, ErrPermissionDenied
	}

	// Check if the problem exists
	problem, err := p.problemDataAccessor.GetProblemByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return UpdateProblemOutput{}, err
	}
	if problem.ID == 0 {
		err := ErrProblemNotFound
		logger.Error("problem not found", zap.Error(err))
		return UpdateProblemOutput{}, err
	}
	if problem.AuthorID != accountID {
		return UpdateProblemOutput{}, ErrPermissionDenied
	}

	// Update the problem in the database
	updatedProblem, err := p.problemDataAccessor.UpdateProblem(ctx, in.ID, database.Problem{
		DisplayName: in.DisplayName,
		Description: in.Description,
		TimeLimit:   in.TimeLimit,
		MemoryLimit: in.MemoryLimit,
	})
	if err != nil {
		logger.Error("failed to update problem", zap.Error(err))
		return UpdateProblemOutput{}, err
	}

	return UpdateProblemOutput{
		Problem: Problem{
			ID:          updatedProblem.ID,
			DisplayName: updatedProblem.DisplayName,
			AuthorId:    updatedProblem.AuthorID,
			AuthorName:  accountName,
			Description: updatedProblem.Description,
			TimeLimit:   updatedProblem.TimeLimit,
			MemoryLimit: updatedProblem.MemoryLimit,
		},
	}, nil
}

func (p *problemLogic) DeleteProblem(ctx context.Context, in DeleteProblemInput) error {
	logger := p.logger.With(zap.String("method", "DeleteProblem"))

	// Check if the problem exists
	problem, err := p.problemDataAccessor.GetProblemByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return err
	}
	if problem.ID == 0 {
		err := ErrProblemNotFound
		logger.Error("problem not found", zap.Error(err))
		return err
	}

	// Delete the problem from the database
	err = p.problemDataAccessor.DeleteProblem(ctx, in.ID)
	if err != nil {
		logger.Error("failed to delete problem", zap.Error(err))
		return err
	}

	return nil
}

func (t *problemLogic) CreateTestCase(ctx context.Context, in CreateTestCaseInput) (CreateTestCaseOutput, error) {
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

func (t *problemLogic) GetTestCase(ctx context.Context, in GetTestCaseInput) (GetTestCaseOutput, error) {
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

func (t *problemLogic) GetProblemTestCaseList(ctx context.Context, in GetProblemTestCaseListInput) (GetProblemTestCaseListOutput, error) {
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

func (t *problemLogic) UpdateTestCase(ctx context.Context, in UpdateTestCaseInput) (UpdateTestCaseOutput, error) {
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

func (t *problemLogic) DeleteTestCase(ctx context.Context, in DeleteTestCaseInput) error {
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

type CreateProblemInput struct {
	Token       string
	DisplayName string
	Description string
	TimeLimit   uint64
	MemoryLimit uint64
}

type CreateProblemOutput struct {
	Problem Problem
}

type Problem struct {
	ID          uint64
	DisplayName string
	AuthorId    uint64
	AuthorName  string
	Description string
	TimeLimit   uint64
	MemoryLimit uint64
}

type GetProblemListInput struct {
	Token  string
	Offset uint64
	Limit  uint64
}

type GetProblemListOutput struct {
	Problems          []Problem
	TotalProblemCount uint64
}

type GetProblemInput struct {
	Token string
	ID    uint64
}

type GetProblemOutput struct {
	Problem Problem
}

type UpdateProblemInput struct {
	Token       string
	ID          uint64
	DisplayName string
	Description string
	TimeLimit   uint64
	MemoryLimit uint64
}

type UpdateProblemOutput struct {
	Problem Problem
}

type DeleteProblemInput struct {
	Token string
	ID    uint64
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
