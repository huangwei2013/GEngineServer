package foo

import (
	"go.etcd.io/etcd/raft"
	"ruletest/api"
)

func RunRaft(serverCtx *api.ServerContext) {

	storage := raft.NewMemoryStorage()
	c := &raft.Config{
		ID:              0x01,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	peers := []raft.Peer{{ID: 0x01}}
	n := raft.StartNode(c, peers)
	serverCtx.Node = &n
}