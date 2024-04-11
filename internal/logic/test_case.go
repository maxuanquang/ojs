package logic

import (
	"context"

	"github.com/maxuanquang/ojs/internal/dataaccess/database"
	"github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	"github.com/mikespook/gorbac"
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
	roleLogic RoleLogic,
) TestCaseLogic {
	return &testCaseLogic{
		logger:                 logger,
		accountDataAccessor:    accountDataAccessor,
		problemDataAccessor:    problemDataAccessor,
		submissionDataAccessor: submissionDataAccessor,
		testCaseDataAccessor:   testCaseDataAccessor,
		tokenLogic:             tokenLogic,
		roleLogic:              roleLogic,
	}
}

type testCaseLogic struct {
	logger                 *zap.Logger
	accountDataAccessor    database.AccountDataAccessor
	problemDataAccessor    database.ProblemDataAccessor
	submissionDataAccessor database.SubmissionDataAccessor
	testCaseDataAccessor   database.TestCaseDataAccessor
	tokenLogic             TokenLogic
	roleLogic              RoleLogic
}

func (t *testCaseLogic) CreateTestCase(ctx context.Context, in CreateTestCaseInput) (CreateTestCaseOutput, error) {
	logger := t.logger.With(zap.Any("create_test_case_input", in))

	_, _, requestingAccountRole, _, err := t.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return CreateTestCaseOutput{}, ErrTokenInvalid
	}

	requiredPermissions := []gorbac.Permission{PermissionTestCasesWriteAll, PermissionTestCasesWriteSelf}
	hasPermission, err := t.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountRole)], requiredPermissions...)
	if err != nil {
		logger.Error("failed to check account permission", zap.Error(err))
		return CreateTestCaseOutput{}, ErrInternal
	}
	if !hasPermission {
		return CreateTestCaseOutput{}, ErrPermissionDenied
	}

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
		TestCase: t.dbTestCaseToLogicTestCase(createdTestCase),
	}, nil
}

func (t *testCaseLogic) GetTestCase(ctx context.Context, in GetTestCaseInput) (GetTestCaseOutput, error) {
	logger := t.logger.With(zap.Any("get_test_case_input", in))

	requestingAccountID, _, _, _, err := t.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return GetTestCaseOutput{}, ErrTokenInvalid
	}

	dbTestCase, err := t.testCaseDataAccessor.GetTestCaseByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get test case", zap.Error(err))
		return GetTestCaseOutput{}, err
	}
	if dbTestCase.ID == 0 {
		err := ErrTestCaseNotFound
		logger.Error("test case not found", zap.Error(err))
		return GetTestCaseOutput{}, err
	}

	requiredPermissions := []gorbac.Permission{PermissionTestCasesReadAll}
	if dbTestCase.ID == requestingAccountID {
		requiredPermissions = append(requiredPermissions, PermissionTestCasesReadSelf)
	}

	hasPermission, err := t.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountID)], requiredPermissions...)
	if err != nil {
		logger.Error("failed to check account permission", zap.Error(err))
		return GetTestCaseOutput{}, ErrInternal
	}
	if !hasPermission {
		return GetTestCaseOutput{}, ErrPermissionDenied
	}

	return GetTestCaseOutput{
		TestCase: t.dbTestCaseToLogicTestCase(dbTestCase),
	}, nil
}

func (t *testCaseLogic) GetProblemTestCaseList(ctx context.Context, in GetProblemTestCaseListInput) (GetProblemTestCaseListOutput, error) {
	logger := t.logger.With(zap.String("method", "GetProblemTestCaseList"))

	requestingAccountID, _, _, _, err := t.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return GetProblemTestCaseListOutput{}, ErrTokenInvalid
	}

	dbProlem, err := t.problemDataAccessor.GetProblemByID(ctx, in.OfProblemID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return GetProblemTestCaseListOutput{}, ErrInternal
	}
	if dbProlem.ID == 0 {
		return GetProblemTestCaseListOutput{}, ErrProblemNotFound
	}

	requiredPermissions := []gorbac.Permission{PermissionProblemsReadAll}
	if dbProlem.AuthorID == requestingAccountID {
		requiredPermissions = append(requiredPermissions, PermissionProblemsReadSelf)
	}

	hasPermission, err := t.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountID)], requiredPermissions...)
	if err != nil {
		logger.Error("failed to check account permission", zap.Error(err))
		return GetProblemTestCaseListOutput{}, ErrInternal
	}
	if !hasPermission {
		return GetProblemTestCaseListOutput{}, ErrPermissionDenied
	}

	testCases, err := t.testCaseDataAccessor.GetProblemTestCaseList(ctx, in.OfProblemID, in.Offset, in.Limit)
	if err != nil {
		logger.Error("failed to get test case list", zap.Error(err))
		return GetProblemTestCaseListOutput{}, ErrInternal
	}

	var testCaseList []TestCase
	for _, tc := range testCases {
		testCaseList = append(testCaseList, t.dbTestCaseToLogicTestCase(tc))
	}

	totalTestCasesCount, err := t.testCaseDataAccessor.GetProblemTestCaseCount(ctx, in.OfProblemID)
	if err != nil {
		logger.Error("failed to get test case count", zap.Error(err))
		return GetProblemTestCaseListOutput{}, ErrInternal
	}

	return GetProblemTestCaseListOutput{
		TestCases:           testCaseList,
		TotalTestCasesCount: totalTestCasesCount,
	}, nil
}

func (t *testCaseLogic) UpdateTestCase(ctx context.Context, in UpdateTestCaseInput) (UpdateTestCaseOutput, error) {
	logger := t.logger.With(zap.String("method", "UpdateTestCase"))

	dbTestCase, err := t.testCaseDataAccessor.GetTestCaseByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get test case", zap.Error(err))
		return UpdateTestCaseOutput{}, ErrInternal
	}
	if dbTestCase.ID == 0 {
		logger.Error("test case not found", zap.Error(err))
		return UpdateTestCaseOutput{}, ErrTestCaseNotFound
	}

	dbProlem, err := t.problemDataAccessor.GetProblemByID(ctx, dbTestCase.OfProblemID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return UpdateTestCaseOutput{}, ErrInternal
	}
	if dbProlem.ID == 0 {
		return UpdateTestCaseOutput{}, ErrProblemNotFound
	}

	requestingAccountID, _, _, _, err := t.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return UpdateTestCaseOutput{}, ErrTokenInvalid
	}

	requiredPermissions := []gorbac.Permission{PermissionProblemsWriteAll}
	if dbProlem.AuthorID == requestingAccountID {
		requiredPermissions = append(requiredPermissions, PermissionProblemsWriteSelf)
	}

	hasPermission, err := t.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountID)], requiredPermissions...)
	if err != nil {
		logger.Error("failed to check account permission", zap.Error(err))
		return UpdateTestCaseOutput{}, ErrInternal
	}
	if !hasPermission {
		return UpdateTestCaseOutput{}, ErrPermissionDenied
	}

	updatedDbTestCase, err := t.testCaseDataAccessor.UpdateTestCase(ctx, database.TestCase{
		ID:       in.ID,
		Input:    in.Input,
		Output:   in.Output,
		IsHidden: in.IsHidden,
	})
	if err != nil {
		logger.Error("failed to update test case", zap.Error(err))
		return UpdateTestCaseOutput{}, ErrInternal
	}

	return UpdateTestCaseOutput{
		TestCase: t.dbTestCaseToLogicTestCase(updatedDbTestCase),
	}, nil
}

func (t *testCaseLogic) DeleteTestCase(ctx context.Context, in DeleteTestCaseInput) error {
	logger := t.logger.With(zap.String("method", "DeleteTestCase"))

	// Check if the test case exists
	dbTestCase, err := t.testCaseDataAccessor.GetTestCaseByID(ctx, in.ID)
	if err != nil {
		logger.Error("failed to get test case", zap.Error(err))
		return err
	}
	if dbTestCase.ID == 0 {
		err := ErrTestCaseNotFound
		logger.Error("test case not found", zap.Error(err))
		return err
	}

	dbProlem, err := t.problemDataAccessor.GetProblemByID(ctx, dbTestCase.OfProblemID)
	if err != nil {
		logger.Error("failed to get problem", zap.Error(err))
		return ErrInternal
	}
	if dbProlem.ID == 0 {
		return ErrProblemNotFound
	}

	requestingAccountID, _, _, _, err := t.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		logger.Error("failed to verify token", zap.Error(err))
		return ErrTokenInvalid
	}

	requiredPermissions := []gorbac.Permission{PermissionProblemsWriteAll}
	if dbProlem.AuthorID == requestingAccountID {
		requiredPermissions = append(requiredPermissions, PermissionProblemsWriteSelf)
	}

	hasPermission, err := t.roleLogic.AccountHasPermission(ctx, ojs.Role_name[int32(requestingAccountID)], requiredPermissions...)
	if err != nil {
		logger.Error("failed to check account permission", zap.Error(err))
		return ErrInternal
	}
	if !hasPermission {
		return ErrPermissionDenied
	}

	err = t.testCaseDataAccessor.DeleteTestCase(ctx, in.ID)
	if err != nil {
		logger.Error("failed to delete test case", zap.Error(err))
		return ErrInternal
	}

	return nil
}

func (t *testCaseLogic) dbTestCaseToLogicTestCase(dbTestCase database.TestCase) TestCase {
	return TestCase{
		ID:          dbTestCase.ID,
		OfProblemID: dbTestCase.OfProblemID,
		Input:       dbTestCase.Input,
		Output:      dbTestCase.Output,
		IsHidden:    dbTestCase.IsHidden,
	}
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
	Token       string
}

type CreateTestCaseOutput struct {
	TestCase TestCase
}

type GetTestCaseInput struct {
	ID    uint64
	Token string
}

type GetTestCaseOutput struct {
	TestCase TestCase
}

type GetProblemTestCaseListInput struct {
	OfProblemID uint64
	Offset      uint64
	Limit       uint64
	Token       string
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
	Token    string
}

type UpdateTestCaseOutput struct {
	TestCase TestCase
}

type DeleteTestCaseInput struct {
	ID    uint64
	Token string
}

type DeleteTestCaseOutput struct{}
