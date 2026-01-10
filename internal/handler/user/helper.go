package user

import (
	"net/http"

	"usermgmt/internal/errorx"
	"usermgmt/pkg/response"
)

// handleError unifies error responses for user handlers.
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
