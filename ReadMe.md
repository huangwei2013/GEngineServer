
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
    - 通过 go plugin 方式，将要动态加载的代码，编译成 so，加载运行
    - 通过网络远程调用/gRPC方式，在被调用端实现动态（即被调用端可采用重启等方式进行加载）
    - [goloader](https://github.com/dearplain/goloader)，利用mmap+go编译器来达到动态加载效果
- [TODO] raft
    - 准备借用 etcd-raft
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