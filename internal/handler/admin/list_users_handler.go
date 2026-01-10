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

func ListUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListUsersRequest
		if err := httpx.ParseForm(r, &req); err != nil {
			response.Error(w, r, http.StatusBadRequest, errorx.ErrValidation.Code, err.Error(), nil)
			return
		}

		logic := adminlogic.NewListUsersLogic(r.Context(), svcCtx)
		resp, err := logic.List(&req)
		if err != nil {
			handleError(w, r, err)
			return
		}

		response.Success(w, r, resp)
	}
}
