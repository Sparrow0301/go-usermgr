package user

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"usermgmt/internal/errorx"
	"usermgmt/internal/logic/common"
	"usermgmt/internal/model"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
	"usermgmt/pkg/contextx"
)

// ProfileLogic fetches current user's profile info.
type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewProfileLogic constructor.
func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) Profile() (*types.ProfileResponse, error) {
	claims := contextx.MustGetClaims(l.ctx)
	if claims == nil {
		return nil, errorx.ErrInvalidCredentials
	}

	var user model.User
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Preload("Roles").
		First(&user, claims.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.ErrInvalidCredentials
		}
		l.Errorf("load profile failed: %v", err)
		return nil, errorx.ErrInternal
	}

	dto := common.ToUserDTO(&user)
	return &types.ProfileResponse{User: dto}, nil
}
