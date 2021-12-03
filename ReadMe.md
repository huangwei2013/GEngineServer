
# 说明

基于 [gengine](https://github.com/bilibili/gengine) 的多租户服务化改造
（还在初期）

## TODO
1、rule文本中备选函数集的通用处理流程
2、rule(Parser后的结构化内容)动态管理
3、集群化处理（raft：rule配置的分布式场景下一致性、节点间池管理)
4、rule 模板 到 rule 字符串的相互转换


### 方案
- [TODO] rule 中指定的对象/函数/方法集，如何管理
    - 运行上下文维护一个可执行的 对象/函数/方法全集，具体 rule 只是使用其中的指定 对象/函数/方法
- [TODO] golang 目前支持的以下几种动态加载代码方式（对应 rule 中新增指定调用的对象/函数/方法）
    - [通过 go plugin 方式]，将要动态加载的代码，编译成 so (linux)，加载运行(https://github.com/bilibili/gengine/wiki/%E9%AB%98%E7%BA%A7%E6%89%A9%E5%B1%95)
    - 通过网络远程调用/gRPC方式，在被调用端实现动态（即被调用端可采用重启等方式进行加载）
    - [goloader](https://github.com/dearplain/goloader)，利用mmap+go编译器来达到动态加载效果
- etcd：通过etcd共享配置
    - [TODO] 通过分布式锁获取写权限
- [TODO] 节点任务分配
    - 考虑 etcd 中类似代码
- [TODO] rule 模板 -- rule 字符串
    - 字符串 -> 模板：通过 DSL parser （gengine 中的 RuleBuilder）
    - 模板 -> 字符串：比较简单，template 自身特性
    

# FAQ
## Rule中引用的函数
### 命名
```golang
会按照对象的方式去找实现，而不是 module：
gengine.DataContext.ExecMethod()
    其中的处理，将 methodName 按 . 拆分成2部分，分别检查两部分在 dataContext 中的存在性
```

因此，[示例代码](https://github.com/bilibili/gengine/wiki/%E6%9C%80%E4%BD%B3%E5%AE%9E%E8%B7%B5)中，将 fmt.Println 赋给 key=println
```
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	msr := NewMyService(10, 20, 1, service_rules, apis)
```

两点不合理
1、未在加载 rule 时校验/处理 函数是否满足
2、仅考虑了对象的查找，对 module 的支持不完整

## etcd 相关的依赖

etcd的依赖关系有些坑，如：
- 引用的 raft，和引用的 raftpb，是来自不同项目 

```go
	"github.com/coreos/etcd/raft/raftpb"
	"go.etcd.io/etcd/raft"
```

### 配置分布式维护流程

```
node管理
    - 租约 和 续租约
    
ruleConf 管理
    - 分布式 lock 写
    - 每个 node 都建立 watch
    - node 对其更新时，根据数据版本和锁来判断是否有权限
  
注册
1.获取数据（包括版本）
2.若1.ok，建立watch通道
3.若2.ok，申请租约
4.若3.ok，本地启动维护租约routine（更新周期=租约 1/3 TTL）
  
(被动)读
1.watch 到变更
2.读取数据
3.合并到本地，合并失败则修改 node 状态为 latency

(主动)写
1.判定是否能获取锁：当前数据版本 vs 锁的版本
2.若1.ok（当前数据版本 = 锁的版本），尝试获取锁
3.若2.ok（获取到锁），获取数据，并判断当前数据版本是否等于锁的版本
4.若3.ok，执行更新

```