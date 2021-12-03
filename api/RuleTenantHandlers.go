package api

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"ruletest/service"
)


func (serverCtx *ServerContext) RuleTenantAdd(r *ghttp.Request) {
	tenantId := r.GetInt("TenantId")
	ruleId := r.GetInt("RuleId")
	if serverCtx.Tenants[tenantId] == nil {
		SendRsp(r, 400, "[Error] Invalided TenantId(%v)", tenantId)
	}
	tenant := serverCtx.Tenants[tenantId]
	if tenantRuleStatus, found := tenant.RulesMap[ruleId] ; !found {
		SendRsp(r, 400, "[Error] Invalided RuleId(%v)", ruleId)
		if tenantRuleStatus != Rule_Normal {
			SendRsp(r, 400, "[Error] Existed Rule(%v) but abnormal Status(%v)", ruleId, tenantRuleStatus)
		}
	} else{
		tenant.RulesMap[ruleId] = Rule_Normal // default status
		serverCtx.Tenants[tenantId] = tenant // update
	}
	SendRsp(r, 200, "rule add successful")
}


func (serverCtx *ServerContext) RuleTenantRun(r *ghttp.Request) {
	tenantId := r.GetInt("TenantId")
	ruleId := r.GetInt("RuleId")

	if serverCtx.Tenants[tenantId] == nil {
		SendRsp(r, 400, "[Error] Invalided TenantId(%v)", tenantId)
	}
	tenant := serverCtx.Tenants[tenantId]
	if tenantRuleStatus, found := tenant.RulesMap[ruleId] ; !found {
		SendRsp(r, 400, "[Error] Invalided RuleId(%v)", ruleId)
		if tenantRuleStatus != Rule_Normal {
			SendRsp(r, 400, "[Error] Invalided Rule(%v) by Status(%v)", ruleId, tenantRuleStatus)
		}
	}
	rule := serverCtx.RuleConfsMap[ruleId]
	if rule == nil{
		SendRsp(r, 400, "Error : Cannot find Rule by RuleId(%v)", ruleId)
	}

	ruleName := rule.RuleName
	ruleStr := rule.RuleContent
	apis := make(map[string]interface{})
	for k,_ := range rule.RuleRunFuncsMap{
		apis[k] = serverCtx.RulesRunFuncsMap[k]
	}
	engineService := service.NewEngineService(1, 2, 1, ruleStr, apis)

	//调用
	req := &service.Request{
		Rid:       int64(ruleId),
		RuleNames: []string{ruleName},
	}
	response, e := engineService.Run(req)
	if e != nil {
		SendRsp(r, 400, fmt.Sprintf("Error : Service Err : %+v", e))
	}
	SendRsp(r, 200, "", response[ruleName])
}