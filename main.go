package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"ruletest/service"
	"ruletest/util"
	"time"
)

type RuleConf struct {
	RuleId      int64
	RuleName    string
	RuleContent string
	RuleVersion string
}

var RuleConfsMap = make(map[int64]*RuleConf)

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

func initEPHandlers(s *ghttp.Server, RuleConfsMap map[int64]*RuleConf){

	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("哈喽世界！")
	})

	s.BindHandler("/engine/{TenantId}/{ProjectId}/{RuleId}/*RuleVersion/run", func(r *ghttp.Request) {
		r.Response.Writeln("TenantId：",r.Get("TenantId"))
		r.Response.Writeln("ProjectId：",r.Get("ProjectId"))
		r.Response.Writeln("RuleId：",r.Get("RuleId"))
		r.Response.Writeln("RuleVersion：",r.Get("RuleVersion"))

		ruleId := r.GetInt64("RuleId")
		rule := RuleConfsMap[ruleId]
		if rule == nil{
			util.SendRsp(r, 400, "Invalid RuleId")
		}

		ruleName := rule.RuleName
		ruleStr := rule.RuleContent

		apis := make(map[string]interface{})
		apis["println"] = fmt.Println
		apis["sleep"] = time.Sleep
		engineService := service.NewEngineService(1, 2, 1, ruleStr, apis)

		//调用
		req := &service.Request{
			Rid:       ruleId,
			RuleNames: []string{ruleName},
		}
		response, e := engineService.Run(req)
		if e != nil {
			util.SendRsp(r, 400, fmt.Sprintf("Service Err : %+v", e))
		}
		util.SendRsp(r, 200, "", response[ruleName])
	})
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

	// NOTE: 实际场景，需要通过 Parser 解析 ruleStr 然后抽取对应字段，用于结构化内容生成
	RuleConfsMap[1] = &RuleConf{1, "demo rule",RuleStr, "1"}

	initEPHandlers(s, RuleConfsMap)
	s.Run()
}
