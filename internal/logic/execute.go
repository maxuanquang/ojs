package logic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

const (
	executeTimeoutPlaceholder         = "$TIMEOUT"
	executeProgramFilePathPlaceholder = "$PROGRAM"
	statusCodeTimeLimitExceeded       = 124
	statusCodeMemoryLimitExceeded     = 137
)

type ExecuteOutput struct {
	ReturnCode          int
	TimeLimitExceeded   bool
	MemoryLimitExceeded bool
	Stdout              string
	Stderr              string
}

type ExecuteLogic interface {
	Execute(
		ctx context.Context,
		programFilePath string,
		programInput string,
	) (ExecuteOutput, error)
}

func NewExecuteLogic(
	logger *zap.Logger,
	dockerClient *client.Client,
	language string,
	executeConfig *configs.Execute,
	appArguments utils.Arguments,
) (ExecuteLogic, error) {

	output := &executeLogic{
		logger:        logger,
		dockerClient:  dockerClient,
		language:      language,
		executeConfig: executeConfig,
		appArguments:  appArguments,
	}

	memory, err := executeConfig.GetMemoryInBytes()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get memory limit")
		return nil, err
	}

	output.memoryLimitInBytes = memory

	timeoutDuration, err := executeConfig.GetTimeoutInTimeDuration()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get timeout")
		return nil, err
	}

	output.timeoutDuration = timeoutDuration

	if appArguments.PullImageAtStartUp {
		if err := output.pullImage(); err != nil {
			return nil, err
		}
	} else {
		go func() {
			output.pullImage()
		}()
	}

	return output, nil
}

type executeLogic struct {
	logger        *zap.Logger
	dockerClient  *client.Client
	language      string
	executeConfig *configs.Execute
	appArguments  utils.Arguments

	timeoutDuration    time.Duration
	memoryLimitInBytes uint64
}

// Execute implements ExecuteLogic.
func (e *executeLogic) Execute(ctx context.Context, programFilePath string, programInput string) (ExecuteOutput, error) {
	logger := e.logger.With(zap.String("program_file_path", programFilePath))
	hostWorkingDir := filepath.Dir(programFilePath)
	programFileName := filepath.Base(programFilePath)

	defer func() {
		if err := os.RemoveAll(hostWorkingDir); err != nil {
			e.logger.With(zap.Error(err)).Error("failed to remove temp dir")
		}
	}()

	containerWorkingDir := hostWorkingDir
	containerProgramFilePath := filepath.Join(containerWorkingDir, programFileName)

	dockerContainerCtx, dockerContainerCancelFunc := context.WithTimeout(ctx, e.timeoutDuration)
	defer dockerContainerCancelFunc()

	containerCreateResponse, err := e.dockerClient.ContainerCreate(
		dockerContainerCtx,
		&container.Config{
			Image:        e.executeConfig.Image,
			WorkingDir:   containerWorkingDir,
			Cmd:          e.getExecuteCommand(e.timeoutDuration, containerProgramFilePath),
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin:    true,
			StdinOnce:    true,
		},
		&container.HostConfig{
			Binds:       []string{fmt.Sprintf("%s:%s", hostWorkingDir, containerWorkingDir)},
			NetworkMode: "none",
			Resources: container.Resources{
				CPUPeriod: defaultCPUPeriod,
				CPUQuota:  int64(e.executeConfig.CPUs * defaultCPUPeriod),
				Memory:    int64(e.memoryLimitInBytes),
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create container")
		return ExecuteOutput{}, err
	}

	defer func() {
		err = e.dockerClient.ContainerRemove(ctx, containerCreateResponse.ID, container.RemoveOptions{})
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to remove execute container")
		}
	}()

	containerID := containerCreateResponse.ID
	containerAttachResponse, err := e.dockerClient.ContainerAttach(
		dockerContainerCtx,
		containerID,
		container.AttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		},
	)
	if err != nil {
		logger.With(zap.String("container_id", containerID)).With(zap.Error(err)).Error("failed to attach container")
		return ExecuteOutput{}, err
	}

	defer containerAttachResponse.Close()

	_, err = containerAttachResponse.Conn.Write(append([]byte(programInput), '\n'))
	if err != nil {
		logger.With(zap.String("container_id", containerID)).With(zap.String("input", programInput)).With(zap.Error(err)).Error("failed to write to container")
		return ExecuteOutput{}, err
	}

	err = e.dockerClient.ContainerStart(
		dockerContainerCtx,
		containerID,
		container.StartOptions{},
	)
	if err != nil {
		logger.With(zap.String("container_id", containerID)).With(zap.Error(err)).Error("failed to start container")
		return ExecuteOutput{}, err
	}

	containerWaitCtx, containerWaitCancelFunc := context.WithTimeout(ctx, time.Minute)
	defer containerWaitCancelFunc()

	dataChan, errChan := e.dockerClient.ContainerWait(containerWaitCtx, containerID, container.WaitConditionNotRunning)
	select {
	case err = <-errChan:
		return e.onContainerWaitError(ctx, containerID, err)
	case <-containerWaitCtx.Done():
		return e.onContainerWaitError(ctx, containerID, containerWaitCtx.Err())
	case data := <-dataChan:
		return e.onContainerWaitData(ctx, data, containerAttachResponse)
	}
}

func (e *executeLogic) onContainerWaitData(
	_ context.Context,
	data container.WaitResponse,
	attachResponse types.HijackedResponse,
) (ExecuteOutput, error) {
	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)
	_, err := stdcopy.StdCopy(stdoutBuffer, stderrBuffer, attachResponse.Reader)
	if err != nil {
		e.logger.With(zap.Error(err)).Error("failed to read from stdout and stderr of container")
		return ExecuteOutput{}, err
	}

	stdOut := utils.TrimSpaceRight(stdoutBuffer.String())
	stdErr := utils.TrimSpaceRight(stderrBuffer.String())

	switch data.StatusCode {
	case 0:
		output := ExecuteOutput{
			Stdout: stdOut,
			Stderr: stdErr,
		}
		e.logger.With(zap.Any("output", output)).Info("test case run successfully")
		return output, nil
	case statusCodeTimeLimitExceeded:
		output := ExecuteOutput{
			ReturnCode:        int(data.StatusCode),
			TimeLimitExceeded: true,
			Stdout:            stdOut,
			Stderr:            stdErr,
		}
		e.logger.With(zap.Any("output", output)).Info("time limit exceeded")
		return output, nil
	case statusCodeMemoryLimitExceeded:
		output := ExecuteOutput{
			ReturnCode:          int(data.StatusCode),
			MemoryLimitExceeded: true,
			Stdout:              stdOut,
			Stderr:              stdErr,
		}
		e.logger.With(zap.Any("output", output)).Info("memory limit exceeded")
		return output, nil
	default:
		output := ExecuteOutput{
			ReturnCode: int(data.StatusCode),
			Stdout:     stdOut,
			Stderr:     stdErr,
		}

		e.logger.With(zap.Any("execute_output", output)).Info("test case run failed")
		return output, nil
	}
}

func (e *executeLogic) onContainerWaitError(
	_ context.Context,
	containerID string,
	err error,
) (ExecuteOutput, error) {
	if errors.Is(err, context.DeadlineExceeded) {
		e.logger.Info("test case execution failed: context deadline exceeded")
		return ExecuteOutput{TimeLimitExceeded: true}, nil
	}

	e.logger.With(zap.String("container_id", containerID)).With(zap.Error(err)).Error("failed to wait for container")
	return ExecuteOutput{}, err
}

func (e *executeLogic) getExecuteCommand(timeout time.Duration, programFilePath string) strslice.StrSlice {
	executeTemplate := make([]string, len(e.executeConfig.CommandTemplate))
	for i := range e.executeConfig.CommandTemplate {
		switch e.executeConfig.CommandTemplate[i] {
		case executeTimeoutPlaceholder:
			executeTemplate[i] = fmt.Sprintf("%d", int64(timeout.Seconds()))
		case compileProgramFilePathPlaceholder:
			executeTemplate[i] = programFilePath
		default:
			executeTemplate[i] = e.executeConfig.CommandTemplate[i]
		}
	}
	return strslice.StrSlice(executeTemplate)
}

func (e *executeLogic) pullImage() error {
	e.logger.Info("pulling image")
	_, err := e.dockerClient.ImagePull(context.Background(), e.executeConfig.Image, image.PullOptions{})
	if err != nil {
		e.logger.With(zap.Error(err)).Error("failed to pull image")
		return err
	}

	e.logger.Info("execute image pulled successfully")
	return nil
}
