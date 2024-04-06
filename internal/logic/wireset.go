package logic

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAccountLogic,
	NewHashLogic,
	NewTokenLogic,
	NewProblemLogic,
	NewSubmissionLogic,
	NewTestCaseLogic,
	NewJudgeLogic,
	NewCompileLogic,
)
