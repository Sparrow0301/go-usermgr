package admin

import (
	"context"
	"math"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"usermgmt/internal/errorx"
	"usermgmt/internal/logic/common"
	"usermgmt/internal/model"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
)

// ListUsersLogic encapsulates pagination & filtering of users.
type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListUsersLogic constructor.
func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) List(req *types.ListUsersRequest) (*types.ListUsersResponse, error) {
	db := l.svcCtx.DB.WithContext(l.ctx)

	page := req.Page
	if page < 1 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = l.svcCtx.Config.Pagination.DefaultPageSize
		if pageSize <= 0 {
			pageSize = 20
		}
	}
	maxSize := l.svcCtx.Config.Pagination.MaxPageSize
	if maxSize <= 0 {
		maxSize = 100
	}
	if pageSize > maxSize {
		pageSize = maxSize
	}
	offset := (page - 1) * pageSize

	baseQuery := db.Model(&model.User{})

	if status := strings.TrimSpace(req.Status); status != "" {
		baseQuery = baseQuery.Where("status = ?", status)
	}

	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		kw := "%" + strings.ToLower(keyword) + "%"
		baseQuery = baseQuery.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(full_name) LIKE ?", kw, kw, kw)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		l.Errorf("count users failed: %v", err)
		return nil, errorx.ErrInternal
	}

	var users []model.User
	if err := baseQuery.
		Preload("Roles").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		l.Errorf("list users failed: %v", err)
		return nil, errorx.ErrInternal
	}

	data := make([]types.UserDTO, 0, len(users))
	for _, user := range users {
		data = append(data, common.ToUserDTO(&user))
	}

	var totalPages int
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	return &types.ListUsersResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}
