package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"ruletest/api"
	"ruletest/foo"
	"ruletest/service/rule/actions"
	"ruletest/util"
	"time"

	"github.com/coreos/etcd/raft/raftpb"

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
	StartAt:         time.Now().Unix(),
	//Port:            8199,
	RuleConfsMap:    make(map[int]*util.RuleConf),
	RulesRunObjsMap: make(map[string]interface{}),
	RulesRunFuncsMap: make(map[string]interface{}),

	Recvc : make(chan raftpb.Message, 10),
}

/**
	最主要的麻烦之处...
 */
func initServerCtx (){

	// 名称信息来自解析入库的内容，但示例相关代码，需要另外处理
	serverCtx.RulesRunObjsMap["room"] = &actions.Room{}

	serverCtx.RulesRunFuncsMap["println"] = fmt.Println
	serverCtx.RulesRunFuncsMap["sleep"] = time.Sleep
}

func initEPHandlers(s *ghttp.Server){

	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Welcome to GEngineServer！")
	})

	s.BindHandler("/engine/rule/check", serverCtx.RuleCheck)
	s.BindHandler("/engine/rule/add", serverCtx.RuleAdd)
	s.BindHandler("/engine/rule/gets", serverCtx.RuleGets)
	//s.BindHandler("/engine/rule/func/add", serverCtx.RuleFuncsAdd)//暂无法支持动态函数加载

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
