package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"usermgmt/internal/errorx"
	adminlogic "usermgmt/internal/logic/admin"
	"usermgmt/internal/svc"
	"usermgmt/internal/types"
	"usermgmt/pkg/response"
)

func UpdateUserStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := parseUserIDFromPath(r)
		if err != nil {
			response.Error(w, r, http.StatusBadRequest, errorx.ErrValidation.Code, err.Error(), nil)
			return
		}

		var req types.UpdateUserStatusRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, r, http.StatusBadRequest, errorx.ErrValidation.Code, err.Error(), nil)
			return
		}

		if err := svcCtx.Validator.StructCtx(r.Context(), req); err != nil {
			appErr := errorx.FromValidationError(err)
			response.Error(w, r, appErr.Status, appErr.Code, appErr.Message, appErr.Details)
			return
		}

		logic := adminlogic.NewUpdateUserStatusLogic(r.Context(), svcCtx)
		resp, err := logic.Update(uint(userID), &req)
		if err != nil {
			handleError(w, r, err)
			return
		}

		response.Success(w, r, resp)
	}
}
