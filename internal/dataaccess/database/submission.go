package database

import (
	"context"
	"errors"

	"github.com/maxuanquang/ojs/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrSubmissionNotFound = errors.New("submission not found")
)

type Submission struct {
	ID          uint64 `gorm:"column:id;primaryKey"`
	OfProblemID uint64 `gorm:"column:of_problem_id"`
	AuthorID    uint64 `gorm:"column:author_id"`
	Content     string `gorm:"column:content"`
	Language    string `gorm:"column:language"`
	Status      int8   `gorm:"column:status"`
	Result      int8   `gorm:"column:result"`
}

type SubmissionDataAccessor interface {
	CreateSubmission(ctx context.Context, submission Submission) (Submission, error)
	GetSubmissionByID(ctx context.Context, id uint64) (Submission, error)
	GetSubmissionList(ctx context.Context, offset, limit uint64) ([]Submission, error)
	GetSubmissionCount(ctx context.Context) (uint64, error)
	GetProblemSubmissionList(ctx context.Context, problemID, offset, limit uint64) ([]Submission, error)
	GetProblemSubmissionCount(ctx context.Context, problemID uint64) (uint64, error)
	GetAccountProblemSubmissionList(ctx context.Context, accountID, problemID, offset, limit uint64) ([]Submission, error)
	GetAccountProblemSubmissionCount(ctx context.Context, accountID, problemID uint64) (uint64, error)
	UpdateSubmission(ctx context.Context, submission Submission) (Submission, error)
	DeleteSubmission(ctx context.Context, id uint64) error
	WithDatabaseTransaction(database Database) SubmissionDataAccessor
}

func NewSubmissionDataAccessor(database Database, logger *zap.Logger) SubmissionDataAccessor {
	return &submissionDataAccessor{
		database: database,
		logger:   logger,
	}
}

type submissionDataAccessor struct {
	database Database
	logger   *zap.Logger
}

// CreateSubmission implements SubmissionDataAccessor.
func (s *submissionDataAccessor) CreateSubmission(ctx context.Context, submission Submission) (Submission, error) {
	createdSubmission := Submission{
		OfProblemID: submission.OfProblemID,
		AuthorID:    submission.AuthorID,
		Content:     submission.Content,
		Language:    submission.Language,
		Status:      submission.Status,
		Result:      submission.Result,
	}
	result := s.database.Create(&createdSubmission)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, s.logger).With(zap.Any("submission", submission))
		logger.Error("error creating submission", zap.Error(result.Error))
		return Submission{}, result.Error
	}

	return createdSubmission, nil
}

// GetSubmissionByID implements SubmissionDataAccessor.
func (s *submissionDataAccessor) GetSubmissionByID(ctx context.Context, id uint64) (Submission, error) {
	var foundSubmission Submission
	result := s.database.First(&foundSubmission, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Submission{}, ErrSubmissionNotFound
		}

		logger := utils.LoggerWithContext(ctx, s.logger).With(zap.Uint64("submission_id", id))
		logger.Error("error getting submission", zap.Error(result.Error))
		return Submission{}, result.Error
	}

	return foundSubmission, nil
}

// GetSubmissionList implements SubmissionDataAccessor.
func (s *submissionDataAccessor) GetSubmissionList(ctx context.Context, offset, limit uint64) ([]Submission, error) {
	var submissions []Submission
	result := s.database.Offset(int(offset)).Limit(int(limit)).Find(&submissions)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, s.logger)
		logger.Error("error getting submission list", zap.Error(result.Error))
		return nil, result.Error
	}

	return submissions, nil
}

func (s *submissionDataAccessor) GetSubmissionCount(ctx context.Context) (uint64, error) {
	var count int64
	if err := s.database.Model(&Submission{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return uint64(count), nil
}

// GetProblemSubmissionList implements SubmissionDataAccessor.
func (s *submissionDataAccessor) GetProblemSubmissionList(ctx context.Context, problemID, offset, limit uint64) ([]Submission, error) {
	var submissions []Submission
	result := s.database.Where("of_problem_id = ?", problemID).Offset(int(offset)).Limit(int(limit)).Find(&submissions)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, s.logger).With(zap.Uint64("problem_id", problemID))
		logger.Error("error getting submission list for problem", zap.Error(result.Error))
		return nil, result.Error
	}

	return submissions, nil
}

func (s *submissionDataAccessor) GetProblemSubmissionCount(ctx context.Context, problemID uint64) (uint64, error) {
	var count int64
	if err := s.database.Model(&Submission{}).Where("of_problem_id = ?", problemID).Count(&count).Error; err != nil {
		return 0, err
	}
	return uint64(count), nil
}

// GetAccountProblemSubmissionList implements SubmissionDataAccessor.
func (s *submissionDataAccessor) GetAccountProblemSubmissionList(ctx context.Context, accountID, problemID, offset, limit uint64) ([]Submission, error) {
	var submissions []Submission
	result := s.database.Where("author_id = ? AND of_problem_id = ?", accountID, problemID).Offset(int(offset)).Limit(int(limit)).Find(&submissions)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, s.logger).With(zap.Uint64("account_id", accountID), zap.Uint64("problem_id", problemID))
		logger.Error("error getting submission list for account and problem", zap.Error(result.Error))
		return nil, result.Error
	}

	return submissions, nil
}

func (s *submissionDataAccessor) GetAccountProblemSubmissionCount(ctx context.Context, accountID, problemID uint64) (uint64, error) {
	var count int64
	if err := s.database.Model(&Submission{}).Where("author_id = ? AND of_problem_id = ?", accountID, problemID).Count(&count).Error; err != nil {
		return 0, err
	}
	return uint64(count), nil
}

// UpdateSubmission implements SubmissionDataAccessor.
func (s *submissionDataAccessor) UpdateSubmission(ctx context.Context, submission Submission) (Submission, error) {
	result := s.database.Save(&submission)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, s.logger).With(zap.Any("submission", submission))
		logger.Error("error updating submission", zap.Error(result.Error))
		return Submission{}, result.Error
	}

	return submission, nil
}

// DeleteSubmission implements SubmissionDataAccessor.
func (s *submissionDataAccessor) DeleteSubmission(ctx context.Context, id uint64) error {
	result := s.database.Delete(&Submission{}, id)
	if result.Error != nil {
		logger := utils.LoggerWithContext(ctx, s.logger).With(zap.Uint64("submission_id", id))
		logger.Error("error deleting submission", zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// WithDatabaseTransaction implements SubmissionDataAccessor.
func (s *submissionDataAccessor) WithDatabaseTransaction(database Database) SubmissionDataAccessor {
	return &submissionDataAccessor{
		database: database,
		logger:   s.logger,
	}
}
