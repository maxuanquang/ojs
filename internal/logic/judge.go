package logic

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/maxuanquang/ojs/internal/configs"
	"github.com/maxuanquang/ojs/internal/dataaccess/database"
	"github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
)

type JudgeLogic interface {
	Judge(ctx context.Context, submission Submission) (ojs.SubmissionResult, error)
}

func NewJudgeLogic(
	problemDataAccessor database.ProblemDataAccessor,
	submissionDataAccessor database.SubmissionDataAccessor,
	testCaseDataAccessor database.TestCaseDataAccessor,
	dockerClient *client.Client,
	judgeConfig configs.Judge,
	appArguments utils.Arguments,
	logger *zap.Logger,
) (JudgeLogic, error) {
	var languageToCompileLogic = make(map[string]CompileLogic)
	var languageToExecuteLogic = make(map[string]ExecuteLogic)

	for _, config := range judgeConfig.Languages {
		language := config.Value
		compileLogic, err := NewCompileLogic(
			logger,
			dockerClient,
			language,
			config.Compile,
			appArguments,
		)
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to create compile logic")
			return nil, err
		}

		languageToCompileLogic[language] = compileLogic

		//TODO: Add languageToExecuteLogic
	}

	return &judgeLogic{
		problemDataAccessor:    problemDataAccessor,
		submissionDataAccessor: submissionDataAccessor,
		testCaseDataAccessor:   testCaseDataAccessor,
		logger:                 logger,
		languageToCompileLogic: languageToCompileLogic,
		languageToExecuteLogic: languageToExecuteLogic,
	}, nil
}

type judgeLogic struct {
	problemDataAccessor    database.ProblemDataAccessor
	submissionDataAccessor database.SubmissionDataAccessor
	testCaseDataAccessor   database.TestCaseDataAccessor

	logger                 *zap.Logger
	languageToCompileLogic map[string]CompileLogic
	languageToExecuteLogic map[string]ExecuteLogic
}

// Judge implements JudgeLogic.
func (j *judgeLogic) Judge(ctx context.Context, submission Submission) (ojs.SubmissionResult, error) {
	compileLogic, ok := j.languageToCompileLogic[submission.Language]
	if !ok {
		j.logger.Error("unsupported language")
		return ojs.SubmissionResult_UnsupportedLanguage, nil
	}

	compileOutput, err := compileLogic.Compile(ctx, submission.Content)
	if err != nil {
		j.logger.With(zap.Error(err)).Error("failed to compile submission")
		return ojs.SubmissionResult_CompileError, nil
	}

	testCases, err := j.testCaseDataAccessor.GetProblemTestCaseListAll(ctx, submission.OfProblemID)
	if err != nil {
		j.logger.With(zap.Error(err)).Error("failed to get test cases")
		return ojs.SubmissionResult_UndefinedResult, err
	}

	executeLogic, ok := j.languageToExecuteLogic[submission.Language]
	if !ok {
		j.logger.Error("unsupported language")
		return ojs.SubmissionResult_UnsupportedLanguage, nil
	}

	for _, testCase := range testCases {
		output, err := executeLogic.Execute(ctx, compileOutput.ProgramFilePath, submission.Language, testCase.Input)
		if err != nil {
			return ojs.SubmissionResult_RuntimeError, nil
		}
		if output != testCase.Output {
			return ojs.SubmissionResult_WrongAnswer, nil
		}
	}

	return ojs.SubmissionResult_OK, nil
}
