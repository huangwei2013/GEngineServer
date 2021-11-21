
# 说明

基于 [gengine](https://github.com/bilibili/gengine) 的多租户服务化改造
（还在初期）

## TODO
1、rule文本中备选函数集的通用处理流程
2、rule(Parser后的结构化内容)动态管理
3、集群化处理（raft、节点间任务分配)
4、rule 模板 到 rule 字符串的相互转换


### 方案
- [TODO] rule 中指定的对象/函数/方法集，如何管理
    - 运行上下文维护一个可执行的 对象/函数/方法全集，具体 rule 只是使用其中的指定 对象/函数/方法
- [TODO] golang 目前支持的两种动态加载代码方式（对应 rule 中新增指定调用的对象/函数/方法）
    - 通过 go plugin 方式，将要动态加载的代码，编译成 so，加载运行
    - 通过网络远程调用/gRPC方式，在被调用端实现动态（即被调用端可采用重启等方式进行加载）
- [TODO] raft
    - 准备借用 etcd-raft
- [TODO] 节点任务分配
    - 考虑 etcd 中类似代码
- [TODO] rule 模板 -- rule 字符串
    - 字符串 -> 模板：通过 DSL parser 
    - 模板 -> 字符串：比较简单，template 自身特性