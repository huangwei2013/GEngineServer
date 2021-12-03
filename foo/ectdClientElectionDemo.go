package foo

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"

	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

//const prefix = "/election/gengine"
//const prop = "local"
//
//func main() {
//	endpoints := []string{"localhost:2379"}
//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
//
//	Campaign(cli, prefix, prop)
//
//}

func InitEctdClient(endpoints []string) (*clientv3.Client){
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func Campaign(c *clientv3.Client, election string, prop string) {
	for {
		ttl := 30
		s, err := concurrency.NewSession(c, concurrency.WithTTL(ttl))
		if err != nil {
			fmt.Println(err)
			return
		}
		e := concurrency.NewElection(s, election)
		ctx := context.TODO()

		log.Println("开始竞选")

		err = e.Campaign(ctx, prop)
		if err != nil {
			log.Println("竞选 leader失败，继续")
			switch {
			case err == context.Canceled:
				return
			default:
				continue
			}
		}

		log.Println("获得leader")

		//TODO：写入leader声明

		if err := runAsLeader(); err != nil {
			log.Println("调用主方法失败，辞去leader，重新竞选")
			_ = e.Resign(ctx)

			//TODO：清理leader声明
			continue
		}
		return
	}
}

func WriteConf(cli *clientv3.Client, key string, value string)(error){
	requestTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := cli.Put(ctx, key, value)
	return  err
}

func ReadConf(cli *clientv3.Client, key string)(*clientv3.GetResponse, error){
	requestTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := cli.Get(ctx, key)
	return resp, err
}

func runAsLeader() error {
	for {
		fmt.Println("doCrontab")
		time.Sleep(time.Second * 4)
		time.Sleep(time.Second * 4)
		time.Sleep(time.Second * 4)
		return fmt.Errorf("sss")
	}
}