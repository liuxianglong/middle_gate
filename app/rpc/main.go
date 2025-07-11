package main

import (
	"middle_srv/app/rpc/internal/cmd"
	_ "middle_srv/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"
	//_ "middle_srv/internal/boot"
	_ "middle_srv/internal/logic"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
