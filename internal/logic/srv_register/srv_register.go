package srv_register

import (
	"context"
	"fmt"
	"github.com/gogf/gf/contrib/registry/consul/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc"
	"io/ioutil"
	"middle_srv/internal/service"
	"middle_srv/utility/code"

	//"net"
	"net/http"
	"time"
)

type (
	sSrvRegister struct{}
)

func init() {
	service.RegisterSrvRegister(New())
}

func New() service.ISrvRegister {
	return &sSrvRegister{}
}

func (s *sSrvRegister) checkConsul(addr string) error {
	url := fmt.Sprintf("http://%s/v1/status/leader", addr)
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("connect to %s failed: %w", addr, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("consul %s status code: %d", addr, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read consul %s response failed: %w", addr, err)
	}
	if len(body) == 0 {
		return fmt.Errorf("consul %s returned empty leader info", addr)
	}
	return nil
}

func (s *sSrvRegister) reg(ctx context.Context, address string) (gsvc.Registry, error) {
	err := s.checkConsul(address)
	if err != nil {
		return nil, err
	}
	registry, err := consul.New(consul.WithAddress(address))

	if err != nil {
		return nil, err
	}

	return registry, nil
}

func (s *sSrvRegister) GetGsvcRegistry(ctx context.Context) (gsvc.Registry, error) {
	addressVar, err := g.Cfg().Get(ctx, "consul.address")
	if err != nil {
		g.Log().Errorf(ctx, "获取consul.address配置错误，err=%v", err)
		return nil, code.CodeError.New(ctx, code.CommonConsulCfgError)
	}
	addressList := addressVar.Strings()
	addressLen := len(addressList)
	if addressLen == 0 {
		g.Log().Error(ctx, "获取consul.address配置为空")
		return nil, code.CodeError.New(ctx, code.CommonConsulCfgError)
	}
	var registry gsvc.Registry
	regSuc := false
	for i := 0; i < addressLen; i++ {
		registry, err = s.reg(ctx, addressList[i])
		if err == nil {
			regSuc = true
			g.Log().Warningf(ctx, "服务注册地址请求成功，配置%s", addressList[i])
			break
		}

		g.Log().Warningf(ctx, "服务注册地址请求失败，请检测配置%s, err=%v", addressList[i], err)
	}
	if !regSuc {
		g.Log().Error(ctx, "服务注册地址请求全部失败")
		return nil, code.CodeError.New(ctx, code.CommonConsulSrvCurlAllError)
	}
	return registry, nil
}

func (s *sSrvRegister) Register(ctx context.Context) *grpcx.GrpcServer {
	registry, err := s.GetGsvcRegistry(ctx)
	if err != nil {
		g.Log().Fatalf(ctx, "初始化注册失败, %v", err)
	}
	grpcx.Resolver.Register(registry)

	c := grpcx.Server.NewConfig()
	c.Options = append(c.Options, []grpc.ServerOption{
		grpcx.Server.ChainUnary(
			grpcx.Server.UnaryValidate,
		)}...,
	)
	grpcServer := grpcx.Server.New(c)

	return grpcServer
}
