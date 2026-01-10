package user

import (
	"net/http"

	"usermgmt/internal/logic/user"
	"usermgmt/internal/svc"
	"usermgmt/pkg/response"
)

func ProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logic := user.NewProfileLogic(r.Context(), svcCtx)
		resp, err := logic.Profile()
		if err != nil {
			handleError(w, r, err)
			return
		}

		response.Success(w, r, resp)
	}
}
