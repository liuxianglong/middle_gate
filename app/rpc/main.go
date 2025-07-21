package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"middle_srv/app/rpc/internal/cmd"
	_ "middle_srv/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"
	//_ "middle_srv/internal/boot"
	_ "middle_srv/internal/logic"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
