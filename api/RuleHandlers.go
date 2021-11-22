package api

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"ruletest/service"
	"time"
)


func (serverCtx *ServerContext) RuleRun(r *ghttp.Request) {

	r.Response.Writeln("TenantId：",r.Get("TenantId"))
	r.Response.Writeln("ProjectId：",r.Get("ProjectId"))
	r.Response.Writeln("RuleId：",r.Get("RuleId"))
	r.Response.Writeln("RuleVersion：",r.Get("RuleVersion"))

	ruleId := r.GetInt64("RuleId")
	rule := serverCtx.RuleConfsMap[ruleId]
	if rule == nil{
		SendRsp(r, 400, "Invalid RuleId")
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
	actionName := "room"
	response, e := engineService.Run(req, actionName, serverCtx.ActionsMap[actionName])
	if e != nil {
		SendRsp(r, 400, fmt.Sprintf("Service Err : %+v", e))
	}
	SendRsp(r, 200, "", response[ruleName])
}