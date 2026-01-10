package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"usermgmt/internal/errorx"
	"usermgmt/internal/model"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
	"usermgmt/pkg/contextx"
	"usermgmt/pkg/security"
)

// ChangePasswordLogic verifies old password then updates new hash.
type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangePasswordLogic constructor.
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) Change(req *types.ChangePasswordRequest) error {
	claims := contextx.MustGetClaims(l.ctx)
	if claims == nil {
		return errorx.ErrInvalidCredentials
	}

	db := l.svcCtx.DB.WithContext(l.ctx)

	var user model.User
	if err := db.First(&user, claims.UserID).Error; err != nil {
		l.Errorf("load user failed: %v", err)
		return errorx.ErrInternal
	}

	if err := security.VerifyPassword(user.PasswordHash, req.OldPassword); err != nil {
		return errorx.ErrInvalidCredentials
	}

	hash, err := security.HashPassword(req.NewPassword, l.svcCtx.Config.Password.BcryptCost)
	if err != nil {
		l.Errorf("hash new password failed: %v", err)
		return errorx.ErrInternal
	}

	if err := db.Model(&model.User{}).
		Where("id = ?", user.ID).
		Update("password_hash", hash).Error; err != nil {
		l.Errorf("update password failed: %v", err)
		return errorx.ErrInternal
	}

	return nil
}
