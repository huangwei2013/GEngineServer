package foo

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"

	"go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"

	"time"
)

const (
	addressServer     = "localhost:12379"
	addressClient     = "localhost:2379"
)
//
//func DemoServer(){
//	// accept first connection so client is created with dial timeout
//	_, err := net.Listen("tcp", addressServer)
//	if err != nil {
//		fmt.Println(err)
//	}
//	//defer ln.Close()
//
//	ep := "unix://"+addressServer
//	cfg := clientv3.Config{
//		Endpoints:   []string{ep},
//		DialTimeout: 30 * time.Second}
//	c, err := clientv3.New(cfg)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	// connect to ipv4 black hole so dial blocks
//	c.SetEndpoints(addressServer)
//
//	// issue Get to force redial attempts
//	getc := make(chan struct{})
//	go func() {
//		defer close(getc)
//		// Get may hang forever on grpc's Stream.Header() if its
//		// context is never canceled.
//		c.Get(c.Ctx(), "abc")
//	}()
//
//	// wait a little bit so client close is after dial starts
//	time.Sleep(100 * time.Millisecond)
//
//	donec := make(chan struct{})
//	go func() {
//		defer close(donec)
//		c.Close()
//	}()
//
//	select {
//	case <-time.After(5 * time.Second):
//		fmt.Println("failed to close")
//	case <-donec:
//	}
//	select {
//	case <-time.After(5 * time.Second):
//		fmt.Println("get failed to exit")
//	case <-getc:
//	}
//}

//func DemoClient(){
//	cli, err := clientv3.New(clientv3.Config{
//		Endpoints: []string{"127.0.0.1:2379"},
//	})
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer cli.Close()
//	kvc := clientv3.NewKV(cli)
//	fmt.Println("client")
//
//	// perform a delete only if key already exists
//	key := "purpleidea"
//	value := "hello world"
//
//	_, err = kvc.Txn(context.Background()).
//		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
//		Then(clientv3.OpPut(key, value)).
//		Commit()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = kvc.Txn(context.Background()).
//		If(clientv3.Compare(clientv3.Version(key), ">", 0)).
//		Then(clientv3.OpDelete("purpleidea")).
//		Commit()
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println("end")
//
//}


func RunLockDemo() {
	c := make(chan os.Signal)
	signal.Notify(c)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer cli.Close()

	fmt.Println("try to lock")
	lockKey := "/lock"

	go func () {
		session, err := concurrency.NewSession(cli)
		if err != nil {
			fmt.Println(err)
		}
		m := concurrency.NewMutex(session, lockKey)
		if err := m.Lock(context.TODO()); err != nil {
			fmt.Println("go1 get mutex failed " + err.Error())
		}
		fmt.Printf("go1 get mutex sucess\n")
		fmt.Println(m)
		time.Sleep(time.Duration(10) * time.Second)
		m.Unlock(context.TODO())
		fmt.Printf("go1 release lock\n")
	}()

	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		session, err := concurrency.NewSession(cli)
		if err != nil {
			fmt.Println(err)
		}
		m := concurrency.NewMutex(session, lockKey)
		if err := m.Lock(context.TODO()); err != nil {
			fmt.Println("go2 get mutex failed " + err.Error())
		}
		fmt.Printf("go2 get mutex sucess\n")
		fmt.Println(m)
		time.Sleep(time.Duration(2) * time.Second)
		m.Unlock(context.TODO())
		fmt.Printf("go2 release lock\n")
	}()

	fmt.Println("running & waiting")

	<-c
}

func RunEctdClient() {

	RunLockDemo()
}

//func main(){
//
//	/**
//		InitEctdClient()
//		go func(){
//			time.Sleep(3 * time.Second)
//
//			for {
//				select{
//					case msg := <-serverCtx.Recvc:
//					case <-serverCtx.Stop:
//						fmt.Println("Server Stop")
//						return
//					default:
//
//				}
//			}
//		}
//	 */
//}