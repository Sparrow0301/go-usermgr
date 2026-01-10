package user

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"usermgmt/internal/errorx"
	"usermgmt/internal/logic/common"
	"usermgmt/internal/model"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
	"usermgmt/pkg/contextx"
)

// UpdateProfileLogic handles email/name changes for current user.
type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateProfileLogic constructor.
func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) Update(req *types.UpdateProfileRequest) (*types.ProfileResponse, error) {
	claims := contextx.MustGetClaims(l.ctx)
	if claims == nil {
		return nil, errorx.ErrInvalidCredentials
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	fullName := strings.TrimSpace(req.FullName)

	db := l.svcCtx.DB.WithContext(l.ctx)

	var count int64
	if err := db.Model(&model.User{}).
		Where("email = ? AND id <> ?", email, claims.UserID).
		Count(&count).Error; err != nil {
		l.Errorf("check email unique failed: %v", err)
		return nil, errorx.ErrInternal
	}
	if count > 0 {
		return nil, errorx.ErrUserExists
	}

	if err := db.Model(&model.User{}).
		Where("id = ?", claims.UserID).
		Updates(map[string]interface{}{
			"email":     email,
			"full_name": fullName,
		}).Error; err != nil {
		l.Errorf("update profile failed: %v", err)
		return nil, errorx.ErrInternal
	}

	var user model.User
	if err := db.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
		l.Errorf("load updated user failed: %v", err)
		return nil, errorx.ErrInternal
	}

	dto := common.ToUserDTO(&user)
	return &types.ProfileResponse{User: dto}, nil
}
