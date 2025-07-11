package cmd

import (
	"context"
	"google.golang.org/grpc/reflection"
	"middle_srv/app/rpc/internal/controller/gate"
	"middle_srv/internal/service"

	//_ "middle_srv/internal/boot"

	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start rpc server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := service.SrvRegister().Register(ctx)
			gate.Register(s)
			
			reflection.Register(s.Server)

			s.Run()
			return nil
		},
	}
)
