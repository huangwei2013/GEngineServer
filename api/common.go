package api

import (
	"context"
	"errors"
	"sync"

	"github.com/coreos/etcd/raft/raftpb"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/raft"

	"net"
	"ruletest/util"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"google.golang.org/grpc"
)



type TenantRuleStatus int32
const (
	Rule_Reg     	TenantRuleStatus = 0 // no defined yet
	Rule_Normal     TenantRuleStatus = 1 // working fine
	Rule_Pending    TenantRuleStatus = 2 // to-be-offline/some other : no accept new request
)


type Tenant struct {
	TenantId 	 int
	RulesMap     map[int]TenantRuleStatus	// 对应 ServerContext.RuleConfsMap 中的 key

	//RuleRuntimeConfsMap map[int]RuleRuntimeConf	// key：serverId, value: config for EnginePool on this server
}

type ServerContext struct {
	ServerId int64
	StartAt int64
	HostIP string
	Port int
	Server *ghttp.Server
	GrpcServer *grpc.Server
	Ctx context.Context
	EctdClient *clientv3.Client

	Mu sync.Mutex

	// NOTE:从 etcd/raft 中将 raft+通讯部分剥离，暂时难度有点大（通讯部分耦合太多），改成直接使用完整 etcd
	Node *raft.Node
	Recvc chan raftpb.Message //从Stream消息通道中读取消息之后，会通过该通道将消息交给Raft接口，然后由它返回给底层etcd-raft模块进行处理。
	Propc chan raftpb.Message //从Stream消息通道中读取到MsgProp类型的消息之后，会通过该通道将MsgProp消息交给Raft接口，然后由它返回给底层etcd-raft模块进行处理

	// 需要动态维护的规则内容相关
	UpdateAt int64
	RuleVersion int
	Ready chan int

	// TODO：补充 last_update_version 等记录，优化 W/R 效率
	RuleConfsMapUpdatedFlag bool
	RuleConfsMap     map[int]*util.RuleConf
	RulesRunFuncsMapUpdatedFlag bool
	RulesRunFuncsMap map[string]interface{}
	RulesRunObjsMapUpdatedFlag bool
	RulesRunObjsMap  map[string]interface{}
	TenantsUpdatedFlag bool
	Tenants map[int]*Tenant

	Stop chan struct{}
	Stopping chan struct{}
	Done chan struct{}
}


func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := GetIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func GetIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}


func GetReqJson(r *ghttp.Request ) *gjson.Json {
	reqJson, err := gjson.DecodeToJson(r.GetRaw())
	if err != nil {
		SendRsp(r, 400, err.Error())
	}
	return reqJson
}

// 标准返回结果数据结构封装。
// 返回固定数据结构的JSON:
// code:  错误码(0:成功, 1:失败, >1:错误码);
// msg:  请求结果信息;
// data: 请求结果,根据不同接口返回结果的数据结构不同;
func SendRsp(r *ghttp.Request, code int, msg string, data ...interface{}){
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(g.Map{
		"code":  code,
		"message":  msg,
		"data": responseData,
	})
	r.Exit()
}