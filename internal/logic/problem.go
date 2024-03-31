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

	CreateSubmission(ctx context.Context, in CreateSubmissionInput) (CreateSubmissionOutput, error)
	GetSubmission(ctx context.Context, in GetSubmissionInput) (GetSubmissionOutput, error)
	GetSubmissionList(ctx context.Context, in GetSubmissionListInput) (GetSubmissionListOutput, error)
	GetAccountProblemSubmissionList(ctx context.Context, in GetAccountProblemSubmissionListInput) (GetAccountProblemSubmissionListOutput, error)
	GetProblemSubmissionList(ctx context.Context, in GetProblemSubmissionListInput) (GetProblemSubmissionListOutput, error)
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

func (p *problemLogic) CreateSubmission(ctx context.Context, in CreateSubmissionInput) (CreateSubmissionOutput, error) {
	// Verify token
	accountID, _, _, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		p.logger.Error("Failed to verify token", zap.Error(err))
		return CreateSubmissionOutput{}, err
	}

	// Create submission in the database
	createdSubmission, err := p.submissionDataAccessor.CreateSubmission(ctx, database.Submission{
		OfProblemID: in.OfProblemID,
		AuthorID:    accountID,
		Content:     in.Content,
		Language:    in.Language,
	})
	if err != nil {
		p.logger.Error("Failed to create submission", zap.Error(err))
		return CreateSubmissionOutput{}, err
	}

	return CreateSubmissionOutput{
		Submission: Submission{
			ID:          createdSubmission.ID,
			OfProblemID: createdSubmission.OfProblemID,
			AuthorID:    createdSubmission.AuthorID,
			Content:     createdSubmission.Content,
			Language:    createdSubmission.Language,
			Status:      ojs.SubmissionStatus(createdSubmission.Status),
			Result:      ojs.SubmissionResult(createdSubmission.Result),
		},
	}, nil
}

func (p *problemLogic) GetSubmission(ctx context.Context, in GetSubmissionInput) (GetSubmissionOutput, error) {
	// Retrieve submission from the database
	submission, err := p.submissionDataAccessor.GetSubmissionByID(ctx, in.ID)
	if err != nil {
		p.logger.Error("Failed to get submission", zap.Error(err))
		return GetSubmissionOutput{}, err
	}

	return GetSubmissionOutput{
		Submission: Submission{
			ID:          submission.ID,
			OfProblemID: submission.OfProblemID,
			AuthorID:    submission.AuthorID,
			Content:     submission.Content,
			Language:    submission.Language,
			Status:      ojs.SubmissionStatus(submission.Status),
			Result:      ojs.SubmissionResult(submission.Result),
		},
	}, nil
}

func (p *problemLogic) GetSubmissionList(ctx context.Context, in GetSubmissionListInput) (GetSubmissionListOutput, error) {
	// Retrieve submission list from the database
	submissions, err := p.submissionDataAccessor.GetSubmissionList(ctx, in.Offset, in.Limit)
	if err != nil {
		p.logger.Error("Failed to get submission list", zap.Error(err))
		return GetSubmissionListOutput{}, err
	}

	var submissionList []Submission
	for _, s := range submissions {
		submissionList = append(submissionList, Submission{
			ID:          s.ID,
			OfProblemID: s.OfProblemID,
			AuthorID:    s.AuthorID,
			Content:     s.Content,
			Language:    s.Language,
			Status:      ojs.SubmissionStatus(s.Status),
			Result:      ojs.SubmissionResult(s.Result),
		})
	}

	totalSubmissionsCount, err := p.submissionDataAccessor.GetSubmissionCount(ctx)
	if err != nil {
		p.logger.Error("Failed to get submission count", zap.Error(err))
		return GetSubmissionListOutput{}, err
	}

	return GetSubmissionListOutput{
		Submissions:           submissionList,
		TotalSubmissionsCount: totalSubmissionsCount,
	}, nil
}

func (p *problemLogic) GetAccountProblemSubmissionList(ctx context.Context, in GetAccountProblemSubmissionListInput) (GetAccountProblemSubmissionListOutput, error) {
	// Verify token
	accountID, _, _, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		p.logger.Error("Failed to verify token", zap.Error(err))
		return GetAccountProblemSubmissionListOutput{}, err
	}

	// Retrieve account's submission list for a specific problem from the database
	submissions, err := p.submissionDataAccessor.GetAccountProblemSubmissionList(ctx, accountID, in.OfProblemID, in.Offset, in.Limit)
	if err != nil {
		p.logger.Error("Failed to get account's problem submission list", zap.Error(err))
		return GetAccountProblemSubmissionListOutput{}, err
	}

	var submissionList []Submission
	for _, s := range submissions {
		submissionList = append(submissionList, Submission{
			ID:          s.ID,
			OfProblemID: s.OfProblemID,
			AuthorID:    s.AuthorID,
			Content:     s.Content,
			Language:    s.Language,
			Status:      ojs.SubmissionStatus(s.Status),
			Result:      ojs.SubmissionResult(s.Result),
		})
	}

	totalSubmissionsCount, err := p.submissionDataAccessor.GetAccountProblemSubmissionCount(ctx, accountID, in.OfProblemID)
	if err != nil {
		p.logger.Error("Failed to get account's problem submission count", zap.Error(err))
		return GetAccountProblemSubmissionListOutput{}, err
	}

	return GetAccountProblemSubmissionListOutput{
		Submissions:           submissionList,
		TotalSubmissionsCount: totalSubmissionsCount,
	}, nil
}

func (p *problemLogic) GetProblemSubmissionList(ctx context.Context, in GetProblemSubmissionListInput) (GetProblemSubmissionListOutput, error) {
	// Retrieve problem's submission list from the database
	submissions, err := p.submissionDataAccessor.GetProblemSubmissionList(ctx, in.OfProblemID, in.Offset, in.Limit)
	if err != nil {
		p.logger.Error("Failed to get problem submission list", zap.Error(err))
		return GetProblemSubmissionListOutput{}, err
	}

	var submissionList []Submission
	for _, s := range submissions {
		submissionList = append(submissionList, Submission{
			ID:          s.ID,
			OfProblemID: s.OfProblemID,
			AuthorID:    s.AuthorID,
			Content:     s.Content,
			Language:    s.Language,
			Status:      ojs.SubmissionStatus(s.Status),
			Result:      ojs.SubmissionResult(s.Result),
		})
	}

	totalSubmissionsCount, err := p.submissionDataAccessor.GetProblemSubmissionCount(ctx, in.OfProblemID)
	if err != nil {
		p.logger.Error("Failed to get problem submission count", zap.Error(err))
		return GetProblemSubmissionListOutput{}, err
	}

	return GetProblemSubmissionListOutput{
		Submissions:           submissionList,
		TotalSubmissionsCount: totalSubmissionsCount,
	}, nil
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

type CreateSubmissionInput struct {
	Token       string
	OfProblemID uint64
	Content     string
	Language    string
}

type Submission struct {
	ID          uint64
	AuthorID    uint64
	OfProblemID uint64
	Content     string
	Language    string
	Status      ojs.SubmissionStatus
	Result      ojs.SubmissionResult
}

type CreateSubmissionOutput struct {
	Submission Submission
}

type GetSubmissionInput struct {
	ID uint64
}

type GetSubmissionOutput struct {
	Submission Submission
}

type GetSubmissionListInput struct {
	Offset uint64
	Limit  uint64
}

type GetSubmissionListOutput struct {
	Submissions           []Submission
	TotalSubmissionsCount uint64
}

type GetAccountProblemSubmissionListInput struct {
	Token       string
	OfProblemID uint64
	Offset      uint64
	Limit       uint64
}

type GetAccountProblemSubmissionListOutput struct {
	Submissions           []Submission
	TotalSubmissionsCount uint64
}

type GetProblemSubmissionListInput struct {
	Token       string
	OfProblemID uint64
	Offset      uint64
	Limit       uint64
	IsHidden    bool
}

type GetProblemSubmissionListOutput struct {
	Submissions           []Submission
	TotalSubmissionsCount uint64
}
