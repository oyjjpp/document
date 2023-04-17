# 分布式

## 集群分布式

## grpc

## 分布式系统优缺点，一致性是如何保证的

## etcd原理

## raft算法是那种一致性算法

## Gossip 协议

Gossip 协议是一种分布式系统中常用的协议，用于在不同节点之间传播信息。它的基本思想是通过节点之间的随机通信来达到信息传播的目的，类似于人们之间的闲话传递。

在 Gossip 协议中，每个节点都维护一个局部的信息列表，并定期随机选择另外一个节点进行通信，将自己的信息列表中的一部分发送给对方。对方收到信息后，会将其与自己的信息列表合并，然后继续通过随机通信向其他节点传播。通过这种方式，节点之间逐渐将信息同步，最终达到全局一致。

Gossip 协议的优点是具有很好的可扩展性和容错性，因为节点之间的通信是随机的，不依赖于固定的拓扑结构。同时，由于信息的传播是分散的、异步的，所以即使部分节点失效，整个系统也可以继续工作。

Gossip 协议在实际应用中被广泛使用，例如分布式数据库系统、P2P 网络、区块链等。


## 分布式id

- [一线大厂的分布式唯一 ID 生成方案是什么样的](https://zhuanlan.zhihu.com/p/140078865)
- [9种分布式ID生成方式](https://zhuanlan.zhihu.com/p/152179727)

### 分布式ID需要满足那些条件？

- 全局唯一：必须保证ID是全局性唯一的，基本要求
- 高性能：高可用低延时，ID生成响应要块，否则反倒会成为业务瓶颈
- 高可用：100%的可用性是骗人的，但是也要无限接近于100%的可用性
- 好接入：要秉着拿来即用的设计原则，在系统设计和实现上要尽可能的简单
- 趋势递增：最好趋势递增，这个要求就得看具体业务场景了，一般不严格要求

### 分布式ID都有哪些生成方式？

- UUID
- 数据库自增ID
- 数据库多主模式
- 号段模式
- Redis
- 雪花算法（SnowFlake）

### 基于UUID

```UUID
在Java/Golang等服务端语言的世界里，想要得到一个具有唯一性的ID，首先被想到可能就是语言自身的一些UUID库，
毕竟它有着全球唯一的特性。
那么UUID可以做分布式ID吗？答案是可以的，但是并不推荐；

// Java
public static void main(String[] args) {
    String uuid = UUID.randomUUID().toString().replaceAll("-","");
    System.out.println(uuid);
}

// Golang
import (
  uuid "github.com/satori/go.uuid"
)

func buildUUID(String[] args) {
    uuid.NewV4()
}

UUID的生成简单到只有一行代码，输出结果 c2b8c2b9e46c47e3b30dca3b0d447718，但UUID却并
不适用于实际的业务需求。像用作订单号UUID这样的字符串没有丝毫的意义，看不出和订单相关的有
用信息；而对于数据库来说用作业务主键ID，它不仅是太长还是字符串，存储性能差查询也很耗时，
所以不推荐用作分布式ID。
```

**UUID优点：**

- 生成足够简单，本地生成无网络消耗，具有唯一性

**UUID缺点：**

- 无序的字符串，不具备趋势自增特性
- 没有具体的业务含义
- UUID的字符串存储，查询效率慢
- 存储空间大

**UUID应用场景：**

- 类似生成token令牌的场景
- 不适用一些要求有趋势递增的ID场景

### 基于数据库自增ID

```mysql
基于数据库的auto_increment自增ID完全可以充当分布式ID;
具体实现：需要一个单独的MySQL实例用来生成ID，建表结构如下：

CREATE TABLE `increasing_id` (
 `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
 `value` CHAR(10) NOT NULL DEFAULT '',
 PRIMARY KEY (`id`) USING BTREE
)ENGINE=INNODB;

INSERT INTO `increasing_id`(VALUE) VALUES('');

当我们需要一个ID的时候，向表中插入一条记录返回主键ID，但这种方式有一个比较致命的缺点，访
问量激增时MySQL本身就是系统的瓶颈，用它来实现分布式服务风险比较大，不推荐！

```

**基于数据库自增ID优点：**

- 实现简单，ID单调自增，数值类型查询速度快
- 查询效率高
- 具有一定的业务可读

**基于数据库自增ID缺点：**

- DB单点存在宕机风险，无法扛住高并发场景

### 基于数据库集群模式

```mysql
这个方案就是解决mysql的单点问题，在auto_increment基本上面，设置step步长。
每台的初始值分别为1,2,3...N，步长为N（这个案例步长为4）
```

![基于数据库集群模式](./image/202211011457.jpg)

|||
|-|-|
|优点：|解决了单点问题|
|缺点：|一旦把步长定好后，就无法扩容；而且单个数据库的压力大，数据库自身性能无法满足高并发|
|应用场景：| 数据不需要扩容的场景|

### 雪花snowflake算法

```mysql
雪花算法（Snowflake）是twitter公司内部分布式项目采用的ID生成算法，开源后广受国内大厂的
好评，在该算法影响下各大公司相继开发出各具特色的分布式生成器。

1位标识符：始终是0

41位时间戳：41位时间截不是存储当前时间的时间截，而是存储时间截的差值（当前时间截 - 开始
时间截 )得到的值，这里的的开始时间截，一般是我们的id生成器开始使用的时间，由我们程序来指定的

10位机器标识码：可以部署在1024个节点，如果机器分机房（IDC）部署，这10位可以由 5位机房ID + 5位机器ID 组成

12位序列：毫秒内的计数，12位的计数顺序号支持每个节点每毫秒(同一机器，同一时间截)产生4096个ID序号
```

![雪花snowflake算法](./image/202211011802.jpeg)

**雪花算法优点：**

- 此方案每秒能够产生409.6万个ID，性能快
- 时间戳在高位，自增序列在低位，整个ID是趋势递增的，按照时间有序递增
- 灵活度高，可以根据业务需求，调整bit位的划分，满足不同的需求

**雪花算法缺点：**

- 依赖机器的时钟，如果服务器时钟回拨，会导致重复ID生成

```mysql
在分布式场景中，服务器时钟回拨会经常遇到，一般存在10ms之间的回拨；小伙伴们就说这点10ms，
很短可以不考虑吧。但此算法就是建立在毫秒级别的生成方案，一旦回拨，就很有可能存在重复ID。
```

### Redis生成方案

```mysql
利用redis的incr原子性操作自增，一般算法为：

年份 + 当天距当年第多少天 + 天数 + 小时 + redis自增
```

|||
|-|-|
|优点： | 有序递增，可读性强 |
|缺点： | 占用带宽，每次要向redis进行请求 |

## 几种分布式锁的实现方式

[几种分布式锁的实现方式](https://juejin.cn/post/6844903863363829767)
[七种方案！探讨Redis分布式锁的正确使用姿势](https://z.itpub.net/article/detail/0A3DCC6FF8BD96C478FF1D7644DBFA57)

### 分布式锁有以下特点

- 可重入
- 同一时间点,只有一个线程持有锁
- 容错性, 当锁节点宕机时, 能及时释放锁
- 高性能
- 无单点问题

|方案|优点|缺点|
|-|-|-|
|数据库|操作简单，容易理解|性能开销大|
|redis|非阻塞，性能好|运维成本高，操作不好容易引起死锁|
|memecached|非阻塞，性能好|运维成本高，操作不好容易引起死锁|
|zookeeper|集群，无单点问题，可冲入，可避免锁无法释放|有性能瓶颈，性能不如redis|

### 基于数据库的分布式

[image](./image/202212101444001.jpg)

基于数据库的分布式锁, 常用的一种方式是使用表的唯一约束特性;当往数据库中成功插入一条数据时, 代表只获取到锁；将这条数据从数据库中删除，则释放送。

**因此需要创建一张锁表：**

```mysql
CREATE TABLE `methodLock` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `method_name` varchar(64) NOT NULL DEFAULT '' COMMENT '锁定的方法名',
  `cust_id` varchar(1024) NOT NULL DEFAULT '客户端唯一编码',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '保存数据时间，自动生成',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_method_name` (`method_name `) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='锁定中的方法';
```

**mysql 添加锁：**

```mysql
insert into methodLock(method_name, cust_id) values ('method_name', 'cust_id');
```

这里cust_id 可以是机器的mac地址+线程编号, 确保一个线程只有唯一的一个编号。通过这个编号， 可以有效的判断是否为锁的创建者，从而进行锁的释放以及重入锁判断

**mysql 释放锁：**

```mysql
delete from methodLock where method_name ='method_name' and cust_id = 'cust_id';
```

**重入锁判断：**

```mysql
select 1 from methodLock where method_name ='method_name' and cust_id = 'cust_id';
```

**加锁以及释放锁的代码示例：**

```golang

func GetGoProcessId() uint64 {
 b := make([]byte, 64)
 b = b[:runtime.Stack(b, false)]
 b = bytes.TrimPrefix(b, []byte("goroutine "))
 b = b[:bytes.IndexByte(b, ' ')]
 n, err := strconv.ParseUint(string(b), 10, 64)
 if err != nil {
  panic(err)
 }
 return n
}

func lock(methodName string) bool {
 success := false
 custId := GetGoProcessId()

 var err error
 success, err = insertLock(methodName, fmt.Sprintf("%d", custId))

 if err != nil {
  return false
 }
 return success
}

func unLock(methodName string) bool {
 success := false
 custId := GetGoProcessId()

 var err error
 success, err = deleteLock(methodName, fmt.Sprintf("%d", custId))

 if err != nil {
  return false
 }
 return success
}

// 是否可以重入锁
func checkReentrantLock(methodName string) bool {
 return true
}

func insertLock(method, custId string) (bool, error) {
 return true, nil
}

func deleteLock(method, custId string) (bool, error) {
 return true, nil
}

// 测试案例
func Test() {
 methodName := "test"
 if !checkReentrantLock(methodName) {
  for !lock(methodName) {
   time.Sleep(time.Second)
  }
 }

 // TODO 业务

 unLock(methodName)
}
```

**以上代码还存在一些问题：**

没有失效时间；解决方案:设置一个定时处理, 定期清理过期锁  
单点问题；解决方案: 弄几个备份数据库，数据库之前双向同步，一旦挂掉快速切换到备库上

### 基于redis的分布式

[image](./image/202210311901.jpg)  
[基于redis的分布式锁](../redis/distributed_lock.md)


