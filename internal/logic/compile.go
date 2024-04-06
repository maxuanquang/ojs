package logic

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

const (
	defaultHostWorkingDir              = "/tmp/ojs-compile"
	defaultContainerWorkingDir         = "/tmp/ojs-compile"
	modeOwnerAllPermission             = 0700
	defaultCPUPeriod                   = 100000
	SourceFilePathPlaceholder          = "$SOURCE"
	CompiledProgramFilePathPlaceholder = "$PROGRAM"
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

	timeoutDuration, err := compileConfig.GetTimeoutInTimeDuration()
	if err != nil {
		c.logger.With(zap.Error(err)).Error("failed to get timeout duration")
		return nil, err
	}

	c.timeoutDuration = timeoutDuration

	memoryInBytes, err := compileConfig.GetMemoryInBytes()
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
	hostWorkingDir := defaultHostWorkingDir

	err := os.MkdirAll(hostWorkingDir, modeOwnerAllPermission)
	if err != nil {
		c.logger.With(zap.Error(err)).Error("failed to create temp dir")
		return CompileOutput{}, err
	}
	c.logger.Info("temp dir created", zap.String("temp_dir", hostWorkingDir))

	sourceFileName := uuid.NewString() + c.compileConfig.SourceFileExtension
	sourceFile, err := c.createSourceFile(ctx, hostWorkingDir, sourceFileName, content)
	if err != nil {
		c.logger.With(zap.Error(err)).Error("failed to create source file")
		return CompileOutput{}, err
	}
	c.logger.Info("source file created", zap.String("source_file", sourceFile.Name()))

	defer func() {
		if err := os.RemoveAll(sourceFile.Name()); err != nil {
			c.logger.With(zap.Error(err)).Error("failed to remove temp dir")
		}
	}()

	// Interpreted languages don't need to be compiled
	if c.compileConfig == nil {
		return CompileOutput{
			ProgramFilePath: sourceFile.Name(),
		}, nil
	}

	compileOutput, err := c.compileSourceFile(ctx, hostWorkingDir, sourceFile)
	if err != nil {
		c.logger.With(zap.Error(err)).Error("failed to compile source file")
		return CompileOutput{}, err
	}

	return compileOutput, nil
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

func (c *compileLogic) createSourceFile(_ context.Context, hostWorkingDir, fileName, content string) (*os.File, error) {
	logger := c.logger.With(zap.String("file_name", fileName)).With(zap.String("host_working_dir", hostWorkingDir))

	sourceFilePath := filepath.Join(hostWorkingDir, fileName)
	sourceFile, err := os.Create(sourceFilePath)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create source file")
		return nil, err
	}

	_, err = sourceFile.WriteString(content)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to write source file")
		return nil, err
	}

	defer sourceFile.Close()

	return sourceFile, nil
}

func (c *compileLogic) compileSourceFile(ctx context.Context, hostWorkingDir string, sourceFile *os.File) (CompileOutput, error) {
	logger := c.logger.With(zap.String("file_name", sourceFile.Name()))

	hostCompiledProgramFilePath := sourceFile.Name() + ".out"
	containerWorkingDir := defaultContainerWorkingDir
	containerSourceFilePath := filepath.Join(containerWorkingDir, filepath.Base(sourceFile.Name()))
	containerCompiledProgramFilePath := filepath.Join(containerWorkingDir, filepath.Base(hostCompiledProgramFilePath))

	dockerContainerCtx, dockerContainerCancelFunc := context.WithTimeout(ctx, c.timeoutDuration)
	defer dockerContainerCancelFunc()

	containerCreateResponse, err := c.dockerClient.ContainerCreate(
		dockerContainerCtx,
		&container.Config{
			Image:        c.compileConfig.Image,
			WorkingDir:   containerWorkingDir,
			Cmd:          c.getCompileCommand(containerSourceFilePath, containerCompiledProgramFilePath),
			AttachStdout: true,
			AttachStderr: true,
		},
		&container.HostConfig{
			Binds:       []string{fmt.Sprintf("%s:%s", hostWorkingDir, containerWorkingDir)},
			NetworkMode: "none",
			Resources: container.Resources{
				CPUPeriod: defaultCPUPeriod,
				CPUQuota:  int64(c.compileConfig.CPUs * defaultCPUPeriod),
				Memory:    c.memoryLimitInBytes,
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create container")
		return CompileOutput{}, err
	}

	defer func() {
		err = c.dockerClient.ContainerRemove(ctx, containerCreateResponse.ID, container.RemoveOptions{})
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to remove container")
		}
	}()

	containerID := containerCreateResponse.ID
	containerAttachResponse, err := c.dockerClient.ContainerAttach(
		dockerContainerCtx,
		containerID,
		container.AttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		},
	)
	if err != nil {
		logger.With(zap.String("container_id", containerID)).With(zap.Error(err)).Error("failed to attach container")
		return CompileOutput{}, err
	}

	defer containerAttachResponse.Close()

	err = c.dockerClient.ContainerStart(
		dockerContainerCtx,
		containerID,
		container.StartOptions{},
	)
	if err != nil {
		logger.With(zap.String("container_id", containerID)).With(zap.Error(err)).Error("failed to start container")
	}

	dataChan, errChan := c.dockerClient.ContainerWait(
		dockerContainerCtx,
		containerID,
		container.WaitConditionNotRunning,
	)
	select {
	case data := <-dataChan:
		if data.StatusCode != 0 {
			stdoutBuffer := new(bytes.Buffer)
			stderrBuffer := new(bytes.Buffer)
			_, err = stdcopy.StdCopy(stdoutBuffer, stderrBuffer, containerAttachResponse.Reader)
			if err != nil {
				logger.With(zap.Error(err)).Error("failed to read container output")
				return CompileOutput{}, err
			}

			compileOutput := CompileOutput{
				ReturnCode: int(data.StatusCode),
				Stdout:     stdoutBuffer.String(),
				Stderr:     stderrBuffer.String(),
			}
			logger.With(zap.Any("compile_output", compileOutput)).Info("failed to compile source file, compiler exited with non-zero code")
			return compileOutput, nil
		}

		logger.Info("source file compiled successfully")
		return CompileOutput{
			ProgramFilePath: hostCompiledProgramFilePath,
		}, nil
	case err := <-errChan:
		logger.With(zap.String("container_id", containerID)).With(zap.Error(err)).Error("failed to wait container")
		return CompileOutput{}, err
	}
}

func (c *compileLogic) getCompileCommand(sourceFilePath, compiledProgramFilePath string) strslice.StrSlice {
	commandTemplate := make([]string, len(c.compileConfig.CommandTemplate))
	for i := range c.compileConfig.CommandTemplate {
		switch c.compileConfig.CommandTemplate[i] {
		case SourceFilePathPlaceholder:
			commandTemplate[i] = sourceFilePath
		case CompiledProgramFilePathPlaceholder:
			commandTemplate[i] = compiledProgramFilePath
		default:
			commandTemplate[i] = c.compileConfig.CommandTemplate[i]
		}
	}
	return strslice.StrSlice(commandTemplate)
}
