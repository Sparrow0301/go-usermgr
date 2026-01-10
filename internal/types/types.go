package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	FullName string `json:"fullName" validate:"required,min=2,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"accessToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	User         UserDTO   `json:"user"`
}

type UserDTO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	Status    string    `json:"status"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProfileResponse struct {
	User UserDTO `json:"user"`
}

type UpdateProfileRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullName" validate:"required,min=2,max=100"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=64"`
}

type ListUsersRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
}

type ListUsersResponse struct {
	Data       []UserDTO `json:"data"`
	Page       int       `json:"page"`
	PageSize   int       `json:"pageSize"`
	TotalItems int64     `json:"totalItems"`
	TotalPages int       `json:"totalPages"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=enabled disabled"`
}

type AssignRolesRequest struct {
	Roles []string `json:"roles" validate:"required,min=1,dive,required"`
}

type JwtClaims struct {
	jwt.RegisteredClaims
	UserID uint     `json:"userId"`
	Roles  []string `json:"roles"`
}
