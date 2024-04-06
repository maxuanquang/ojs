package logic

import (
	"context"
	"os"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	// "github.com/google/uuid"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

const (
	TempDirPattern = "ojs-compile-"
)

type CompileOutput struct {
	ProgramFilePath string
	ReturnCode      int
	Stdout          string
	Stderr          string
}

type CompileLogic interface {
	Compile(ctx context.Context, content string) (CompileOutput, error)
}

func NewCompileLogic(
	logger *zap.Logger,
	dockerClient *client.Client,
	language string,
	compileConfig *configs.Compile,
	appArguments utils.Arguments,
) (CompileLogic, error) {
	c := &compileLogic{
		logger:        logger.With(zap.String("language", language)).With(zap.Any("compile_config", compileConfig)),
		dockerClient:  dockerClient,
		language:      language,
		compileConfig: compileConfig,
		appArguments:  appArguments,
	}

	if compileConfig == nil {
		return c, nil
	}

	timeoutDuration, err := compileConfig.GetTimeoutDuration()
	if err != nil {
		c.logger.With(zap.Error(err)).Error("failed to get timeout duration")
		return nil, err
	}

	c.timeoutDuration = timeoutDuration

	memoryInBytes, err := compileConfig.GetMemoryLimitInBytes()
	if err != nil {
		c.logger.With(zap.Error(err)).Error("failed to get memory in bytes")
		return nil, err
	}

	c.memoryLimitInBytes = int64(memoryInBytes)

	if c.appArguments.PullImageAtStartUp {
		if err := c.pullImage(); err != nil {
			return nil, err
		}
	} else {
		go func() {
			c.pullImage()
		}()
	}

	return c, nil
}

type compileLogic struct {
	logger        *zap.Logger
	dockerClient  *client.Client
	language      string
	compileConfig *configs.Compile
	appArguments  utils.Arguments

	timeoutDuration    time.Duration
	memoryLimitInBytes int64
}

// Compile implements CompileLogic.
func (c *compileLogic) Compile(ctx context.Context, content string) (CompileOutput, error) {
	hostWorkingDir, err := os.MkdirTemp("", TempDirPattern)
	if err != nil {
		c.logger.With(zap.Error(err)).Error("failed to create temp dir")
		return CompileOutput{}, err
	}

	defer func() {
		if err := os.RemoveAll(hostWorkingDir); err != nil {
			c.logger.With(zap.Error(err)).Error("failed to remove temp dir")
		}
	}()

	// sourceFile, err = c.createSourceFile(ctx, hostWorkingDir, uuid.NewString(), content)
	// if err != nil {
	// 	c.logger.With(zap.Error(err)).Error("failed to create source file")
	// 	return CompileOutput{}, err
	// }

	// return c.compileSourceFile(ctx, sourceFile)
	return CompileOutput{}, nil
}

func (c *compileLogic) pullImage() error {
	c.logger.Info("pulling image")
	_, err := c.dockerClient.ImagePull(context.Background(), c.compileConfig.Image, image.PullOptions{})
	if err != nil {
		c.logger.With(zap.Error(err)).Error("failed to pull image")
		return err
	}

	c.logger.Info("compile image pulled successfully")
	return nil
}
