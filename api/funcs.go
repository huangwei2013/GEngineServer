package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ruletest/util"

	"go.etcd.io/etcd/client/v3"
)

func ReadConf(cli *clientv3.Client, key string)(*clientv3.GetResponse, error){
	requestTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := cli.Get(ctx, key)
	return resp, err
}

func RuleReader1(cli *clientv3.Client, key string)(map[int]*util.RuleConf, error){
	resp , err := ReadConf(cli, key)
	if err != nil{
		fmt.Printf("Error : %v \n", err)
	} else{
		if len(resp.Kvs) == 0 { return map[int]*util.RuleConf{}, nil }

		if string(resp.Kvs[0].Key[:]) == key {
			//val := string(resp.Kvs[0].Value[:])
			//var tempMap map[string]interface{}
			var tempMap map[int]*util.RuleConf
			err = json.Unmarshal(resp.Kvs[0].Value[:], &tempMap)
			if err != nil {
				fmt.Printf("Error : %v \n", err)
				return nil, err
			}
			return tempMap,nil
		}
	}
	return nil, err
}

func RuleReader2(cli *clientv3.Client, key string)(map[string]interface{}, error){
	resp , err := ReadConf(cli, key)
	if err != nil{
		fmt.Printf("Error : %v \n", err)
	} else{
		if string(resp.Kvs[0].Key[:]) == key {
			if  resp.Kvs[0].Value == nil { return map[string]interface{}{}, nil }

			var tempMap map[string]interface{}
			err = json.Unmarshal(resp.Kvs[0].Value[:], &tempMap)
			if err != nil {
				fmt.Printf("Error : %v \n", err)
				return nil, err
			}
			return tempMap, nil
		}
	}
	return nil, err
}

func RuleReader3(cli *clientv3.Client, key string)(map[int]*Tenant, error){
	resp , err := ReadConf(cli, key)
	if err != nil{
		fmt.Printf("Error : %v \n", err)
	} else{
		if len(resp.Kvs) == 0 { return map[int]*Tenant{}, nil }

		if string(resp.Kvs[0].Key[:]) == key {
			var tempMap map[int]*Tenant
			err = json.Unmarshal(resp.Kvs[0].Value[:], &tempMap)
			if err != nil {
				fmt.Printf("Error : %v \n", err)
				return nil, err
			}
			return tempMap, nil
		}
	}
	return nil, err
}
