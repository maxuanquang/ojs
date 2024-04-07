package logic

import (
	"context"
	"errors"

	"github.com/maxuanquang/ojs/internal/dataaccess/database"
	"github.com/maxuanquang/ojs/internal/dataaccess/mq/producer"
	"github.com/maxuanquang/ojs/internal/generated/grpc/ojs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SubmissionLogic interface {
	CreateSubmission(ctx context.Context, in CreateSubmissionInput) (CreateSubmissionOutput, error)
	GetSubmission(ctx context.Context, in GetSubmissionInput) (GetSubmissionOutput, error)
	GetSubmissionList(ctx context.Context, in GetSubmissionListInput) (GetSubmissionListOutput, error)
	GetAccountProblemSubmissionList(ctx context.Context, in GetAccountProblemSubmissionListInput) (GetAccountProblemSubmissionListOutput, error)
	GetProblemSubmissionList(ctx context.Context, in GetProblemSubmissionListInput) (GetProblemSubmissionListOutput, error)

	ExecuteSubmission(ctx context.Context, in ExecuteSubmissionInput) error
}

func NewSubmissionLogic(
	logger *zap.Logger,
	accountDataAccessor database.AccountDataAccessor,
	problemDataAccessor database.ProblemDataAccessor,
	submissionDataAccessor database.SubmissionDataAccessor,
	testCaseDataAccessor database.TestCaseDataAccessor,
	tokenLogic TokenLogic,
	judgeLogic JudgeLogic,
	submissionCreatedProducer producer.SubmissionCreatedProducer,
	database database.Database,
) SubmissionLogic {
	return &submissionLogic{
		logger:                    logger,
		accountDataAccessor:       accountDataAccessor,
		problemDataAccessor:       problemDataAccessor,
		submissionDataAccessor:    submissionDataAccessor,
		testCaseDataAccessor:      testCaseDataAccessor,
		tokenLogic:                tokenLogic,
		judgeLogic:                judgeLogic,
		submissionCreatedProducer: submissionCreatedProducer,
		database:                  database,
	}
}

type submissionLogic struct {
	logger                    *zap.Logger
	accountDataAccessor       database.AccountDataAccessor
	problemDataAccessor       database.ProblemDataAccessor
	submissionDataAccessor    database.SubmissionDataAccessor
	testCaseDataAccessor      database.TestCaseDataAccessor
	tokenLogic                TokenLogic
	judgeLogic                JudgeLogic
	submissionCreatedProducer producer.SubmissionCreatedProducer
	database                  database.Database
}

func (p *submissionLogic) CreateSubmission(ctx context.Context, in CreateSubmissionInput) (CreateSubmissionOutput, error) {
	var (
		err               error
		createdSubmission database.Submission
		txErr             error
		accountID         uint64
	)

	// Verify token
	accountID, _, _, _, err = p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		p.logger.Error("Failed to verify token", zap.Error(err))
		return CreateSubmissionOutput{}, err
	}

	// Create submission in the database
	txErr = p.database.Transaction(func(tx *gorm.DB) error {
		createdSubmission, err = p.submissionDataAccessor.WithDatabaseTransaction(tx).CreateSubmission(ctx, database.Submission{
			OfProblemID: in.OfProblemID,
			AuthorID:    accountID,
			Content:     in.Content,
			Language:    in.Language,
			Status:      int8(ojs.SubmissionStatus_Submitted),
		})
		if err != nil {
			p.logger.Error("failed to create submission", zap.Error(err))
			return err
		}

		return nil
	})
	if txErr != nil {
		p.logger.Error("create submission transaction failed", zap.Error(err))
		return CreateSubmissionOutput{}, err
	}

	// Produce a message to the submission created queue
	err = p.submissionCreatedProducer.Produce(ctx, createdSubmission.ID)
	if err != nil {
		p.logger.Error("failed to send message to submission created queue", zap.Error(err))
	}

	return CreateSubmissionOutput{
		Submission: Submission{
			ID:          createdSubmission.ID,
			OfProblemID: createdSubmission.OfProblemID,
			AuthorID:    createdSubmission.AuthorID,
			Content:     createdSubmission.Content,
			Language:    createdSubmission.Language,
			Status:      ojs.SubmissionStatus(createdSubmission.Status),
			Result:      ojs.SubmissionResult(createdSubmission.Result),
		},
	}, nil
}

func (p *submissionLogic) GetSubmission(ctx context.Context, in GetSubmissionInput) (GetSubmissionOutput, error) {
	// Retrieve submission from the database
	submission, err := p.submissionDataAccessor.GetSubmissionByID(ctx, in.ID)
	if err != nil {
		p.logger.Error("Failed to get submission", zap.Error(err))
		return GetSubmissionOutput{}, err
	}

	return GetSubmissionOutput{
		Submission: Submission{
			ID:          submission.ID,
			OfProblemID: submission.OfProblemID,
			AuthorID:    submission.AuthorID,
			Content:     submission.Content,
			Language:    submission.Language,
			Status:      ojs.SubmissionStatus(submission.Status),
			Result:      ojs.SubmissionResult(submission.Result),
		},
	}, nil
}

func (p *submissionLogic) GetSubmissionList(ctx context.Context, in GetSubmissionListInput) (GetSubmissionListOutput, error) {
	// Retrieve submission list from the database
	submissions, err := p.submissionDataAccessor.GetSubmissionList(ctx, in.Offset, in.Limit)
	if err != nil {
		p.logger.Error("Failed to get submission list", zap.Error(err))
		return GetSubmissionListOutput{}, err
	}

	var submissionList []Submission
	for _, s := range submissions {
		submissionList = append(submissionList, Submission{
			ID:          s.ID,
			OfProblemID: s.OfProblemID,
			AuthorID:    s.AuthorID,
			Content:     s.Content,
			Language:    s.Language,
			Status:      ojs.SubmissionStatus(s.Status),
			Result:      ojs.SubmissionResult(s.Result),
		})
	}

	totalSubmissionsCount, err := p.submissionDataAccessor.GetSubmissionCount(ctx)
	if err != nil {
		p.logger.Error("Failed to get submission count", zap.Error(err))
		return GetSubmissionListOutput{}, err
	}

	return GetSubmissionListOutput{
		Submissions:           submissionList,
		TotalSubmissionsCount: totalSubmissionsCount,
	}, nil
}

func (p *submissionLogic) GetAccountProblemSubmissionList(ctx context.Context, in GetAccountProblemSubmissionListInput) (GetAccountProblemSubmissionListOutput, error) {
	// Verify token
	accountID, _, _, _, err := p.tokenLogic.VerifyTokenString(ctx, in.Token)
	if err != nil {
		p.logger.Error("Failed to verify token", zap.Error(err))
		return GetAccountProblemSubmissionListOutput{}, err
	}

	// Retrieve account's submission list for a specific problem from the database
	submissions, err := p.submissionDataAccessor.GetAccountProblemSubmissionList(ctx, accountID, in.OfProblemID, in.Offset, in.Limit)
	if err != nil {
		p.logger.Error("Failed to get account's problem submission list", zap.Error(err))
		return GetAccountProblemSubmissionListOutput{}, err
	}

	var submissionList []Submission
	for _, s := range submissions {
		submissionList = append(submissionList, Submission{
			ID:          s.ID,
			OfProblemID: s.OfProblemID,
			AuthorID:    s.AuthorID,
			Content:     s.Content,
			Language:    s.Language,
			Status:      ojs.SubmissionStatus(s.Status),
			Result:      ojs.SubmissionResult(s.Result),
		})
	}

	totalSubmissionsCount, err := p.submissionDataAccessor.GetAccountProblemSubmissionCount(ctx, accountID, in.OfProblemID)
	if err != nil {
		p.logger.Error("Failed to get account's problem submission count", zap.Error(err))
		return GetAccountProblemSubmissionListOutput{}, err
	}

	return GetAccountProblemSubmissionListOutput{
		Submissions:           submissionList,
		TotalSubmissionsCount: totalSubmissionsCount,
	}, nil
}

func (p *submissionLogic) GetProblemSubmissionList(ctx context.Context, in GetProblemSubmissionListInput) (GetProblemSubmissionListOutput, error) {
	// Retrieve problem's submission list from the database
	submissions, err := p.submissionDataAccessor.GetProblemSubmissionList(ctx, in.OfProblemID, in.Offset, in.Limit)
	if err != nil {
		p.logger.Error("Failed to get problem submission list", zap.Error(err))
		return GetProblemSubmissionListOutput{}, err
	}

	var submissionList []Submission
	for _, s := range submissions {
		submissionList = append(submissionList, Submission{
			ID:          s.ID,
			OfProblemID: s.OfProblemID,
			AuthorID:    s.AuthorID,
			Content:     s.Content,
			Language:    s.Language,
			Status:      ojs.SubmissionStatus(s.Status),
			Result:      ojs.SubmissionResult(s.Result),
		})
	}

	totalSubmissionsCount, err := p.submissionDataAccessor.GetProblemSubmissionCount(ctx, in.OfProblemID)
	if err != nil {
		p.logger.Error("Failed to get problem submission count", zap.Error(err))
		return GetProblemSubmissionListOutput{}, err
	}

	return GetProblemSubmissionListOutput{
		Submissions:           submissionList,
		TotalSubmissionsCount: totalSubmissionsCount,
	}, nil
}

// ExecuteSubmission implements SubmissionLogic.
func (s *submissionLogic) ExecuteSubmission(ctx context.Context, in ExecuteSubmissionInput) error {
	var (
		err        error
		txErr      error
		submission database.Submission
	)

	txErr = s.database.Transaction(func(tx *gorm.DB) error {
		submission, err = s.submissionDataAccessor.WithDatabaseTransaction(tx).GetSubmissionByID(ctx, in.ID)
		if err != nil {
			s.logger.Error("Failed to get submission", zap.Error(err))
			return err
		}

		if submission.Status != int8(ojs.SubmissionStatus_Submitted) {
			s.logger.Error("Submission is not submitted", zap.Uint64("submission_id", submission.ID))
			return errors.New("submission is not submitted")
		}

		submission.Status = int8(ojs.SubmissionStatus_Executing)
		_, err = s.submissionDataAccessor.WithDatabaseTransaction(tx).UpdateSubmission(ctx, submission)
		if err != nil {
			s.logger.Error("Failed to update submission", zap.Error(err))
			return err
		}

		return nil
	})
	if txErr != nil {
		s.logger.Error("execute submission transaction failed", zap.Error(txErr))
		return txErr
	}

	result, err := s.judgeLogic.Judge(
		ctx,
		databaseSubmissionToLogicSubmission(submission),
	)
	if err != nil {
		s.logger.Error("Failed to judge submission", zap.Error(err))
	}

	// Update submission result and status in the database
	submission.Result = int8(result)
	submission.Status = int8(ojs.SubmissionStatus_Finished)
	s.submissionDataAccessor.UpdateSubmission(
		ctx,
		submission,
	)

	return nil
}

func databaseSubmissionToLogicSubmission(submission database.Submission) Submission {
	return Submission{
		ID:          submission.ID,
		OfProblemID: submission.OfProblemID,
		AuthorID:    submission.AuthorID,
		Content:     submission.Content,
		Language:    submission.Language,
		Status:      ojs.SubmissionStatus(submission.Status),
		Result:      ojs.SubmissionResult(submission.Result),
	}
}

type CreateSubmissionInput struct {
	Token       string
	OfProblemID uint64
	Content     string
	Language    string
}

type Submission struct {
	ID          uint64
	AuthorID    uint64
	OfProblemID uint64
	Content     string
	Language    string
	Status      ojs.SubmissionStatus
	Result      ojs.SubmissionResult
}

type CreateSubmissionOutput struct {
	Submission Submission
}

type GetSubmissionInput struct {
	ID uint64
}

type GetSubmissionOutput struct {
	Submission Submission
}

type GetSubmissionListInput struct {
	Offset uint64
	Limit  uint64
}

type GetSubmissionListOutput struct {
	Submissions           []Submission
	TotalSubmissionsCount uint64
}

type GetAccountProblemSubmissionListInput struct {
	Token       string
	OfProblemID uint64
	Offset      uint64
	Limit       uint64
}

type GetAccountProblemSubmissionListOutput struct {
	Submissions           []Submission
	TotalSubmissionsCount uint64
}

type GetProblemSubmissionListInput struct {
	Token       string
	OfProblemID uint64
	Offset      uint64
	Limit       uint64
	IsHidden    bool
}

type GetProblemSubmissionListOutput struct {
	Submissions           []Submission
	TotalSubmissionsCount uint64
}

type ExecuteSubmissionInput struct {
	ID uint64
}
