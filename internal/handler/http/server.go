package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/maxuanquang/ojs/internal/configs"
	gw "github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	grpcHandler "github.com/maxuanquang/ojs/internal/handler/grpc"
	"github.com/maxuanquang/ojs/internal/handler/http/servemuxoption"
)

const (
	AuthCookieName = "OJS_AUTH"
)

type Server interface {
	Start(ctx context.Context) error
}

func NewServer(
	httpConfig configs.HTTP,
	grpcConfig configs.GRPC,
	authConfig configs.Auth,
	logger *zap.Logger,
) Server {
	return &server{
		httpConfig: httpConfig,
		grpcConfig: grpcConfig,
		authConfig: authConfig,
		logger:     logger,
	}
}

type server struct {
	httpConfig configs.HTTP
	grpcConfig configs.GRPC
	authConfig configs.Auth
	logger     *zap.Logger
}

func (s *server) Start(ctx context.Context) error {
	mux := runtime.NewServeMux(
		servemuxoption.WithAuthCookieToAuthMetadata(AuthCookieName, grpcHandler.AuthTokenMetadataName),
		servemuxoption.WithAuthMetadataToAuthCookie(AuthCookieName, grpcHandler.AuthTokenMetadataName, s.authConfig.Token.GetTokenDuration()),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterOjsServiceHandlerFromEndpoint(
		ctx,
		mux,
		s.grpcConfig.Address,
		opts,
	)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	fmt.Printf("http server is running on %s\n", s.httpConfig.Address)
	return http.ListenAndServe(s.httpConfig.Address, mux)
}
