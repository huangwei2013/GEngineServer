package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"ruletest/api"
	"ruletest/foo"
	"ruletest/service/rule/actions"
	"ruletest/util"
	"time"
)


const RuleStr string = `
rule "测试规则名称1" "rule desc"
begin
	ReqData = 10 + 7 + 8
	//sleep(1000)   
end
rule "1"
begin
	ReqData = 10 + 7 + 8
	//sleep(1000)
	println("1")
end
rule "demo rule" "rule desc"
begin
	ReqData = 10 + 7 + 8
	//sleep(1000)
	return ReqData
end
`

var serverCtx = api.ServerContext{
	StartAt: time.Now().Unix(),
	Port: 8199,
	RuleConfsMap: make(map[int64]*util.RuleConf),
	ActionsMap: make(map[string]interface{}),
}

/**
	最主要的麻烦之处...
 */
func initServerCtx (){

	// NOTE: 实际场景，需要通过 Parser 解析 ruleStr 然后抽取对应字段，用于结构化内容生成

	// Parser 时需入库，这里从库中读取即可
	serverCtx.RuleConfsMap[1] = &util.RuleConf{1, "demo rule",RuleStr, "1"}

	// 名称信息来自解析入库的内容，但示例相关代码，需要另外处理
	serverCtx.ActionsMap["room"] = &actions.Room{}
}

func initEPHandlers(s *ghttp.Server){

	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Welcome to GEngineServer！")
	})

	s.BindHandler("/engine/{TenantId}/{ProjectId}/{RuleId}/*RuleVersion/run", serverCtx.RuleRun)
}


func main() {
	s := g.Server()

	s.SetConfigWithMap(g.Map{
		"address":          ":8081",
		"accessLogEnabled": true,
		"errorLogEnabled":  true,
		"pprofEnabled":     true,
		"logPath":          "./tmp/log/ServerLog",
		"sessionIdName":    "MySessionId",
		"sessionPath":      "./tmp/session",
		"sessionMaxAge":    24 * time.Hour,
		"dumpRouterMap":    false,
	})

	initServerCtx()
	initEPHandlers(s)

	go foo.RunRaft(&serverCtx)

	s.Run()
}
