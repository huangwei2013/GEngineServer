package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"log"
	"net"
	"ruletest/api"
	"ruletest/foo"
	"ruletest/service/rule/actions"
	"ruletest/util"
	"time"

	"github.com/coreos/etcd/raft/raftpb"
	"google.golang.org/grpc"

)



var serverCtx api.ServerContext


/**
	最主要的麻烦之处...
 */
func initServerCtx (){

	serverCtx = api.ServerContext{
		ServerId: 		 1,
		StartAt:         time.Now().Unix(),

		// 从 etcd 中读取/同步完成后，才会更新该簇配置
		UpdateAt:         time.Now().Unix(),
		RuleVersion:	  1,
		Ready:			  make(chan int),

		//Port:            8199,
		RuleConfsMap:    make(map[int]*util.RuleConf),
		RulesRunObjsMap: make(map[string]interface{}),
		RulesRunFuncsMap: make(map[string]interface{}),
		Tenants : make(map[int]*api.Tenant),

		RuleConfsMapUpdatedFlag : false,
		RulesRunFuncsMapUpdatedFlag : false,
		RulesRunObjsMapUpdatedFlag : false,
		TenantsUpdatedFlag :false,

		Recvc : make(chan raftpb.Message, 10),

		Stop: make(chan struct{}),
		Stopping: make(chan struct{}),
		Done: make(chan struct{}),

	}

	// TODO:
	//  应改为从 etcd 获取一系列已经注册的.so，依次加载so以获取这里要调用的自定义对象
	//  但对于 go 内置库，暂时没有更好办法

	// 名称信息来自解析入库的内容，但示例相关代码，需要另外处理
	serverCtx.RulesRunObjsMap["room"] = &actions.Room{}

	serverCtx.RulesRunFuncsMap["println"] = fmt.Println
	serverCtx.RulesRunFuncsMap["sleep"] = time.Sleep
}

func initPeersServer(ctx context.Context, protocol string, address string) {
	// 监听本地端口
	lis, err := net.Listen(protocol, address)
	if err != nil {
		log.Printf("监听端口失败: %s", err)
		return
	}
	// 创建gRPC服务器
	serverCtx.GrpcServer = grpc.NewServer()

	// 注册服务
	//pb.RegisterHelloServer(server.GrpcServer, server)
	//pb.RegisterHeartBeatServer(server.GrpcServer, server)
	//pb.RegisterDataServer(server.GrpcServer, server)
	//pb.RegisterByeServer(server.GrpcServer, server)

	//reflection.Register(serverCtx.GrpcServer)

	// server Run...
	go func() {
		log.Println("服务启动中")
		err = serverCtx.GrpcServer.Serve(lis)
		if err != nil {
			log.Printf("开启服务失败: %s", err)
			return
		}
		log.Println("服务退出中")
	}()

	// ctx.Done() --> server.GracefulStop()
	go func(){
		for {
			select {
			case <-ctx.Done():
				serverCtx.GrpcServer.GracefulStop()
				log.Printf("Quiting by GracefulStop")
				return
			}
		}
	}()
}

func initEPHandlers(s *ghttp.Server){

	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Welcome to GEngineServer！")	})

	s.BindHandler("/engine/server/status", serverCtx.ServerStatus)

	s.BindHandler("/engine/rule/check", serverCtx.RuleCheck)
	s.BindHandler("/engine/rule/add", serverCtx.RuleAdd) // add rule
	s.BindHandler("/engine/rule/gets", serverCtx.RuleGets)
	//s.BindHandler("/engine/rule/func/add", serverCtx.RuleFuncsAdd)//暂无法支持动态函数加载

	s.BindHandler("/engine/{TenantId}/{ProjectId}/{RuleId}/*RuleVersion/add", serverCtx.RuleTenantAdd) // add an exist rule to this tenant
	s.BindHandler("/engine/{TenantId}/{ProjectId}/{RuleId}/*RuleVersion/run", serverCtx.RuleTenantRun) // run an exist rule which belongs to this tenant
}

func initCluster(){
	time.Sleep(1 * time.Second)

	const prefix = "/election/gengine"
	const prop = "cluster"

	endpoints := []string{"localhost:2379"}
	serverCtx.EctdClient = foo.InitEctdClient(endpoints)

	// TODO:续约routine

	// NOTE:从etcd中同步配置，同步完成执行 <-serverCtx.Ready
	syncRRuleConf()

	//<-serverCtx.Ready

	//go foo.Campaign(cli, prefix, prop)
}

// TODO: temp code ,to be remove
func syncWRuleConfRoutine(){
	t := time.NewTicker(5 * time.Second)
	for {
		<-t.C
		syncWRuleConf()
	}
}

func syncWRuleConf(){
	var bytes []byte
	var err error

	serverCtx.Mu.Lock()
	defer func(){
		serverCtx.RuleConfsMapUpdatedFlag = false
		serverCtx.RulesRunFuncsMapUpdatedFlag = false
		serverCtx.RulesRunObjsMapUpdatedFlag = false
		serverCtx.TenantsUpdatedFlag = false
		serverCtx.Mu.Unlock()
	}()

	if serverCtx.RuleConfsMapUpdatedFlag {
		bytes, _ = json.Marshal(serverCtx.RuleConfsMap)
		err = foo.WriteConf(serverCtx.EctdClient, "RuleConfsMap", string(bytes))
		if err != nil {
			fmt.Printf("Errer : %v \n", err)
			return
		}
	}

	if serverCtx.RulesRunFuncsMapUpdatedFlag {
		bytes, _ = json.Marshal(serverCtx.RulesRunFuncsMap)
		err = foo.WriteConf(serverCtx.EctdClient, "RulesRunFuncsMap", string(bytes))
		if err != nil {
			fmt.Printf("Errer : %v \n", err)
			return
		}
	}

	if serverCtx.RulesRunObjsMapUpdatedFlag {
		bytes, _ = json.Marshal(serverCtx.RulesRunObjsMap)
		err = foo.WriteConf(serverCtx.EctdClient, "RulesRunObjsMap", string(bytes))
		if err != nil {
			fmt.Printf("Errer : %v \n", err)
			return
		}
	}

	if serverCtx.TenantsUpdatedFlag {
		bytes, _ = json.Marshal(serverCtx.Tenants)
		err = foo.WriteConf(serverCtx.EctdClient, "Tenants", string(bytes))
		if err != nil {
			fmt.Printf("Errer : %v \n", err)
			return
		}
	}

	fmt.Printf("[%v] sync w done \n", time.Now())
}

// TODO: temp code ,to be remove
func syncRRuleConfRoutine(){
	t := time.NewTicker(10 * time.Second)
	for {
		<-t.C
		syncRRuleConf()
	}
}

func syncRRuleConf(){
	serverCtx.Mu.Lock()
	defer serverCtx.Mu.Unlock()

	if conf, ok := api.RuleReader1(serverCtx.EctdClient, "RuleConfsMap"); ok == nil {
		serverCtx.RuleConfsMap = conf
	}

	if conf, ok := api.RuleReader2(serverCtx.EctdClient, "RulesRunFuncsMap"); ok == nil{
		serverCtx.RulesRunFuncsMap = conf
	}
	if conf, ok := api.RuleReader2(serverCtx.EctdClient, "RulesRunObjsMap"); ok == nil{
		serverCtx.RulesRunObjsMap = conf
	}

	if conf, ok := api.RuleReader3(serverCtx.EctdClient, "Tenants"); ok == nil{
		serverCtx.Tenants = conf
	}

	fmt.Printf("[%v] sync r done \n", time.Now())
}

func initWatch(){

	wch := serverCtx.EctdClient.Watch(context.TODO(), "RuleConfsMap")

	go func(){
		for {
			select {
				case <-wch:
					serverCtx.Mu.Lock()
					if conf, ok := api.RuleReader1(serverCtx.EctdClient, "RuleConfsMap"); ok == nil {
						serverCtx.RuleConfsMap = conf
					}
					serverCtx.Mu.Unlock()
			}
		}
	}()
}


func main() {
	initServerCtx()

	serverCtx.Server = g.Server()
	serverCtx.Server.SetConfigWithMap(g.Map{
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

	initEPHandlers(serverCtx.Server)
	initCluster()
	//serverCtx.Ready <- 1

	initWatch()

	fmt.Println("GEngineServer init finished. Running... ")

	//go foo.RunRaft(&serverCtx)
	//go syncWRuleConfRoutine()
	//go syncRRuleConfRoutine()
	serverCtx.Server.Run()
}
