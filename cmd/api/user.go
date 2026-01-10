package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"usermgmt/internal/config"
	"usermgmt/internal/handler"
	"usermgmt/internal/svc"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

// main bootstraps the REST server, wiring config, dependencies and routes.
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	svcCtx := svc.NewServiceContext(c)
	if err := svcCtx.AutoMigrate(); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	server := rest.MustNewServer(c.RestConf, rest.WithCors(c.Security.AllowOrigins...))

	defer server.Stop()

	handler.RegisterHandlers(server, svcCtx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
