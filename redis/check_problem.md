# 如何排查redis性能问题

![image](./image/redis_slowlog.png)

## redis真的变慢了嘛？

有完整的链路跟踪日志，也就是在服务访问外部依赖的出入口，记录下每次请求外部依赖的响应延时；

### 如果确认是redis链路问题？那我们在分析可能的原因？

1、业务服务器到redis服务器之间的网络存在问题  
2、redis本身存在问题，需要进一步排查是什么原因导致redis变慢  

### 对redis进行基准性测试

了解redis在生产环境服务器上的基准性能  
基准性能就是指redis在一台负载正常的机器上，其最大响应延迟和平均响应延时分别是多少？  

```redis
[root@eac5f6fb4074 redis]# ./redis-cli -h 127.0.0.1 -p 6379 --intrinsic-latency 60
Max latency so far: 1 microseconds.
Max latency so far: 15 microseconds.
Max latency so far: 45 microseconds.
Max latency so far: 46 microseconds.
Max latency so far: 1053 microseconds.

Max latency so far: 4280 microseconds.
Max latency so far: 4449 microseconds.
Max latency so far: 5397 microseconds.

1644822599 total runs (avg latency: 0.0365 microseconds / 364.78 nanoseconds per run).
Worst run took 147952x longer than the average latency.
```

从输出结果看到，60s内的最大响应延迟为5397微妙，平均0.0365微妙  

还可以使用如下命令，查看一段时间内Redis的最小、最大、平均访问延迟  

```redis
[root@eac5f6fb4074 redis]# ./redis-cli -h 127.0.0.1 -p 6379 --latency-history -i 1
min: 0, max: 4, avg: 0.23 (96 samples) -- 1.00 seconds range
min: 0, max: 1, avg: 0.27 (96 samples) -- 1.01 seconds range
min: 0, max: 1, avg: 0.22 (96 samples) -- 1.01 seconds range
min: 0, max: 1, avg: 0.27 (95 samples) -- 1.00 seconds range
min: 0, max: 1, avg: 0.28 (96 samples) -- 1.01 seconds range
...
```

了解了基准性能测试方法，就可以按照以下几步，来判断Redis是否真的变慢了：

在相同配置的服务器上，测试一个正常Redis实例的基准性能  
找到认为可能变慢的Redis实例，测试这个实例的基准性能  
如果你观察到，这个实例的运行延迟是正常Redis基准性能的2倍以上，即可认为这个 Redis 实例确实变慢了。

## 使用复杂度过高的命令

查看redis的慢日志（slowlog）

```redis
# 查看当前慢日志时间的阈值
127.0.0.1:6379> config get slowlog-log-slower-than
1) "slowlog-log-slower-than"
2) "10000"
# 查看当前保存的日志数目
127.0.0.1:6379> config get slowlog-max-len
1) "slowlog-max-len"
2) "128"
```

然后通过slowlog查看记录的慢查询日志

```redis
127.0.0.1:6379> slowlog get 5
(empty list or set)
```

### Redis 命令有以下特点，那么有可能会导致操作延迟变大

1、经常使用 O(N) 以上复杂度的命令，例如 SORT、SUNION、ZUNIONSTORE 聚合类命令  
2、使用O(N)复杂度的命令，并且N的值非常大  

第一种情况导致变慢的原因在于，Redis在操作内存数据时，时间复杂度过高，要花费更多的 CPU资源。  

第二种情况导致变慢的原因在于，Redis一次需要返回给客户端的数据过多，更多时间花费在数据协议的组装和网络传输过程中。

### 解决方案

1、尽量不使用 O(N) 以上复杂度过高的命令，对于数据的聚合操作，放在客户端做  
2、执行O(N)命令，保证N尽量的小（推荐 N <= 300），每次获取尽量少的数据，让Redis可以及时处理返回  

## 操作bigkey

如果通过慢查询发现，并不是复杂度过高的命令导致，而是简单的set/del命令，那就可以查看是否写入了biekey；  

redis在写入数据时，需要为新的数分配内存，相对应的，当从redis中删除数据时，它会释放对应的内存空间；  

如果一个key写入的value非常大，那么redis在“分配内存时就会比较耗时”，同样的当删除这个key时，“释放内存也会比较耗时”，这种类型的key，我们一般称之为bigkey；

### 如果已经写入了bigkey，有没有办法扫描出bigkey的分布情况呢？

```golang
[root@eac5f6fb4074 redis]# redis-cli -h 127.0.0.1 -p 6379 --bigkeys -i 0.01

# Scanning the entire keyspace to find biggest keys as well as
# average sizes per key type.  You can use -i 0.1 to sleep 0.1 sec
# per 100 SCAN commands (not usually needed).

[00.00%] Biggest hash   found so far '1101000' with 1 fields
[00.00%] Biggest string found so far 'key:__rand_int__' with 3 bytes
[00.00%] Biggest set    found so far 'db' with 4 members
[00.00%] Biggest string found so far 'counter:__rand_int__' with 6 bytes
[00.00%] Biggest string found so far 'name' with 9 bytes
[00.00%] Biggest string found so far '1101000070' with 10 bytes

-------- summary -------

Sampled 10 keys in the keyspace!
Total key length in bytes is 85 (avg len 8.50)

Biggest string found '1101000070' has 10 bytes
Biggest    set found 'db' has 4 members
Biggest   hash found '1101000' has 1 fields

6 strings with 33 bytes (60.00% of keys, avg size 5.50)
0 lists with 0 items (00.00% of keys, avg size 0.00)
3 sets with 8 members (30.00% of keys, avg size 2.67)
1 hashs with 1 fields (10.00% of keys, avg size 1.00)
0 zsets with 0 members (00.00% of keys, avg size 0.00)
```

从结果我们可以很清晰地看到，每种数据类型所占用最大内存/拥有最多元素key是哪一个；  
此命令的原理就是在redis内部执行了scan命令，遍历整个实例中所有的key，然后对key的类型分别执行strlen/llen/hlen/scard/zcard命令，来获取string类型的长度  

当执行这个命令时，要注意2个问题：

1、对线上实例进行bigkey扫描时，Redis的OPS会突增，为了降低扫描过程中对Redis的影响，最好控制一下扫描的频率，指定-i参数即可，它表示扫描过程中每次扫描后休息的时间间隔，单位是秒  
2、扫描结果中，对于容器类型（List、Hash、Set、ZSet）的key，只能扫描出元素最多的 key。但一个key的元素多，不一定表示占用内存也多，你还需要根据业务情况，进一步评估内存占用情况

### 针对bigkey有哪方面优化呢？

1、业务应用尽量避免写入bigkey  
2、如果使用的Redis4.0以上版本，针对主动删除操作，用UNLINK命令替代DEL，此命令可以把释放key内存的操作，放到后台线程中去执行，从而降低对Redis的影响  
3、如果使用的Redis6.0以上版本，针对主动删除操作，可以开启lazy-free机制（lazyfree-lazy-user-del = yes），在执行DEL命令时，释放内存也会放到后台线程中执行

## 集中过期

如果发现，平时在操作Redis时，并没有延迟很大的情况发生，但在某个时间点突然出现一波延时，其现象表现为：“变慢的时间点很有规律，例如某个整点，或者每间隔多久就会发生一波延迟”。

如果是出现这种情况，那么你需要排查一下，业务代码中是否存在设置大量key集中过期的情况。

### 为什么集中过期会导致 Redis 延迟变大？

Redis的过期数据采用被动过期 + 主动过期两种策略：

1、被动过期：只有当访问某个key时，才判断这个key是否已过期，如果已过期，则从实例中删除；  
2、主动过期：Redis内部维护了一个定时任务，默认每隔100毫秒（1秒10次）就会从全局的过期哈希表中随机取出20个key，然后删除其中过期的key，如果过期key的比例超过了25%，则继续重复此过程，直到过期key的比例下降到25%以下，或者这次任务的执行耗时超过了25毫秒，才会退出循环。

**注意，这个主动过期key的定时任务，是在Redis主线程中执行的。**
![image](./image/redis-slowlog.webp)

### 那遇到这种情况，如何分析和排查？

一般集中过期使用的是expireat/pexpireat命令，你需要在代码中搜索这个关键字。

### 针对集中过期的情况如何避免呢？

1、集中过期key增加一个随机过期时间，把集中过期的时间打散，降低Redis清理过期key的压力；  
2、如果你使用的Redis4.0以上版本，可以开启lazy-free机制，当删除过期key时，把释放内存的操作放到后台线程中执行，避免阻塞主线程；

## 实例内存达到上限

如果你的Redis实例设置了内存上限maxmemory，那么也有可能导致Redis变慢。  

当Redis内存达到maxmemory后，每次写入新的数据之前，Redis必须先从实例中踢出一部分数据，让整个实例的内存维持在maxmemory之下，然后才能把新数据写进来。

### redis支持的淘汰策略

|算法|说明|
|-----|-----|
|allkeys-lru|不管key是否设置了过期，淘汰最近最少访问的key|
|volatile-lru|只淘汰最近最少访问、并设置了过期时间的key|
|allkeys-random|不管key是否设置了过期，随机淘汰key|
|volatile-random|只随机淘汰设置了过期时间的key|
|allkeys-ttl|不管key是否设置了过期，淘汰即将过期的key|
|noeviction|不淘汰任何key，实例内存达到maxmeory后，再写入新数据直接返回错误|
|allkeys-lfu|不管key是否设置了过期，淘汰访问频率最低的key（4.0+版本支持）|
|volatile-lfu|只淘汰访问频率最低、并设置了过期时间 key（4.0+版本支持）|

## fork耗时严重

一般我们会通过RDB快照和AOF追加日志两种方式备份数据，但是两种方式在运行中都会fork子进程；而fork在执行过程中，“主进程需要拷贝自己的内存页表给子进程”，如果这个实例很大，那么拷贝过程也会比较耗时；而且fork过程会消耗大量CPU资源，在完成fork之前，整个redis实例会被阻塞住，无法处理任何客户端请求。

### 如何确认因为fork耗时导致的redis延迟变大呢？

通过info命令查看latest_fork_usec，单位微妙

```golang
latest_fork_usec:0
```

![image](./image/redis_fork.webp)

### 如何避免这种情况呢？

1、控制redis实例的内存，尽量在10G以下，执行fork的耗时与实例大小有关，实例越大，耗时越久；  
2、合理配置数据持久化策略：可以采用主-从-从架构，如果数据不敏感业务，可以关闭AOF和AOD重写；  
3、降低主从库全量同步的概率：适当调大 repl-backlog-size 参数，避免主从全量同步

## 开启内存大页

### 什么是内存大页？

应用程序向操作系统申请内存时，时按“内存页”进行申请的，而常规的内存页大小时4KB；linux内核从2.6.38开始，支持了内存大页机制，该机制允许引用程序以2MB大小为单位，向操作系统申请内存；应用程序每次向操作系统申请的内存单位变大了，这也就意味着申请内存的耗时边长了。

## 开启AOF

## 绑定CPU

## 使用swap

## 碎片整理

## 网络带宽过载

## 其他原因

## 原文

[Redis为什么变慢了？一文讲透如何排查Redis性能问题](https://mp.weixin.qq.com/s/mUibP9pYtQ8Af1ztzt0k-A)
