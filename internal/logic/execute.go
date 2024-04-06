package logic

import "context"

type ExecuteLogic interface {
	Execute(ctx context.Context, programFilePath string, language string, programInput string) (string, error)
}
