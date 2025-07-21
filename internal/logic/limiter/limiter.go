package limiter

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/time/rate"
	"middle_srv/internal/model"
	"middle_srv/internal/service"
	"sync"
	"time"
)

// sLimiter 限流器
type (
	sLimiter struct {
		limiterMap sync.Map
		version    string
	}
)

func init() {
	service.RegisterLimiter(New())
}

func New() service.ILimiter {
	return &sLimiter{}
}

func (s *sLimiter) getCfg(ctx context.Context) (*model.LimitedCfg, error) {
	cc, _ := g.Config().Get(ctx, "limited")
	limitedCfg := &model.LimitedCfg{}
	err := cc.Struct(limitedCfg)
	if err != nil {
		return nil, err
	}
	return limitedCfg, nil
}

// Init 初始化
func (s *sLimiter) Init(ctx context.Context) {

	limitedCfg, err := s.getCfg(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "sLimiter.Init 失败， err=%v", err)
		return
	}
	s.version = limitedCfg.Version

	//通过consul拿到对应配置
	s.limiterMap = sync.Map{}
	for serverName, v := range limitedCfg.Server {
		limiter := rate.NewLimiter(rate.Limit(v.LimitNum), v.OutNum)
		limitedMapData := &model.LimitedMapData{
			Limiter:  limiter,
			LimitNum: v.LimitNum,
			OutNum:   v.OutNum,
		}
		s.limiterMap.Store(serverName, limitedMapData)
	}

	return
}

func (s *sLimiter) GetLimiter(ctx context.Context, service string) *rate.Limiter {
	v, ok := s.limiterMap.Load(service)
	if ok {
		return v.(*model.LimitedMapData).Limiter
		//v.(*rate.Limiter)
	}
	return nil
}

// Lookup 监听,更新配置
func (s *sLimiter) Lookup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			limitedCfg, err := s.getCfg(ctx)
			if err != nil {
				g.Log().Errorf(ctx, "sLimiter.Lookup 失败， err=%v", err)
				return
			}
			if limitedCfg.Version != s.version {
				for serverName, v := range limitedCfg.Server {
					//和现有的做对比，如果没变化则不动
					curCfg, ok := s.limiterMap.Load(serverName)
					if ok && curCfg.(*model.LimitedMapData).LimitNum == v.LimitNum &&
						curCfg.(*model.LimitedMapData).OutNum == v.OutNum {
						continue
					}
					limiter := rate.NewLimiter(rate.Limit(v.LimitNum), v.OutNum)
					limitedMapData := &model.LimitedMapData{
						Limiter:  limiter,
						LimitNum: v.LimitNum,
						OutNum:   v.OutNum,
					}
					s.limiterMap.Store(serverName, limitedMapData)
					
					g.Log().Debugf(ctx, "sLimiter.Lookup 更改%s成功，limit=%d，out=%d",
						serverName, v.LimitNum, v.OutNum)
				}
			}
		}
	}
}
