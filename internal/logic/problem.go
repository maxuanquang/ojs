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

