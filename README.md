go
==
generic code, RTFSC


Amqp : 用以连接 RabbitMQ 等接受 amqp 协议的 message queue


BitMapCache 系列 : 一个内存位图缓存服务。比如需要表示QQ号是否在线等7种状态，可以创建一个4位缓存。如仅表示两种状态，1位的缓存足矣。在使用1位缓存时，每10亿用户仅用128M内存


BloomFIlter : 布隆过滤器的Go语言实现，默认提供 8 个算子


Common : 公用代码，如Init，CheckPanic, 退出信号， 时间转换，日志， 网络，Zip，Hash，加密等


DctDst : 离散正弦变换及其逆变换，离散余弦变换及其逆变换


GoPool : 协程池，用以使用固定数量的 goroutine 顺序处理大量事件的场景


SimpleMsgChan : 简单的 Pub/Sub message queue


TrieTree : TrieTree (字典树)的 Go 语言实现，使用 map 和递归


WatchDog : 看门狗实现，用以监控容易失控的循环或超时


keepalive.sh : 保活脚本模板，使用 crontab / systemd 执行

