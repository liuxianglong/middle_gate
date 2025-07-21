package srv_register

import (
	"context"
	consul "github.com/gogf/gf/contrib/config/consul/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"middle_srv/utility/code"
)

func (s *sSrvRegister) getAdapter(ctx context.Context, address, configPath string) (adapter gcfg.Adapter, err error) {
	err = s.checkConsul(address)
	if err != nil {
		return nil, err
	}
	consulConfig := api.Config{
		Address:    address,                            // Consul server address
		Scheme:     "http",                             // Connection scheme (http/https)
		Datacenter: "dc1",                              // Datacenter name
		Transport:  cleanhttp.DefaultPooledTransport(), // HTTP transport with connection pooling
		//Token:      "3f8aeba2-f1f7-42d0-b912-fcb041d4546d", // ACL token for authentication
	}
	adapter, err = consul.New(ctx, consul.Config{
		ConsulConfig: consulConfig, // Consul client configuration
		Path:         configPath,   // Configuration path in KV store
		Watch:        true,         // Enable configuration watching for updates
	})
	if err != nil {
		return nil, err
	}

	return adapter, nil
}
func (s *sSrvRegister) Config(ctx context.Context, configPath string) (adapter gcfg.Adapter, err error) {
	addressList, err := s.getConsulAddressList(ctx)
	if err != nil {
		return nil, err
	}
	addressLen := len(addressList)
	if addressLen == 0 {
		g.Log().Error(ctx, "获取consul.address配置为空")
		return nil, code.CodeError.New(ctx, code.CommonConsulCfgError)
	}
	configSuc := false
	for i := 0; i < addressLen; i++ {
		adapter, err = s.getAdapter(ctx, addressList[i], configPath)
		if err == nil {
			configSuc = true
			g.Log().Debugf(ctx, "服务配置地址请求成功，配置%s", addressList[i])
			break
		}

		g.Log().Warningf(ctx, "服务配置地址请求失败，请检测配置%s, err=%v", addressList[i], err)
	}
	if !configSuc {
		g.Log().Error(ctx, "服务配置地址请求全部失败")
		return nil, code.CodeError.New(ctx, code.CommonConsulCfgCurlAllError)
	}
	return adapter, nil
}
