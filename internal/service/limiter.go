// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"golang.org/x/time/rate"
)

type (
	ILimiter interface {
		// Init 初始化
		Init(ctx context.Context)
		GetLimiter(ctx context.Context, service string) *rate.Limiter
		// Lookup 监听,更新配置
		Lookup(ctx context.Context)
	}
)

var (
	localLimiter ILimiter
)

func Limiter() ILimiter {
	if localLimiter == nil {
		panic("implement not found for interface ILimiter, forgot register?")
	}
	return localLimiter
}

func RegisterLimiter(i ILimiter) {
	localLimiter = i
}
