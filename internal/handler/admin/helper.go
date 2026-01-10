package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"usermgmt/internal/errorx"
	"usermgmt/pkg/response"
)

// handleError unifies error responses for admin handlers.
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	if appErr, ok := err.(*errorx.AppError); ok {
		response.Error(w, r, appErr.Status, appErr.Code, appErr.Message, appErr.Details)
		return
	}
	response.Error(w, r, errorx.ErrInternal.Status, errorx.ErrInternal.Code, errorx.ErrInternal.Message, nil)
}

func parseUserIDFromPath(r *http.Request) (uint64, error) {
	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	for i := 0; i < len(segments)-1; i++ {
		if segments[i] == "users" {
			return strconv.ParseUint(segments[i+1], 10, 64)
		}
	}
	return 0, errors.New("用户ID缺失")
}
