# Wireshark

开源网络包分析器

## 常见过滤规则

### 过滤MAC地址

eth.addr == e0:3f:49:28:87:88 //过滤目标或源地址是e0:3f:49:28:87:88的数据包  
eth.src == e0:3f:49:28:87:88 //过滤源地址是e0:3f:49:28:87:88的数据包  
eth.dst == e0:3f:49:28:87:88 //过滤目标地址是e0:3f:49:28:87:88的数据包  

### 过滤VLAN

vlan.id == 2594

### IP过滤

#### 源地址过滤

ip.src == 192.168.0.208  
ip.src eq 192.168.0.208  

#### 目标地址过滤

ip.dst == 192.168.0.208  
ip.dst eq 192.168.0.208  

#### ip地址过滤，不论源还是目标

ip.addr == 192.168.0.208  
ip.addr eq 192.168.0.208  

### 端口过滤

tcp.port == 80  
udp.port eq 80  

tcp.dstport == 80 // 只显tcp协议的目标端口80  
tcp.srcport == 80 // 只显tcp协议的来源端口80  

#### 过滤端口范围

tcp.port >= 1 and tcp.port <= 80

### 常用协议过滤

```network
tcp  
udp  
arp  
icmp  
http  
smtp  
ftp  
dns  
msnms  
ip  
ssl  
```

排除ssl包，如!ssl 或者 not ssl

### http模式过滤

http.request.method == "GET"  
http.request.method == "POST"  
http.request.uri == "/img/logo-edu.gif"  
http contains "GET"  
http contains "HTTP/1."

#### GET包

http.request.method == "GET" && http contains "Host: "  
http.request.method == "GET" && http contains "User-Agent: "  

#### POST包

http.request.method == "POST" && http contains "Host: "  
http.request.method == POST && http contains "User-Agent: "

#### 响应包

http contains "HTTP/1.1 200 OK" && http contains "Content-Type: "  
http contains "HTTP/1.0 200 OK" && http contains "Content-Type: "

#### 域名过滤

http.host == "www.baidu.com"

### 运算符

less than：lt  
less and equal： le  
equal：eq  
great then：gt  
great and equal：ge  
not equal：ne  

### 过滤内容

tcp[20] 表示从20开始，取1个字符  
tcp[20:]表示从20开始，取1个字符以上  
tcp[20:8]表示从20开始，取8个字符  
tcp[offset,n]  

udp[8:3]==81:60:03 // 偏移8个bytes,再取3个数，是否与==后面的数据相等？  
udp[8:1]==32 如果我猜的没有错的话，应该是udp[offset:截取个数]=nValue  
eth.addr[0:3]==00:06:5B  

#### 案例一

判断upd下面那块数据包前三个是否等于0x20 0x21 0x22  
我们都知道udp固定长度为8  
udp[8:3]==20:21:22  

#### 案例二

判断tcp那块数据包前三个是否等于0x20 0x21 0x22  
tcp一般情况下，长度为20,但也有不是20的时候  
tcp[8:3]==20:21:22  

#### 案例三

matches(匹配)和contains(包含某字符串)语法  
ip.src==192.168.1.107 and udp[8:5] matches "\x02\x12\x21\x00\x22"  
ip.src==192.168.1.107 and udp contains 02:12:21:00:22  
ip.src==192.168.1.107 and tcp contains "GET"  
udp contains 7c:7c:7d:7d 匹配payload中含有0x7c7c7d7d的UDP数据包，不一定是从第一字节匹配。  
