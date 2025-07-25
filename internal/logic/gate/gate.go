package gate

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"math/rand"
	v1 "middle_srv/app/rpc/api/gate/v1"
	"middle_srv/internal/consts"
	"middle_srv/internal/service"
	"middle_srv/utility/code"
	"time"
)

type (
	sGate struct{}
)

func init() {
	service.RegisterGate(New())
}

func New() service.IGate {
	return &sGate{}
}
func (s *sGate) getConn(ctx context.Context, req *v1.CallRequest) (*grpc.ClientConn, error) {
	cacheKey := fmt.Sprintf(consts.SrvTmpHost, req.RegService)
	sHostVar, err := g.Redis().Get(ctx, cacheKey)
	if err != nil {
		g.Log().Warningf(ctx, "[gate] getConn get cache key:%s err:%v", cacheKey, err)
	}
	sHost := sHostVar.String()
	if sHost == "" {
		registry, err := service.SrvRegister().GetGsvcRegistry(ctx)
		if err != nil {
			return nil, err
		}

		instances, err := registry.Search(ctx, gsvc.SearchInput{
			Prefix:   "",
			Name:     req.RegService,
			Metadata: nil,
		})
		if err != nil { //搜索报错
			g.Log().Errorf(ctx, "Search RegService, RegService=%s, err=%v", req.RegService, err)
			return nil, code.CodeError.New(ctx, code.GateSearchRegServiceFail)
		}
		instanceLen := len(instances)
		if instanceLen == 0 { //相当于没注册过实例或全下线了
			g.Log().Errorf(ctx, "RegService len=0, RegService=%s", req.RegService)
			return nil, code.CodeError.New(ctx, code.GateSearchRegServiceLenError)
		}
		sHost = instances[rand.Intn(instanceLen)].GetEndpoints().String()
		err = g.Redis().SetEX(ctx, cacheKey, sHost, 30)
		if err != nil {
			g.Log().Warningf(ctx, "[gate] getConn SetEX key:%s err:%v", cacheKey, err)
		}
	}

	conn := grpcx.Client.MustNewGrpcClientConn(sHost)
	return conn, nil
}

func (s *sGate) verifyField(ctx context.Context, req *v1.CallRequest) error {
	fmt.Println(req)
	if req.RegService == "" {
		return code.CodeError.New(ctx, code.CommonRequiredError, "RegService")
	}
	if req.Service == "" {
		return code.CodeError.New(ctx, code.CommonRequiredError, "Service")
	}
	if req.Method == "" {
		return code.CodeError.New(ctx, code.CommonRequiredError, "Method")
	}
	return nil
}

func (s *sGate) Call(ctx context.Context, req *v1.CallRequest) (*v1.CallReply, error) {
	//cc, _ := g.Config().Get(ctx, "limited")
	//limitedCfg := &model.LimitedCfg{}
	//err := cc.Struct(limitedCfg)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(limitedCfg.Server["gate.service"].LimitNum)
	startTime := time.Now()
	err := s.verifyField(ctx, req)
	if err != nil {
		return nil, err
	}
	//限流处理
	//@todo 可以视情况而定，限制哪块，目前测试阶段先用RegService
	limiterKey := req.RegService
	limiter := service.Limiter().GetLimiter(ctx, req.RegService)
	if limiter != nil && !limiter.Allow() {
		//拒绝
		g.Log().Warningf(ctx, "触发限流, regService=%s", limiterKey)
		return nil, code.CodeError.New(ctx, code.GateLimiterError)
	}
	payload := req.GetPayload()
	conn, err := s.getConn(ctx, req)
	if err != nil {
		return nil, err
	}

	// 1. 反射客户端
	rc := grpcreflect.NewClientV1Alpha(ctx, reflectpb.NewServerReflectionClient(conn))
	defer rc.Reset()

	// 2. 查找方法描述
	svcDesc, err := rc.ResolveService(req.Service)
	if err != nil {
		g.Log().Errorf(ctx, "no found Service, regService=%s, service=%s, err=%v", req.RegService, req.Service, err)
		return nil, code.CodeError.New(ctx, code.GateSearchServiceFail)
	}

	mDesc := svcDesc.FindMethodByName(req.Method)
	if mDesc == nil {
		g.Log().Errorf(ctx, "no found Method, regService=%s, service=%s, method=%s",
			req.RegService, req.Service, req.Method)
		return nil, code.CodeError.New(ctx, code.GateSearchMethodFail)
	}
	if payload == "" {
		payload = "{}"
	}

	// 3. 构造请求消息
	reqMsg := dynamic.NewMessage(mDesc.GetInputType())
	if err := reqMsg.UnmarshalJSON([]byte(payload)); err != nil {
		g.Log().Errorf(ctx, "req payload UnmarshalJSON err,err=%v", err)
		return nil, code.CodeError.New(ctx, code.GatePayloadParamsError)
	}

	// 4. 动态调用
	stub := grpcdynamic.NewStub(conn)

	//超时控制
	dur := 5 * time.Second

	ctx, cancel := context.WithTimeout(ctx, dur)

	defer cancel()
	startTime2 := time.Now()
	defer func() {
		if time.Since(startTime) >= time.Millisecond*800 {
			var reqPayload string

			if len(req.Payload) > 1000 {
				reqPayload = fmt.Sprintf("%s...", hex.EncodeToString([]byte(req.Payload[:1000])))
			} else {
				reqPayload = hex.EncodeToString([]byte(req.Payload))
			}
			g.Log().Warningf(ctx, "RPC request slow log,regService=%s, service=%s, method=%s,cost:%v,rpc_cost:%v,reqPayload=%s",
				req.RegService, req.Service, req.Method, time.Since(startTime), time.Since(startTime2), reqPayload)
		}
	}()

	resp, err := stub.InvokeRpc(ctx, mDesc, reqMsg)

	if err != nil {
		codeErr := gerror.Code(err)

		if codes.Code(codeErr.Code()) == codes.DeadlineExceeded {
			g.Log().Errorf(ctx, "InvokeRpc 请求服务超时错误,regService=%s, service=%s, method=%s,err=%v",
				req.RegService, req.Service, req.Method, err)
			return nil, code.CodeError.New(ctx, code.GateRpcTimeout, req.RegService, req.Service, req.Method)
		}

		g.Log().Errorf(ctx, "InvokeRpc 错误,err=%v", err)
		return nil, err
	}
	return &v1.CallReply{
		Payload: resp.String(),
	}, nil
}
