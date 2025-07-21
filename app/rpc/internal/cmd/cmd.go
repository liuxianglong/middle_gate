package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc/reflection"
	"middle_srv/app/rpc/internal/controller/gate"
	"middle_srv/internal/service"

	_ "middle_srv/internal/boot"

	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start rpc server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			adapter, err := service.SrvRegister().Config(ctx, "server/config")
			if err != nil {
				g.Log().Fatalf(ctx, `New consul adapter error: %+v`, err)
			}

			// Set Consul adapter as the configuration adapter
			// This enables GoFrame to use Consul for configuration management
			g.Cfg().SetAdapter(adapter)

			service.Limiter().Init(ctx)

			s := service.SrvRegister().Register(ctx)
			gate.Register(s)

			reflection.Register(s.Server)
			go service.Limiter().Lookup(ctx)
			s.Run()
			return nil
		},
	}
)
