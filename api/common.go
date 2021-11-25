package api

import (
	"context"
	"errors"
	"github.com/coreos/etcd/raft/raftpb"
	"go.etcd.io/etcd/raft"
	"net"
	"ruletest/util"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)


type Cluster struct {

}

type ServerContext struct {
	StartAt int64
	HostIP string
	Port int
	Server *ghttp.Server
	Ctx context.Context
	Cluster Cluster

	Node *raft.Node
	Recvc chan raftpb.Message //从Stream消息通道中读取消息之后，会通过该通道将消息交给Raft接口，然后由它返回给底层etcd-raft模块进行处理。
	Propc chan raftpb.Message //从Stream消息通道中读取到MsgProp类型的消息之后，会通过该通道将MsgProp消息交给Raft接口，然后由它返回给底层etcd-raft模块进行处理

	// 需要动态维护的规则内容相关
	RuleConfsMap     map[int]*util.RuleConf
	RulesRunFuncsMap map[string]interface{}
	RulesRunObjsMap  map[string]interface{}

	// stop signals the run goroutine should shutdown.
	Stop chan struct{}
	// stopping is closed by run goroutine on shutdown.
	Stopping chan struct{}
	// done is closed when all goroutines from start() complete.
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