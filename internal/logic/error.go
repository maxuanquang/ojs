package logic

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrProblemNotFound  = status.Error(codes.NotFound, "problem not found")
	ErrTestCaseNotFound = status.Error(codes.NotFound, "test case not found")
	ErrAccountNotFound  = status.Error(codes.NotFound, "account not found")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "permission denied")
	ErrInternal         = status.Error(codes.Internal, "internal error")
	ErrTokenInvalid     = status.Error(codes.Unauthenticated, "invalid authentication token")
)
