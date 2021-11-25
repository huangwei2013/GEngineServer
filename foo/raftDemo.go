package foo

import (
	"fmt"
	"github.com/coreos/etcd/raft/raftpb"
	"go.etcd.io/etcd/raft"
	"ruletest/api"
	"time"
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

	// fake recving
	go func() {
		t := time.NewTicker(2 * time.Second)
		msg := raftpb.Message{From: 1, To: 1, Type: raftpb.MsgProp, Entries: []raftpb.Entry{{Data: []byte("somedata")}}}

		for {
			<- t.C
			fmt.Printf("%v \n", time.Now())
			serverCtx.Recvc <- msg
			//fmt.Println(msg)
		}
	}()

	time.Sleep(3 * time.Second)

	for {
		select{
			case msg := <-serverCtx.Recvc:
				fmt.Println("Raft Msg Recv")
				recvMsg(n, msg)
			case <-serverCtx.Stop:
				fmt.Println("Server Stop")
				return
		}
	}
}

// 接收外部msg，处理
func recvMsg(n raft.Node, msg raftpb.Message){
	fmt.Println(n.Status())

	//n.Step(context.TODO(), msg)
	//fmt.Println(n.Status())
	//
	//n.Campaign(context.TODO())
	//fmt.Println(n.Status())
}

// 根据自己状态，向外界发送msg
func sendMsg(msg raftpb.Message){

}

