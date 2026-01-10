package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"usermgmt/internal/errorx"
	"usermgmt/internal/logic/user"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
	"usermgmt/pkg/response"
)

func UpdateProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateProfileRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, r, http.StatusBadRequest, errorx.ErrValidation.Code, err.Error(), nil)
			return
		}

		if err := svcCtx.Validator.StructCtx(r.Context(), req); err != nil {
			appErr := errorx.FromValidationError(err)
			response.Error(w, r, appErr.Status, appErr.Code, appErr.Message, appErr.Details)
			return
		}

		logic := user.NewUpdateProfileLogic(r.Context(), svcCtx)
		resp, err := logic.Update(&req)
		if err != nil {
			handleError(w, r, err)
			return
		}

		response.Success(w, r, resp)
	}
}
