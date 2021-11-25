package foo


import (
	"fmt"
	"github.com/bilibili/gengine/engine"
)

//业务规则
const Service_rules string = `
rule "1" "1"
begin
	//resp.At = room.GetAttention()
	println("rule 1...")
end 

rule "2" "2"
begin
	resp.Num = room.GetNum()
	println("rule 2...")
	fmt.println("rule 2...")
end
`

//业务接口
type EngineService struct {
	//gengine pool
	Pool *engine.GenginePool

	//other params
}

//request
type Request struct {
	Rid       int64
	RuleNames []string
	//other params
}

//resp
type Response struct {
	At  int64
	Num int64
	//other params
}

//特定的场景服务
type Room struct {
}

func (r *Room) GetAttention( /*params*/ ) int64 {
	// logic
	return 100
}

func (r *Room) GetNum( /*params*/ ) int64 {
	//logic
	return 111
}

//初始化业务服务
//apiOuter这里最好仅注入一些无状态函数，方便应用中的状态管理
func NewMyService(poolMinLen, poolMaxLen int64, em int, rulesStr string, apiOuter map[string]interface{}) *EngineService {
	pool, e := engine.NewGenginePool(poolMinLen, poolMaxLen, em, rulesStr, apiOuter)
	if e != nil {
		panic(fmt.Sprintf("初始化gengine失败，err:%+v", e))
	}

	myService := &EngineService{Pool: pool}
	return myService
}

//service
func (ms *EngineService) Service(req *Request) (*Response, error) {

	resp := &Response{}

	//基于需要注入接口或数据,data这里最好仅注入与本次请求相关的结构体或数据，便于状态管理
	data := make(map[string]interface{})
	data["req"] = req
	data["resp"] = resp

	//模块化业务逻辑,api
	room := &Room{}
	data["room"] = room

	//
	e,_ := ms.Pool.ExecuteSelectedRules(data, req.RuleNames)
	if e != nil {
		println(fmt.Sprintf("pool execute rules error: %+v", e))
		return nil, e
	}

	return resp, nil
}
