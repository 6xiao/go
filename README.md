go
==
generic code, RTFSC

Common : 公用代码，如Init，CheckPanic,  时间， 网络

keepalive.sh : 保活脚本

BitMapCache系列 : 一个内存位图缓存服务。比如需要表示QQ号是否在线等7种状态，可以创建一个4位缓存。如仅表示两种状态，1位的缓存足矣。在使用1位缓存时，每10亿用户仅用128M内存。


DctDst : 离散正弦变换及其逆变换，离散余弦变换及其逆变换


SimpleMsgChan : 可实现简单的消息队列功能

Amqp : 用以连接 RabbitMQ 等接受 amqp 协议的消息队列
