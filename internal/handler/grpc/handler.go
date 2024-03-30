package grpc

import (
	"context"
	// "strings"

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
) ojs.OjsServiceServer {
	return &Handler{
		accountLogic: accountLogic,
	}
}

type Handler struct {
	ojs.UnimplementedOjsServiceServer
	accountLogic logic.AccountLogic
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