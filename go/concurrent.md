# 并发编程

## 同步原语和锁

锁是一种并发编程中的同步原语（Synchronization Primitives），它能保证多个Goroutine在访问同一片内存时不会出现竞争条件（Race condition）等问题。

### 基本原语

```golang
sync.Mutex
sync.RWmutex
sync.WaitGroup
sync.Once
sync.Cond
```

### 竞争检测机制

```golang
go test -race main
go run -race main.go
go build -race main
```

### 提供哪些方式安全读写共享变量

1、atomic包提供了原子性内存原语  
2、可以使用channel在go协程中安全读写共享变量  
3、sync包提供了互斥锁、读写锁

## channel

Channel是支撑Go语言高性能并发编程模型的重要结构

### channel作用？

goroutine之间的通信方式

### channel如何做到线程安全？

使用互斥锁解决

```golang
type hchan struct {
 qcount   uint
 dataqsiz uint
 buf      unsafe.Pointer
 elemsize uint16
 closed   uint32
 elemtype *_type
 sendx    uint
 recvx    uint
 recvq    waitq
 sendq    waitq

 lock mutex
}
```

### channel发送和接收是否同步

```golang
ch := make(chan int) // 无缓冲的channel由于没有缓冲发送和接收需要同步
ch := make(chan int, 2) // 有缓冲channel不要求发送和操作同步
```

channel无缓冲时：发送阻塞直到数据被接收，接收阻塞直到读到数据  
channel有缓冲时：当缓冲满时发送阻塞，当缓冲空时接收阻塞  

## golang中常用的并发模型

1、使用同步原语和锁实现的共享内存的并发模型  
2、基于channel的CSP模型

## CSP

通信顺序过程（CSP）是一种描述并发系统中交互模式的正式语言，它是并发数据理论家族中的一个成员，被称为过程算法，是基于消息的通道传递理论。

不同传统的多线程共享内存来通信，CPS讲究的是“以通信的方式来共享内存”；用于描述多个独立的并发实体通过共享的通讯管道（channel）进行通信的并发模型，CSP中channel是第一类对象，它不关注发送消息的实体，而关注与发送消息时使用的channel。

## 调度器

[Golang调度器](./scheduler.md)

## 网络轮询器

[Golang网络轮询器](./network_poller.md)
