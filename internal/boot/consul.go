package boot

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	consul "github.com/gogf/gf/contrib/config/consul/v2"
)

func init() {
	var (
		ctx = gctx.GetInitCtx()
		// Configure Consul client settings
		consulConfig = api.Config{
			Address:    "consul1:8500",                     // Consul server address
			Scheme:     "http",                             // Connection scheme (http/https)
			Datacenter: "dc1",                              // Datacenter name
			Transport:  cleanhttp.DefaultPooledTransport(), // HTTP transport with connection pooling
			//Token:      "3f8aeba2-f1f7-42d0-b912-fcb041d4546d", // ACL token for authentication
		}
		// Path in Consul's Key-Value store where configurations are stored
		configPath = "server/message2"
	)

	// Create Consul adapter with configuration
	// The adapter implements gcfg.Adapter interface for configuration management
	adapter, err := consul.New(ctx, consul.Config{
		ConsulConfig: consulConfig, // Consul client configuration
		Path:         configPath,   // Configuration path in KV store
		Watch:        true,         // Enable configuration watching for updates
	})
	if err != nil {
		// Log fatal error if client initialization fails
		g.Log().Fatalf(ctx, `New consul adapter error: %+v`, err)
	}

	// Set Consul adapter as the configuration adapter
	// This enables GoFrame to use Consul for configuration management
	g.Cfg().SetAdapter(adapter)
}
