package logic

import (
	"context"

	"github.com/mikespook/gorbac"
	"go.uber.org/zap"
)

const (
	AccountRoleAdmin         = "admin"
	AccountRoleProblemSetter = "problem_setter"
	AccountRoleContester     = "contester"
	AccountRoleWorker        = "worker"
)

var (
	PermissionAccountsReadSelf  = gorbac.NewLayerPermission("accounts:read:self")
	PermissionAccountsWriteSelf = gorbac.NewLayerPermission("accounts:write:self")
	PermissionAccountsReadAll   = gorbac.NewLayerPermission("accounts:read")
	PermissionAccountsWriteAll  = gorbac.NewLayerPermission("accounts:write")

	PermissionTestCasesReadSelf  = gorbac.NewLayerPermission("testcases:read:self")
	PermissionTestCasesWriteSelf = gorbac.NewLayerPermission("testcases:write:self")
	PermissionTestCasesReadAll   = gorbac.NewLayerPermission("testcases:read")
	PermissionTestCasesWriteAll  = gorbac.NewLayerPermission("testcases:write")

	PermissionSubmissionsReadSelf  = gorbac.NewLayerPermission("submissions:read:self")
	PermissionSubmissionsWriteSelf = gorbac.NewLayerPermission("submissions:write:self")
	PermissionSubmissionsReadAll   = gorbac.NewLayerPermission("submissions:read")
	PermissionSubmissionsWriteAll  = gorbac.NewLayerPermission("submissions:write")

	PermissionProblemsReadSelf  = gorbac.NewLayerPermission("problems:read:self")
	PermissionProblemsWriteSelf = gorbac.NewLayerPermission("problems:write:self")
	PermissionProblemsReadAll   = gorbac.NewLayerPermission("problems:read")
	PermissionProblemsWriteAll  = gorbac.NewLayerPermission("problems:write")
)

type RoleLogic interface {
	AccountHasPermission(ctx context.Context, accountRole string, permissions ...gorbac.Permission) (bool, error)
}

func NewRoleLogic(logger *zap.Logger) RoleLogic {
	rbac := gorbac.New()

	roleAdmin := gorbac.NewStdRole(AccountRoleAdmin)
	roleAdmin.Assign(PermissionAccountsReadAll)
	roleAdmin.Assign(PermissionAccountsWriteAll)
	roleAdmin.Assign(PermissionProblemsReadAll)
	roleAdmin.Assign(PermissionProblemsWriteAll)
	roleAdmin.Assign(PermissionTestCasesReadAll)
	roleAdmin.Assign(PermissionTestCasesWriteAll)
	roleAdmin.Assign(PermissionSubmissionsReadAll)
	roleAdmin.Assign(PermissionSubmissionsWriteAll)

	roleProblemSetter := gorbac.NewStdRole(AccountRoleProblemSetter)
	roleProblemSetter.Assign(PermissionAccountsReadAll)
	roleProblemSetter.Assign(PermissionAccountsWriteSelf)
	roleProblemSetter.Assign(PermissionProblemsReadAll)
	roleProblemSetter.Assign(PermissionProblemsWriteSelf)
	roleProblemSetter.Assign(PermissionTestCasesReadAll)
	roleProblemSetter.Assign(PermissionTestCasesWriteSelf)
	roleProblemSetter.Assign(PermissionSubmissionsReadAll)
	roleProblemSetter.Assign(PermissionSubmissionsWriteSelf)

	roleContester := gorbac.NewStdRole(AccountRoleContester)
	roleContester.Assign(PermissionAccountsReadSelf)
	roleContester.Assign(PermissionAccountsWriteSelf)
	roleContester.Assign(PermissionProblemsReadAll)
	roleContester.Assign(PermissionSubmissionsReadAll)
	roleContester.Assign(PermissionSubmissionsWriteSelf)

	roleWorker := gorbac.NewStdRole(AccountRoleWorker)
	roleWorker.Assign(PermissionProblemsReadAll)
	roleWorker.Assign(PermissionTestCasesReadAll)
	roleWorker.Assign(PermissionSubmissionsReadAll)
	roleWorker.Assign(PermissionSubmissionsWriteAll)

	rbac.Add(roleAdmin)
	rbac.Add(roleProblemSetter)
	rbac.Add(roleContester)
	rbac.Add(roleWorker)

	return &roleLogic{
		rbac:   rbac,
		logger: logger,
	}
}

type roleLogic struct {
	rbac   *gorbac.RBAC
	logger *zap.Logger
}

// AccountHasPermission returns true if account role has any of required permissions.
func (r *roleLogic) AccountHasPermission(
	ctx context.Context,
	accountRole string,
	requiredPermissions ...gorbac.Permission,
) (bool, error) {
	accountRBACRole, _, err := r.rbac.Get(accountRole)
	if err != nil {
		r.logger.With(zap.Error(err)).Error("failed to get account role: role does not exist")
		return false, err
	}

	for i := range requiredPermissions {
		if accountRBACRole.Permit(requiredPermissions[i]) {
			return true, nil
		}
	}

	return false, nil
}
