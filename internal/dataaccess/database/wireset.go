package database

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAccountDataAccessor,
	NewAccountPasswordDataAccessor,
	NewTokenPublicKeyDataAccessor,
	NewMigrator,
	InitializeDB,
	NewProblemDataAccessor,
	NewSubmissionDataAccessor,
	NewTestCaseDataAccessor,
)
