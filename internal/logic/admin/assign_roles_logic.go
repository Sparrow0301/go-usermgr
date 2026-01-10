package admin

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"usermgmt/internal/errorx"
	"usermgmt/internal/logic/common"
	"usermgmt/internal/model"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
)

// AssignRolesLogic rebinds role memberships for a user.
type AssignRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewAssignRolesLogic constructor.
func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesLogic {
	return &AssignRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRolesLogic) Assign(userID uint, req *types.AssignRolesRequest) (*types.ProfileResponse, error) {
	db := l.svcCtx.DB.WithContext(l.ctx)

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrUserNotFound
		}
		l.Errorf("load user for role assignment failed: %v", err)
		return nil, errorx.ErrInternal
	}

	roleNames := normalizeRoles(req.Roles)
	if len(roleNames) == 0 {
		return nil, errorx.ErrValidation.WithDetails("角色列表不能为空")
	}

	var roles []model.Role
	if err := db.Where("name IN ?", roleNames).Find(&roles).Error; err != nil {
		l.Errorf("load roles failed: %v", err)
		return nil, errorx.ErrInternal
	}

	if len(roles) != len(roleNames) {
		existing := make(map[string]struct{})
		for _, role := range roles {
			existing[strings.ToLower(role.Name)] = struct{}{}
		}
		missing := make([]string, 0)
		for _, role := range roleNames {
			if _, ok := existing[strings.ToLower(role)]; !ok {
				missing = append(missing, role)
			}
		}
		return nil, errorx.ErrValidation.WithDetails(map[string]interface{}{"missingRoles": missing})
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		userRoles := make([]model.UserRole, 0, len(roles))
		for _, role := range roles {
			userRoles = append(userRoles, model.UserRole{UserID: userID, RoleID: role.ID})
		}
		if len(userRoles) > 0 {
			if err := tx.Create(&userRoles).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		l.Errorf("assign roles transaction failed: %v", err)
		return nil, errorx.ErrInternal
	}

	if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
		l.Errorf("reload user after role assignment failed: %v", err)
		return nil, errorx.ErrInternal
	}

	dto := common.ToUserDTO(&user)
	return &types.ProfileResponse{User: dto}, nil
}

func normalizeRoles(roles []string) []string {
	result := make([]string, 0, len(roles))
	seen := make(map[string]struct{})
	for _, role := range roles {
		role = strings.TrimSpace(role)
		if role == "" {
			continue
		}
		key := strings.ToLower(role)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, role)
	}
	return result
}
