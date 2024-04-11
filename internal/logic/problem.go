package logic

import (
	"context"

	"github.com/maxuanquang/ojs/internal/dataaccess/database"
	"github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	"github.com/mikespook/gorbac"
	"go.uber.org/zap"
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
	roleLogic RoleLogic,
) ProblemLogic {
	return &problemLogic{
		logger:                 logger,
		accountDataAccessor:    accountDataAccessor,
		problemDataAccessor:    problemDataAccessor,
		submissionDataAccessor: submissionDataAccessor,
		testCaseDataAccessor:   testCaseDataAccessor,
		tokenLogic:             tokenLogic,
		roleLogic:              roleLogic,
	}
}

type problemLogic struct {
	logger                 *zap.Logger
	accountDataAccessor    database.AccountDataAccessor
	problemDataAccessor    database.ProblemDataAccessor
	submissionDataAccessor database.SubmissionDataAccessor
	testCaseDataAccessor   database.TestCaseDataAccessor
	roleLogic              RoleLogic
	tokenLogic             TokenLogic
}

func (p *problemLogic) CreateProblem(ctx context.Context, in CreateProblemInput) (CreateProblemOutput, error) {
	logger := p.logger.With(zap.Any("create_problem_input", in))

	requestingAccountID, requestingAccountName, requestingAccountRole, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return CreateProblemOutput{}, ErrTokenInvalid
	}

	requiredPermissions := []gorbac.Permission{PermissionProblemsWriteAll, PermissionProblemsWriteSelf}
	hasPermission, err := p.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountRole)], requiredPermissions...)
	if err != nil {
		logger.Error("failed to check permission", zap.Error(err))
		return CreateProblemOutput{}, ErrInternal
	}
	if !hasPermission {
		return CreateProblemOutput{}, ErrPermissionDenied
	}

	createdProblem, err := p.problemDataAccessor.CreateProblem(ctx, database.Problem{
		DisplayName: in.DisplayName,
		Description: in.Description,
		TimeLimit:   in.TimeLimit,
		MemoryLimit: in.MemoryLimit,
		AuthorID:    requestingAccountID,
	})
	if err != nil {
		logger.Error("failed to create problem", zap.Error(err))
		return CreateProblemOutput{}, ErrInternal
	}

	return CreateProblemOutput{
		Problem: p.dbProblemToLogicProblem(createdProblem, requestingAccountName),
	}, nil
}

func (p *problemLogic) GetProblem(ctx context.Context, in GetProblemInput) (GetProblemOutput, error) {
	logger := p.logger.With(zap.Any("get_problem_input", in))

	requestingAccountID, _, requestingAccountRole, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return GetProblemOutput{}, ErrTokenInvalid
	}

	requiredPermissions := []gorbac.Permission{PermissionProblemsReadAll}
	if requestingAccountID == in.ID {
		requiredPermissions = append(requiredPermissions, PermissionProblemsReadSelf)
	}

	hasPermission, err := p.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountRole)], requiredPermissions...)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check permission")
		return GetProblemOutput{}, ErrInternal
	}
	if !hasPermission {
		return GetProblemOutput{}, ErrPermissionDenied
	}

	problem, err := p.problemDataAccessor.GetProblemByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return GetProblemOutput{}, ErrInternal
	}
	if problem.ID == 0 {
		logger.Error("problem not found", zap.Error(err))
		return GetProblemOutput{}, ErrProblemNotFound
	}

	author, err := p.accountDataAccessor.GetAccountByID(ctx, problem.AuthorID)
	if err != nil {
		logger.Error("failed to get author", zap.Error(err))
		return GetProblemOutput{}, ErrInternal
	}
	if author.ID == 0 {
		logger.Error("author not found", zap.Error(err))
		return GetProblemOutput{}, ErrAccountNotFound
	}

	return GetProblemOutput{
		Problem: p.dbProblemToLogicProblem(problem, author.Name),
	}, nil
}

func (p *problemLogic) GetProblemList(ctx context.Context, in GetProblemListInput) (GetProblemListOutput, error) {
	logger := p.logger.With(zap.String("method", "GetProblemList"))

	_, _, requestingAccountRole, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return GetProblemListOutput{}, ErrTokenInvalid
	}

	requiredPermissions := []gorbac.Permission{PermissionProblemsReadAll}
	hasPermission, err := p.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountRole)], requiredPermissions...)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check permission")
		return GetProblemListOutput{}, ErrInternal
	}
	if !hasPermission {
		return GetProblemListOutput{}, ErrPermissionDenied
	}

	dbProblemList, err := p.problemDataAccessor.GetProblemList(ctx, in.Offset, in.Limit)
	if err != nil {
		logger.Error("failed to get problem list", zap.Error(err))
		return GetProblemListOutput{}, err
	}
	totalProblemCount, err := p.problemDataAccessor.GetProblemCount(ctx)
	if err != nil {
		logger.Error("failed to get problem count", zap.Error(err))
		return GetProblemListOutput{}, err
	}

	var problemList []Problem
	for _, pb := range dbProblemList {
		author, err := p.accountDataAccessor.GetAccountByID(ctx, pb.AuthorID)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to get author")
			return GetProblemListOutput{}, ErrInternal
		}
		if author.ID == 0 {
			logger.With(zap.Error(err)).Error("author not found")
			return GetProblemListOutput{}, ErrAccountNotFound
		}

		problemList = append(problemList, p.dbProblemToLogicProblem(pb, author.Name))
	}

	return GetProblemListOutput{
		Problems:          problemList,
		TotalProblemCount: totalProblemCount,
	}, nil
}

func (p *problemLogic) UpdateProblem(ctx context.Context, in UpdateProblemInput) (UpdateProblemOutput, error) {
	logger := p.logger.With(zap.String("method", "UpdateProblem"))

	requestingAccountID, _, requestingAccountRole, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return UpdateProblemOutput{}, ErrTokenInvalid
	}

	problem, err := p.problemDataAccessor.GetProblemByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return UpdateProblemOutput{}, err
	}
	if problem.ID == 0 {
		logger.Error("problem not found", zap.Error(err))
		return UpdateProblemOutput{}, ErrProblemNotFound
	}

	requiredPermissions := []gorbac.Permission{PermissionAccountsWriteAll}
	if requestingAccountID == problem.AuthorID {
		requiredPermissions = append(requiredPermissions, PermissionAccountsWriteSelf)
	}

	hasPermission, err := p.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountRole)], requiredPermissions...)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check permission")
		return UpdateProblemOutput{}, ErrInternal
	}
	if !hasPermission {
		return UpdateProblemOutput{}, ErrPermissionDenied
	}

	author, err := p.accountDataAccessor.GetAccountByID(ctx, problem.AuthorID)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get author")
		return UpdateProblemOutput{}, ErrInternal
	}
	if author.ID == 0 {
		logger.With(zap.Error(err)).Error("author not found")
		return UpdateProblemOutput{}, ErrAccountNotFound
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
		Problem: p.dbProblemToLogicProblem(updatedProblem, author.Name),
	}, nil
}

func (p *problemLogic) DeleteProblem(ctx context.Context, in DeleteProblemInput) error {
	logger := p.logger.With(zap.String("method", "DeleteProblem"))

	problem, err := p.problemDataAccessor.GetProblemByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return ErrInternal
	}
	if problem.ID == 0 {
		logger.Error("problem not found", zap.Error(err))
		return ErrProblemNotFound
	}

	requestingAccountID, _, requestingAccountRole, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return ErrTokenInvalid
	}

	requiredPermissions := []gorbac.Permission{PermissionAccountsWriteAll}
	if requestingAccountID == problem.AuthorID {
		requiredPermissions = append(requiredPermissions, PermissionAccountsWriteSelf)
	}

	hasPermission, err := p.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountRole)], requiredPermissions...)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check permission")
		return ErrInternal
	}
	if !hasPermission {
		return ErrPermissionDenied
	}

	err = p.problemDataAccessor.DeleteProblem(ctx, in.ID)
	if err != nil {
		logger.Error("failed to delete problem", zap.Error(err))
		return ErrInternal
	}

	return nil
}

func (p *problemLogic) dbProblemToLogicProblem(dbProblem database.Problem, accountName string) Problem {
	return Problem{
		ID:          dbProblem.ID,
		DisplayName: dbProblem.DisplayName,
		AuthorId:    dbProblem.AuthorID,
		AuthorName:  accountName,
		Description: dbProblem.Description,
		TimeLimit:   dbProblem.TimeLimit,
		MemoryLimit: dbProblem.MemoryLimit,
	}
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
