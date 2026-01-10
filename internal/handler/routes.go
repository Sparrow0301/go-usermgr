package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"usermgmt/internal/handler/admin"
	"usermgmt/internal/handler/auth"
	userhandler "usermgmt/internal/handler/user"
	"usermgmt/internal/svc"
)

// RegisterHandlers wires up all HTTP routes.
func RegisterHandlers(server *rest.Server, ctx *svc.ServiceContext) {
	authGroup := []rest.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/auth/register",
			Handler: auth.RegisterHandler(ctx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/auth/login",
			Handler: auth.LoginHandler(ctx),
		},
	}

	userGroup := []rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/me",
			Handler: ctx.AuthMiddleware(userhandler.ProfileHandler(ctx)),
		},
		{
			Method:  http.MethodPut,
			Path:    "/api/v1/me",
			Handler: ctx.AuthMiddleware(userhandler.UpdateProfileHandler(ctx)),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/me/password",
			Handler: ctx.AuthMiddleware(userhandler.ChangePasswordHandler(ctx)),
		},
	}

	adminGroup := []rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/v1/admin/users",
			Handler: ctx.AuthMiddleware(ctx.RoleGuard("admin")(admin.ListUsersHandler(ctx))),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/api/v1/admin/users/:id/status",
			Handler: ctx.AuthMiddleware(ctx.RoleGuard("admin")(admin.UpdateUserStatusHandler(ctx))),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/admin/users/:id/roles",
			Handler: ctx.AuthMiddleware(ctx.RoleGuard("admin")(admin.AssignRolesHandler(ctx))),
		},
	}

	server.AddRoutes(authGroup)
	server.AddRoutes(userGroup)
	server.AddRoutes(adminGroup)
}
