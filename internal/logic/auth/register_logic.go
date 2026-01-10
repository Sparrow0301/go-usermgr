package auth

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"usermgmt/internal/errorx"
	"usermgmt/internal/logic/common"
	"usermgmt/internal/model"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
	"usermgmt/pkg/security"
)

// RegisterLogic handles user sign-up, including validation and hashing.
type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic constructs the logic layer with request context.
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (*types.UserDTO, error) {
	db := l.svcCtx.DB.WithContext(l.ctx)

	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	fullName := strings.TrimSpace(req.FullName)

	var count int64
	if err := db.Model(&model.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error; err != nil {
		l.Errorf("check user exists failed: %v", err)
		return nil, errorx.ErrInternal
	}

	if count > 0 {
		return nil, errorx.ErrUserExists
	}

	hash, err := security.HashPassword(req.Password, l.svcCtx.Config.Password.BcryptCost)
	if err != nil {
		l.Errorf("hash password failed: %v", err)
		return nil, errorx.ErrInternal
	}

	user := model.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		FullName:     fullName,
		Status:       model.UserStatusEnabled,
	}

	if err := db.Create(&user).Error; err != nil {
		l.Errorf("create user failed: %v", err)
		return nil, errorx.ErrInternal
	}

	dto := common.ToUserDTO(&user)
	return &dto, nil
}
