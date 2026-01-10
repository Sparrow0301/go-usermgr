package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"usermgmt/internal/errorx"
	"usermgmt/internal/logic/common"
	"usermgmt/internal/model"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
	"usermgmt/pkg/security"
)

// LoginLogic validates credentials and issues JWT tokens.
type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginLogic constructs the login logic with request context.
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	db := l.svcCtx.DB.WithContext(l.ctx)
	username := strings.TrimSpace(req.Username)

	var user model.User
	if err := db.Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.ErrInvalidCredentials
		}
		l.Errorf("query user failed: %v", err)
		return nil, errorx.ErrInternal
	}

	if user.Status == model.UserStatusDisabled {
		return nil, errorx.ErrUserDisabled
	}

	if err := security.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errorx.ErrInvalidCredentials
	}

	roleNames := common.ExtractRoleNames(user.Roles)
	accessToken, accessExpire, err := security.GenerateToken(user.ID, roleNames, l.svcCtx.Config.JWT.AccessSecret, l.svcCtx.Config.JWT.AccessExpire)
	if err != nil {
		l.Errorf("generate access token failed: %v", err)
		return nil, errorx.ErrInternal
	}

	refreshToken := ""
	if l.svcCtx.Config.JWT.RefreshExpire > 0 {
		if refreshTokenValue, _, err := security.GenerateToken(user.ID, roleNames, l.svcCtx.Config.JWT.AccessSecret, l.svcCtx.Config.JWT.RefreshExpire); err != nil {
			l.Errorf("generate refresh token failed: %v", err)
			return nil, errorx.ErrInternal
		} else {
			refreshToken = refreshTokenValue
		}
	}

	if err := db.Model(&model.User{}).
		Where("id = ?", user.ID).
		Update("last_login_at", time.Now()).Error; err != nil {
		l.Errorf("update last login failed: %v", err)
	}

	dto := common.ToUserDTO(&user)

	return &types.LoginResponse{
		AccessToken:  accessToken,
		ExpiresAt:    accessExpire,
		RefreshToken: refreshToken,
		User:         dto,
	}, nil
}
