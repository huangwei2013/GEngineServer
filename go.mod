module ruletest

go 1.16

require (
	github.com/bilibili/gengine v1.5.7
	github.com/coreos/etcd v3.3.27+incompatible
	github.com/gogf/gf v1.16.6
	go.etcd.io/etcd v3.3.27+incompatible
	go.etcd.io/etcd/client/v3 v3.5.1
	go.uber.org/zap v1.19.1 // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1 // indirect
	google.golang.org/grpc v1.40.0
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3
