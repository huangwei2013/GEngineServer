package foo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/coreos/etcd/raft/raftpb"
	"go.etcd.io/etcd/raft"
	"net"
	"ruletest/api"
	"time"
)

func RunRaft(serverCtx *api.ServerContext) {
	serverId := uint64(serverCtx.ServerId)
	storage := raft.NewMemoryStorage()
	c := &raft.Config{
		ID:              serverId,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	peers := []raft.Peer{{ID: serverId}}
	n := raft.StartNode(c, peers)
	serverCtx.Node = &n
	var term uint64 = 1

	go func(){
		t := time.NewTicker(120 * time.Second)

		for {
			<- t.C
			term += term
			fmt.Println("[term] ", term)
		}
	}()

	// fake recving
	go func(term uint64) {
		t := time.NewTicker(5 * time.Second)
		msg := raftpb.Message{From: 1, To: 1, Type: raftpb.MsgHup}

		for {
			<- t.C
			fmt.Printf("%s %v \n", msg.Type, time.Now())
			serverCtx.Recvc <- msg
			msg = raftpb.Message{From: 1, To: 1, Term: term, Type: raftpb.MsgTimeoutNow, LogTerm: 11, Index: 11}
			//fmt.Println(msg)
		}
	}(term)

	go func() {
		t := time.NewTicker(2 * time.Second)
		msg := raftpb.Message{From: 1, To: 1, Term: term, Type: raftpb.MsgBeat, Entries: []raftpb.Entry{{Data: []byte("somedata")}}}

		for {
			<- t.C
			fmt.Printf("MsgBeat %v \n", time.Now())
			serverCtx.Recvc <- msg
			//fmt.Println(msg)
		}
	}()

	go func() {
		t := time.NewTicker(60 * time.Second)
		msg := raftpb.Message{From: 1, To: 1, Term: term, Type: raftpb.MsgHup}

		for {
			<- t.C
			fmt.Printf("%s %v \n", msg.Type, time.Now())
			msg.Term += msg.Term
			serverCtx.Recvc <- msg
		}
	}()


	time.Sleep(3 * time.Second)

	for {
		select{
			case msg := <-serverCtx.Recvc:
				fmt.Println("Raft Msg Recv")
				recvMsg(n, msg)
			case <-serverCtx.Stop:
				n.Stop()
				fmt.Println("Server Stop")
				return
		}
	}
}

// 接收外部msg，处理
func recvMsg(n raft.Node, msg raftpb.Message){
	fmt.Println("【recvMsg】", n.Status())

	//n.Campaign(context.TODO())
	//fmt.Println("【Campaign】",n.Status())

	err := n.Step(context.TODO(), msg)
	if err != nil{
		fmt.Println("【Step-Error】",n.Status())
	}
	fmt.Println("【Step】",n.Status())
}

// 根据自己状态，向外界发送msg
func sendMsg(msg raftpb.Message){

}


func send (data []byte, scheme, addr string) {
	var out net.Conn
	var err error

	out, err = net.Dial(scheme, addr)
	if err != nil {
		fmt.Println(err)
	}
	if _, err = out.Write(data); err != nil {
		fmt.Println(err)
	}
	if err = out.Close(); err != nil {
		fmt.Println(err)
	}
}


func listen(scheme, addr string) (ln net.Listener) {
	ln, err := net.Listen(scheme, addr)
	if err != nil {
		fmt.Println(err)
	}
	return ln
}

/**
	ln := listen(scheme, dstAddr, tlsInfo)
	defer ln.Close()
*/

func receive(ln net.Listener) (data []byte) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	for {
		in, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		var n int64
		n, err = buf.ReadFrom(in)
		if err != nil {
			fmt.Println(err)
		}
		if n > 0 {
			break
		}
	}
	return buf.Bytes()
}