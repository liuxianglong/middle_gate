// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gcfg"
)

type (
	ISrvRegister interface {
		Config(ctx context.Context, configPath string) (adapter gcfg.Adapter, err error)
		GetGsvcRegistry(ctx context.Context) (gsvc.Registry, error)
		Register(ctx context.Context) *grpcx.GrpcServer
	}
)

var (
	localSrvRegister ISrvRegister
)

func SrvRegister() ISrvRegister {
	if localSrvRegister == nil {
		panic("implement not found for interface ISrvRegister, forgot register?")
	}
	return localSrvRegister
}

func RegisterSrvRegister(i ISrvRegister) {
	localSrvRegister = i
}
