# pbqp

PBQP(Protobuf Quick Protocol), the simple protocol based on [Protobuf](github.com/golang/protobuf) with [SMUX](github.com/xtaci/smux)

## Great-Ghost Protocol(s) 幽灵计划协议族

### 链路层
+ 实现于：gnode

+ RSH 远程终端协议
+ CABLE 幽灵链路协议

### 网络层
+ 实现于：gnet

+ ARP 地址解析协议
+ DHCP 动态主机设置协议

### 传输层
+ 实现于：gnet

+ TCP
+ UDP
+ P2P

### 应用层
+ 实现于：gnet

+ PROXY 代理协议
