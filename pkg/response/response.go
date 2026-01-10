package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// JSON writes a JSON response with custom status code.
func JSON(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	httpx.WriteJsonCtx(r.Context(), w, status, payload)
}

// Error writes a standardized error response body with code and message.
func Error(w http.ResponseWriter, r *http.Request, status int, code, message string, details interface{}) {
	httpx.WriteJsonCtx(r.Context(), w, status, ErrorBody{Code: code, Message: message, Details: details})
}

// Success wraps payload with default status 200.
func Success(w http.ResponseWriter, r *http.Request, payload interface{}) {
	httpx.OkJsonCtx(r.Context(), w, payload)
}
