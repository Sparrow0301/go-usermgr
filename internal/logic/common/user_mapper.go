package common

import (
	"usermgmt/internal/model"
	"usermgmt/internal/types"
)

// ToUserDTO maps model.User to API DTO.
func ToUserDTO(user *model.User) types.UserDTO {
	if user == nil {
		return types.UserDTO{}
	}
	return types.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Status:    user.Status,
		Roles:     ExtractRoleNames(user.Roles),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ExtractRoleNames returns role name slice.
func ExtractRoleNames(roles []model.Role) []string {
	if len(roles) == 0 {
		return []string{}
	}
	result := make([]string, 0, len(roles))
	for _, role := range roles {
		result = append(result, role.Name)
	}
	return result
}
