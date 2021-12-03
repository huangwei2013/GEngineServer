package api

import (
	"github.com/gogf/gf/net/ghttp"
)


func (serverCtx *ServerContext) ServerStatus(r *ghttp.Request) {

	SendRsp(r, 200, "Server Status get successful", serverCtx)
}
