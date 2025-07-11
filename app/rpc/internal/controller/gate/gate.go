package gate

import (
	"context"
	v1 "middle_srv/app/rpc/api/gate/v1"
	"middle_srv/internal/service"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedGateServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterGateServer(s.Server, &Controller{})
}

func (*Controller) Call(ctx context.Context, req *v1.CallRequest) (res *v1.CallReply, err error) {
	call, err := service.Gate().Call(ctx, req)
	if err != nil {
		return nil, err
	}
	return call, nil
}
