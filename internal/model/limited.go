package model

import "golang.org/x/time/rate"

type LimitedCfg struct {
	Version string                   `yaml:"version"`
	Server  map[string]limitedServer `yaml:"server"`
}

type limitedServer struct {
	LimitNum int `yaml:"limit_num"` //限流数
	OutNum   int `yaml:"out_num"`   //允许超出数
}

type LimitedMapData struct {
	Limiter  *rate.Limiter
	LimitNum int
	OutNum   int
}

//func NewLimitedCfg() *LimitedCfg {}
