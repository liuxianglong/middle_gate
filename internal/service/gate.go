// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "middle_srv/app/rpc/api/gate/v1"
)

type (
	IGate interface {
		Call(ctx context.Context, req *v1.CallRequest) (*v1.CallReply, error)
	}
)

var (
	localGate IGate
)

func Gate() IGate {
	if localGate == nil {
		panic("implement not found for interface IGate, forgot register?")
	}
	return localGate
}

func RegisterGate(i IGate) {
	localGate = i
}
