package api

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"ruletest/util"
)

func (serverCtx *ServerContext) RuleCheck(r *ghttp.Request) {
	ruleStr := r.GetRequestString("rule")
	dataContext := context.NewDataContext()
	rb := builder.NewRuleBuilder(dataContext)
	if ruleStr != "" {
		if e := rb.BuildRuleFromString(ruleStr); e != nil {
			SendRsp(r, 400,  fmt.Sprintf("Error : build rule from string err: %+v", e), ruleStr)
		}

		for ruleName, ruleEntity := range rb.Kc.RuleEntities {
			for _, statement  := range ruleEntity.RuleContent.Statements.StatementList {
				if statement.MethodCall != nil {
					if serverCtx.RulesRunFuncsMap[statement.MethodCall.MethodName] == nil {
						SendRsp(r, 400, fmt.Sprintf("Error : rule [%v] func [%v] unsupported ", ruleName, statement.MethodCall.MethodName))
					}
				}
				if statement.FunctionCall != nil {
					if serverCtx.RulesRunFuncsMap[statement.FunctionCall.FunctionName] == nil {
						SendRsp(r, 400, fmt.Sprintf("Error : rule [%v] func [%v] unsupported ", ruleName, statement.FunctionCall.FunctionName))
					}
				}
			}
		}
	} else {
		SendRsp(r, 400, "Error : ruleStr is empty", ruleStr)
	}
	SendRsp(r, 200, "rule check successful", ruleStr)
}

func (serverCtx *ServerContext) RuleAdd(r *ghttp.Request) {
	ruleStr := r.GetRequestString("rule")
	dataContext := context.NewDataContext()
	rb := builder.NewRuleBuilder(dataContext)
	if ruleStr != "" {
		if e := rb.BuildRuleFromString(ruleStr); e != nil {
			SendRsp(r, 400, fmt.Sprintf("Error : build rule from string err: %+v", e), ruleStr)
		}

		// TODO： 加锁 && 排重 && raft
		for ruleName, ruleEntity := range rb.Kc.RuleEntities {
			rule := util.RuleConf{0, ruleName,ruleStr, "1", make(map[string]string), make(map[string]string)}
			rule.RuleId = len(serverCtx.RuleConfsMap) + 1
			for _, statement  := range ruleEntity.RuleContent.Statements.StatementList {
				if statement.MethodCall != nil {
					methodName := statement.MethodCall.MethodName
					rule.RuleRunFuncsMap[methodName] = methodName // TODO：需根据名称映射成具体方法/函数(go的难点，目前只能手工维护名字=》函数的映射了)
				}
				if statement.FunctionCall != nil {
					functionName := statement.FunctionCall.FunctionName
					rule.RuleRunFuncsMap[functionName] = functionName // 同上 TODO
				}
			}
			serverCtx.RuleConfsMap[rule.RuleId] = &rule
		}
	}
	SendRsp(r, 200, "rule add successful", ruleStr)
}

func (serverCtx *ServerContext) RuleGets(r *ghttp.Request) {

	SendRsp(r, 200, "rule get successful", serverCtx.RuleConfsMap)
}

/**
 * NOTE：目前支持不了动态设定可用函数，该函数无用(参看 ReadMe.md)
 *
 */
func (serverCtx *ServerContext) RuleFuncsAdd(r *ghttp.Request) {

	ruleFuncsJson := gjson.New(r.Get("ruleFuncs"))
	for k,v  := range ruleFuncsJson.Map(){
		serverCtx.RulesRunFuncsMap[k] = v
	}
	SendRsp(r, 200, "rule funcs add successful", serverCtx.RulesRunFuncsMap)
}
