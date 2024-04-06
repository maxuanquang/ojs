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
	submissionLogic logic.SubmissionLogic,
	testCaseLogic logic.TestCaseLogic,
) ojs.OjsServiceServer {
	return &Handler{
		accountLogic:    accountLogic,
		problemLogic:    problemLogic,
		submissionLogic: submissionLogic,
		testCaseLogic:   testCaseLogic,
	}
}

type Handler struct {
	ojs.UnimplementedOjsServiceServer
	accountLogic    logic.AccountLogic
	problemLogic    logic.ProblemLogic
	submissionLogic logic.SubmissionLogic
	testCaseLogic   logic.TestCaseLogic
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
			AuthorName:  resp.Problem.AuthorName,
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
	output, err := h.testCaseLogic.CreateTestCase(
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
	output, err := h.testCaseLogic.GetTestCase(
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
	output, err := h.testCaseLogic.GetProblemTestCaseList(
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
	err := h.testCaseLogic.DeleteTestCase(
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
	updatedTestCase, err := h.testCaseLogic.UpdateTestCase(
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
func (h *Handler) CreateSubmission(ctx context.Context, in *ojs.CreateSubmissionRequest) (*ojs.CreateSubmissionResponse, error) {
	// Call the corresponding method of h.submissionLogic
	output, err := h.submissionLogic.CreateSubmission(
		ctx,
		logic.CreateSubmissionInput{
			Token:       h.getAuthTokenFromMetadata(ctx),
			OfProblemID: in.GetOfProblemId(),
			Content:     in.GetContent(),
			Language:    in.GetLanguage(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	response := &ojs.CreateSubmissionResponse{
		Submission: &ojs.Submission{
			Id:          output.Submission.ID,
			OfProblemId: output.Submission.OfProblemID,
			AuthorId:    output.Submission.AuthorID,
			Content:     output.Submission.Content,
			Language:    output.Submission.Language,
			Status:      output.Submission.Status,
			Result:      output.Submission.Result,
		},
	}

	return response, nil
}

// GetProblemSubmissionList implements ojs.OjsServiceServer.
func (h *Handler) GetProblemSubmissionList(ctx context.Context, in *ojs.GetProblemSubmissionListRequest) (*ojs.GetProblemSubmissionListResponse, error) {
	// Call the corresponding method of h.submissionLogic
	output, err := h.submissionLogic.GetProblemSubmissionList(
		ctx,
		logic.GetProblemSubmissionListInput{
			OfProblemID: in.GetId(),
			Offset:      in.GetOffset(),
			Limit:       in.GetLimit(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	var submissions []*ojs.Submission
	for _, submission := range output.Submissions {
		submissions = append(submissions, &ojs.Submission{
			Id:          submission.ID,
			AuthorId:    submission.AuthorID,
			OfProblemId: submission.OfProblemID,
			Content:     submission.Content,
			Language:    submission.Language,
			Status:      submission.Status,
			Result:      submission.Result,
		})
	}

	response := &ojs.GetProblemSubmissionListResponse{
		Submissions:           submissions,
		TotalSubmissionsCount: output.TotalSubmissionsCount,
	}

	return response, nil
}

// GetSubmission implements ojs.OjsServiceServer.
func (h *Handler) GetSubmission(ctx context.Context, in *ojs.GetSubmissionRequest) (*ojs.GetSubmissionResponse, error) {
	// Call the corresponding method of h.submissionLogic
	output, err := h.submissionLogic.GetSubmission(
		ctx,
		logic.GetSubmissionInput{
			ID: in.GetId(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	response := &ojs.GetSubmissionResponse{
		Submission: &ojs.Submission{
			Id:          output.Submission.ID,
			AuthorId:    output.Submission.AuthorID,
			OfProblemId: output.Submission.OfProblemID,
			Content:     output.Submission.Content,
			Language:    output.Submission.Language,
			Status:      output.Submission.Status,
			Result:      output.Submission.Result,
		},
	}

	return response, nil
}

// GetSubmissionList implements ojs.OjsServiceServer.
func (h *Handler) GetSubmissionList(ctx context.Context, in *ojs.GetSubmissionListRequest) (*ojs.GetSubmissionListResponse, error) {
	// Call the corresponding method of h.submissionLogic
	output, err := h.submissionLogic.GetSubmissionList(
		ctx,
		logic.GetSubmissionListInput{
			Offset: in.GetOffset(),
			Limit:  in.GetLimit(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	var submissions []*ojs.Submission
	for _, submission := range output.Submissions {
		submissions = append(submissions, &ojs.Submission{
			Id:          submission.ID,
			AuthorId:    submission.AuthorID,
			OfProblemId: submission.OfProblemID,
			Content:     submission.Content,
			Language:    submission.Language,
			Status:      submission.Status,
			Result:      submission.Result,
		})
	}

	response := &ojs.GetSubmissionListResponse{
		Submissions:           submissions,
		TotalSubmissionsCount: output.TotalSubmissionsCount,
	}

	return response, nil
}

// GetAccountProblemSubmissionList implements ojs.OjsServiceServer.
func (h *Handler) GetAccountProblemSubmissionList(ctx context.Context, in *ojs.GetAccountProblemSubmissionListRequest) (*ojs.GetAccountProblemSubmissionListResponse, error) {
	// Call the corresponding method of h.submissionLogic
	output, err := h.submissionLogic.GetAccountProblemSubmissionList(
		ctx,
		logic.GetAccountProblemSubmissionListInput{
			Token:       h.getAuthTokenFromMetadata(ctx),
			OfProblemID: in.GetProblemId(),
			Offset:      in.GetOffset(),
			Limit:       in.GetLimit(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Format the response based on the result obtained
	var submissions []*ojs.Submission
	for _, submission := range output.Submissions {
		submissions = append(submissions, &ojs.Submission{
			Id:          submission.ID,
			AuthorId:    submission.AuthorID,
			OfProblemId: submission.OfProblemID,
			Content:     submission.Content,
			Language:    submission.Language,
			Status:      submission.Status,
			Result:      submission.Result,
		})
	}

	response := &ojs.GetAccountProblemSubmissionListResponse{
		Submissions:           submissions,
		TotalSubmissionsCount: output.TotalSubmissionsCount,
	}

	return response, nil
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
