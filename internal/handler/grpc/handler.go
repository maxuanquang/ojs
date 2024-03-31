package grpc

import (
	"context"

	ojs "github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	"github.com/maxuanquang/ojs/internal/logic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	AuthTokenMetadataName         = "OJS_AUTH"
	GRPCGatewayCookieMetadataName = "grpcgateway-cookie"
)

func NewHandler(
	accountLogic logic.AccountLogic,
	problemLogic logic.ProblemLogic,
) ojs.OjsServiceServer {
	return &Handler{
		accountLogic: accountLogic,
		problemLogic: problemLogic,
	}
}

type Handler struct {
	ojs.UnimplementedOjsServiceServer
	accountLogic logic.AccountLogic
	problemLogic logic.ProblemLogic
}

// CreateProblem implements ojs.OjsServiceServer.
func (h *Handler) CreateProblem(ctx context.Context, in *ojs.CreateProblemRequest) (*ojs.CreateProblemResponse, error) {
	resp, err := h.problemLogic.CreateProblem(
		ctx,
		logic.CreateProblemInput{
			Token:       h.getAuthTokenFromMetadata(ctx),
			DisplayName: in.DisplayName,
			Description: in.Description,
			TimeLimit:   in.TimeLimit,
			MemoryLimit: in.MemoryLimit,
		},
	)
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &ojs.CreateProblemResponse{
		Problem: &ojs.Problem{
			Id:          resp.Problem.ID,
			AuthorId:    resp.Problem.AuthorId,
			DisplayName: resp.Problem.DisplayName,
			Description: resp.Problem.Description,
			TimeLimit:   resp.Problem.TimeLimit,
			MemoryLimit: resp.Problem.MemoryLimit,
		},
	}, nil

}

// CreateAccount implements ojs.OjsServiceServer.
func (h *Handler) CreateAccount(ctx context.Context, in *ojs.CreateAccountRequest) (*ojs.CreateAccountResponse, error) {
	account, err := h.accountLogic.CreateAccount(ctx, logic.CreateAccountInput{
		Name:     in.Name,
		Password: in.Password,
		Role:     in.Role,
	})
	if err != nil {
		return nil, clientResponseError(err)
	}
	return &ojs.CreateAccountResponse{
		Account: &ojs.Account{
			Id:   account.ID,
			Name: account.Name,
			Role: account.Role,
		},
	}, nil
}

// GetAccount implements ojs.OjsServiceServer.
func (h *Handler) GetAccount(ctx context.Context, in *ojs.GetAccountRequest) (*ojs.GetAccountResponse, error) {
	output, err := h.accountLogic.GetAccount(ctx, logic.GetAccountInput{
		ID: in.Id,
	})
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &ojs.GetAccountResponse{
		Account: &ojs.Account{
			Id:   output.ID,
			Name: output.Name,
			Role: output.Role,
		},
	}, nil
}

// CreateSession implements ojs.OjsServiceServer.
func (h *Handler) CreateSession(ctx context.Context, in *ojs.CreateSessionRequest) (*ojs.CreateSessionResponse, error) {
	session, err := h.accountLogic.CreateSession(
		ctx,
		logic.CreateSessionInput{
			Name:     in.Name,
			Password: in.Password,
		},
	)
	if err != nil {
		return nil, clientResponseError(err)
	}

	err = grpc.SendHeader(ctx, metadata.Pairs(AuthTokenMetadataName, session.Token))
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &ojs.CreateSessionResponse{
		Account: &ojs.Account{
			Id:   session.ID,
			Name: session.Name,
			Role: session.Role,
		},
	}, nil
}

// DeleteSession implements ojs.OjsServiceServer.
func (h *Handler) DeleteSession(ctx context.Context, in *ojs.DeleteSessionRequest) (*ojs.DeleteSessionResponse, error) {
	err := h.accountLogic.DeleteSession(
		ctx,
		logic.DeleteSessionInput{
			Token: h.getAuthTokenFromMetadata(ctx),
		},
	)
	if err != nil {
		return nil, clientResponseError(err)
	}

	err = grpc.SendHeader(ctx, metadata.Pairs(AuthTokenMetadataName, ""))
	if err != nil {
		return nil, clientResponseError(err)
	}

	return &ojs.DeleteSessionResponse{}, nil
}

func (h *Handler) GetProblem(ctx context.Context, in *ojs.GetProblemRequest) (*ojs.GetProblemResponse, error) {
	// Call the corresponding method of h.problemLogic
	output, err := h.problemLogic.GetProblem(
		ctx,
		logic.GetProblemInput{
			Token: h.getAuthTokenFromMetadata(ctx),
			ID:    in.GetId(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	response := &ojs.GetProblemResponse{
		Problem: &ojs.Problem{
			Id:          output.Problem.ID,
			DisplayName: output.Problem.DisplayName,
			AuthorId:    output.Problem.AuthorId,
			Description: output.Problem.Description,
			TimeLimit:   output.Problem.TimeLimit,
			MemoryLimit: output.Problem.MemoryLimit,
		},
	}

	return response, nil
}

// GetProblemList implements ojs.OjsServiceServer.
func (h *Handler) GetProblemList(ctx context.Context, in *ojs.GetProblemListRequest) (*ojs.GetProblemListResponse, error) {
	// Call the corresponding method of h.problemLogic
	output, err := h.problemLogic.GetProblemList(
		ctx,
		logic.GetProblemListInput{
			Token:  h.getAuthTokenFromMetadata(ctx),
			Offset: in.GetOffset(),
			Limit:  in.GetLimit(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	response := &ojs.GetProblemListResponse{
		TotalProblemCount: output.TotalProblemCount,
	}
	for _, pb := range output.Problems {
		response.Problems = append(response.Problems, &ojs.Problem{
			Id:          pb.ID,
			DisplayName: pb.DisplayName,
			AuthorId:    pb.AuthorId,
			Description: pb.Description,
			TimeLimit:   pb.TimeLimit,
			MemoryLimit: pb.MemoryLimit,
		})
	}

	return response, nil
}

// UpdateProblem implements ojs.OjsServiceServer.
func (h *Handler) UpdateProblem(ctx context.Context, in *ojs.UpdateProblemRequest) (*ojs.UpdateProblemResponse, error) {
	// Call the corresponding method of h.problemLogic
	output, err := h.problemLogic.UpdateProblem(
		ctx,
		logic.UpdateProblemInput{
			Token:       h.getAuthTokenFromMetadata(ctx),
			ID:          in.GetId(),
			DisplayName: in.DisplayName,
			Description: in.Description,
			TimeLimit:   in.TimeLimit,
			MemoryLimit: in.MemoryLimit,
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	response := &ojs.UpdateProblemResponse{
		Problem: &ojs.Problem{
			Id:          output.Problem.ID,
			DisplayName: output.Problem.DisplayName,
			AuthorId:    output.Problem.AuthorId,
			AuthorName:  output.Problem.AuthorName,
			Description: output.Problem.Description,
			TimeLimit:   output.Problem.TimeLimit,
			MemoryLimit: output.Problem.MemoryLimit,
		},
	}

	return response, nil
}

// DeleteProblem implements ojs.OjsServiceServer.
func (h *Handler) DeleteProblem(ctx context.Context, in *ojs.DeleteProblemRequest) (*ojs.DeleteProblemResponse, error) {
	// Call the corresponding method of h.problemLogic
	err := h.problemLogic.DeleteProblem(
		ctx,
		logic.DeleteProblemInput{
			Token: h.getAuthTokenFromMetadata(ctx),
			ID:    in.GetId(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Return success response
	return &ojs.DeleteProblemResponse{}, nil
}

func (h *Handler) CreateTestCase(ctx context.Context, in *ojs.CreateTestCaseRequest) (*ojs.CreateTestCaseResponse, error) {
	// Call the corresponding method of h.testCaseLogic
	output, err := h.problemLogic.CreateTestCase(
		ctx,
		logic.CreateTestCaseInput{
			OfProblemID: in.GetOfProblemId(),
			Input:       in.GetInput(),
			Output:      in.GetOutput(),
			IsHidden:    in.GetIsHidden(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	response := &ojs.CreateTestCaseResponse{
		TestCase: &ojs.TestCase{
			Id:          output.TestCase.ID,
			OfProblemId: output.TestCase.OfProblemID,
			Input:       output.TestCase.Input,
			Output:      output.TestCase.Output,
			IsHidden:    output.TestCase.IsHidden,
		},
	}

	return response, nil
}

func (h *Handler) GetTestCase(ctx context.Context, in *ojs.GetTestCaseRequest) (*ojs.GetTestCaseResponse, error) {
	// Call the corresponding method of h.testCaseLogic
	output, err := h.problemLogic.GetTestCase(
		ctx,
		logic.GetTestCaseInput{
			ID: in.GetId(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	response := &ojs.GetTestCaseResponse{
		TestCase: &ojs.TestCase{
			Id:          output.TestCase.ID,
			OfProblemId: output.TestCase.OfProblemID,
			Input:       output.TestCase.Input,
			Output:      output.TestCase.Output,
			IsHidden:    output.TestCase.IsHidden,
		},
	}

	return response, nil
}

func (h *Handler) GetProblemTestCaseList(ctx context.Context, in *ojs.GetProblemTestCaseListRequest) (*ojs.GetProblemTestCaseListResponse, error) {
	// Call the corresponding method of h.testCaseLogic
	output, err := h.problemLogic.GetProblemTestCaseList(
		ctx,
		logic.GetProblemTestCaseListInput{
			OfProblemID: in.GetId(),
			Offset:      in.GetOffset(),
			Limit:       in.GetLimit(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	var testCases []*ojs.TestCase
	for _, testCase := range output.TestCases {
		testCases = append(testCases, &ojs.TestCase{
			Id:          testCase.ID,
			OfProblemId: testCase.OfProblemID,
			Input:       testCase.Input,
			Output:      testCase.Output,
		})
	}

	response := &ojs.GetProblemTestCaseListResponse{
		TestCases:           testCases,
		TotalTestCasesCount: output.TotalTestCasesCount,
	}

	return response, nil
}

func (h *Handler) DeleteTestCase(ctx context.Context, in *ojs.DeleteTestCaseRequest) (*ojs.DeleteTestCaseResponse, error) {
	// Call the corresponding method of h.testCaseLogic
	err := h.problemLogic.DeleteTestCase(
		ctx,
		logic.DeleteTestCaseInput{
			ID: in.GetId(),
		},
	)
	if err != nil {
		return nil, err
	}

	// No need to return any response for delete operation
	return &ojs.DeleteTestCaseResponse{}, nil
}

func (h *Handler) UpdateTestCase(ctx context.Context, in *ojs.UpdateTestCaseRequest) (*ojs.UpdateTestCaseResponse, error) {
	// Call the corresponding method of h.testCaseLogic
	updatedTestCase, err := h.problemLogic.UpdateTestCase(
		ctx,
		logic.UpdateTestCaseInput{
			ID:       in.GetId(),
			Input:    in.GetInput(),
			Output:   in.GetOutput(),
			IsHidden: in.GetIsHidden(),
		},
	)
	if err != nil {
		return nil, err
	}

	// No need to return any response for update operation
	return &ojs.UpdateTestCaseResponse{
		TestCase: &ojs.TestCase{
			Id:          updatedTestCase.TestCase.ID,
			OfProblemId: updatedTestCase.TestCase.OfProblemID,
			Input:       updatedTestCase.TestCase.Input,
			Output:      updatedTestCase.TestCase.Output,
			IsHidden:    updatedTestCase.TestCase.IsHidden,
		},
	}, nil
}

// CreateSubmission implements ojs.OjsServiceServer.
func (h *Handler) CreateSubmission(context.Context, *ojs.CreateSubmissionRequest) (*ojs.CreateSubmissionResponse, error) {
	panic("unimplemented")
}

// GetProblemSubmissionList implements ojs.OjsServiceServer.
func (h *Handler) GetProblemSubmissionList(context.Context, *ojs.GetProblemSubmissionListRequest) (*ojs.GetProblemSubmissionListResponse, error) {
	panic("unimplemented")
}

// GetSubmission implements ojs.OjsServiceServer.
func (h *Handler) GetSubmission(context.Context, *ojs.GetSubmissionRequest) (*ojs.GetSubmissionResponse, error) {
	panic("unimplemented")
}

// GetSubmissionList implements ojs.OjsServiceServer.
func (h *Handler) GetSubmissionList(context.Context, *ojs.GetSubmissionListRequest) (*ojs.GetSubmissionListResponse, error) {
	panic("unimplemented")
}

// GetAccountProblemSubmissionList implements ojs.OjsServiceServer.
func (h *Handler) GetAccountProblemSubmissionList(context.Context, *ojs.GetAccountProblemSubmissionListRequest) (*ojs.GetAccountProblemSubmissionListResponse, error) {
	panic("unimplemented")
}

// GetServerInfo implements ojs.OjsServiceServer.
func (h *Handler) GetServerInfo(context.Context, *ojs.GetServerInfoRequest) (*ojs.GetServerInfoResponse, error) {
	panic("unimplemented")
}

// UpdateSetting implements ojs.OjsServiceServer.
func (h *Handler) UpdateSetting(context.Context, *ojs.UpdateSettingRequest) (*ojs.UpdateSettingResponse, error) {
	panic("unimplemented")
}

// GetAndUpdateFirstSubmittedSubmissionToExecuting implements ojs.OjsServiceServer.
func (h *Handler) GetAndUpdateFirstSubmittedSubmissionToExecuting(context.Context, *ojs.GetAndUpdateFirstSubmittedSubmissionToExecutingRequest) (*ojs.GetAndUpdateFirstSubmittedSubmissionToExecutingResponse, error) {
	panic("unimplemented")
}

func (h *Handler) getAuthTokenFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	authTokenValues := md.Get(AuthTokenMetadataName)
	if len(authTokenValues) == 0 {
		return ""
	}

	return authTokenValues[0]
}
