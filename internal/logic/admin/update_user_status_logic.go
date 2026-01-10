package admin

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"usermgmt/internal/errorx"
	"usermgmt/internal/logic/common"
	"usermgmt/internal/model"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
)

// UpdateUserStatusLogic toggles enabled/disabled state.
type UpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateUserStatusLogic constructor.
func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserStatusLogic) Update(userID uint, req *types.UpdateUserStatusRequest) (*types.ProfileResponse, error) {
	db := l.svcCtx.DB.WithContext(l.ctx)

	result := db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("status", req.Status)
	if result.Error != nil {
		l.Errorf("update status failed: %v", result.Error)
		return nil, errorx.ErrInternal
	}
	if result.RowsAffected == 0 {
		return nil, errorx.ErrUserNotFound
	}

	var user model.User
	if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
		l.Errorf("load user after status update failed: %v", err)
		return nil, errorx.ErrInternal
	}

	dto := common.ToUserDTO(&user)
	return &types.ProfileResponse{User: dto}, nil
}
