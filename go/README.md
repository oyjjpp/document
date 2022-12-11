# Golang常见问题

## 编译原理

[Golang编译原理](./compile.md)

### go的init函数是什么时候执行的？

### 多个init函数执行顺序能保证吗？

## 数据结构

[Golang数据结构](./struct.md)

### 切片的底层实现说一下？

```golang
type SliceHeader struct {
 Data uintptr
 Len  int
 Cap  int
}
```

- Data 是指向数组的指针;  
- Len 是当前切片的长度；  
- Cap 是当前切片的容量，即 Data 数组的大小。

### Go中对nil的Slice和空的Slice的处理是一致的吗？

nil slice和empty slice是不一致的

通常错误的用法，会报数组越界的错误，因为只是声明了slice，却没有给实例化的对象。

```golang
var slice []int
slice[1] = 0
```

此时slice的值是nil，这种情况可以用于需要返回slice的函数，当函数出现异常的时候，保证函数依然会有nil的返回值。

empty slice 是指slice不为nil，但是slice没有值，slice的底层的空间是空的，此时的定义如下：

```golang
slice := make([]int,0）
slice := []int{}
```

当我们查询或者处理一个空的列表的时候，这非常有用，它会告诉我们返回的是一个列表，但是列表内没有任何值。

### slice，len，cap，共享，扩容

- 如果期望容量大于当前容量的两倍就会使用期望容量；
- 如果当前切片的长度小于 1024 就会将容量翻倍；
- 如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；

### slice和array区别

- 数组长度不能改变，初始化后长度就是固定的；切片的长度是不固定的，可以追加元素，在追加时可能使切片的容量增大。
- 结构不同，数组是一串固定数据，切片描述的是截取数组的一部分数据，从概念上说是一个结构体。
- 初始化方式不同，在声明时的时候：声明数组时，方括号内写明了数组的长度或使用...自动计算长度，而声明slice时，方括号内没有任何字符。
- unsafe.sizeof的取值不同，unsafe.sizeof(slice)返回的大小是切片的描述符，不管slice里的元素有多少，返回的数据都是24字节。unsafe.sizeof(arr)的值是在随着arr的元素的个数的增加而增加，是数组所存储的数据内存的大小。

### 如何把数组转化成一个切片

```golang
arrData := [10]int{}
sliceData := arrData[0:3]
```

### make一个slice参数怎么写？

```golang
data := make([]int, 5, 10）
```

当前切片长度=len=5  
当前切片容量=cap=10  

### map会遇到一些并发安全的问题，为什么就并发不安全了？

go 语言的设计者认为，在大部分场景中，对 map 的操作都非线程安全的；**我们不可能为了那小部分的需求，而牺牲大部分人的性能。**  

```golang
func main() {
    s := make(map[int]int)
    // 启 100 个异步线程去写 map
    for i := 0; i < 100; i++ {
        go func(i int) {
            s[i] = i
        }(i)
    }
    // 启 100 个一步线程去读 map
    for i := 0; i < 100; i++ {
        go func(i int) {
            fmt.Printf("map 第 %d 个元素值是 %d", i, s[i])
        }(i)
    }
    // 睡眠 3 秒钟，方便前面的协程打印
    time.Sleep(3 * time.Second)
}
```

因为是异步协程，所以 map 同一时刻就可能被多个写的协程操作。  
那么运行后就会报错：fatal error: concurrent map writes

### map和sync.Map是有什么区别？看过源码吗，可以介绍一下吗？

[一文解决Map并发问题](https://cloud.tencent.com/developer/article/1539049)

#### sync.Map的特点

```golang
type Map struct {
   mu Mutex
   read atomic.Value // readOnly
   dirty map[interface{}]*entry
   misses int
}
```

可以无锁访问read map，而且会优先操作read map，倘若只操作read map就可以满足要求，那就不用去操作write map(dirty)。

**sync.Map 特点：**

- 空间换时间，通过冗余的两个数据结构(read、dirty)，实现加锁对性能的影响。  
- 使用只读数据(read)，避免读写冲突。  
- 动态调整，miss次数多了之后，将dirty数据提升为read。  
- double-checking。  
- 延迟删除，删除一个键值只是打标记，只有在提升dirty的时候才清理删除的数据。  
- 优先从read读取、更新、删除，因为对read的读取不需要锁。  

**sync缺点：**
如果是写多的场景，会导致 read map 缓存失效，需要加锁，冲突变多，性能急剧下降。

### map里面解决hash冲突怎么做的，冲突了元素放在头还是尾

Go 语言中使用拉链法来解决哈希碰撞的问题实现了哈希表

根据key定位到指定桶中，循环桶中的元素

- 1、先定位如果tophash与桶中的top是否不相等 如果不相等则返回地址  
- 2、如果与桶中tophash相等，则判断key是否相等，如果相等则返回地址  
- 3、如果1/2 寻找后还没有合适地址，则去查找溢出桶，循环1/2  
- 4、如果溢出桶也未寻找到，则需要扩容  

### map取一个key，然后修改这个值，原map数据的值会不会变化

- 1、如果单独赋值，并且修改了原map的值时不会发生变化  
- 2、如果将整个map赋值，在修改map，原来map的数据值会发生变化

### map如何顺序读取

因为map通过rang遍历的时候加入了随机种子，所以遍历是无须的，  
如果想通过顺序遍历可以将key放入到slice，通过遍历切片的形式去遍历map

### 并发读写map会发生什么？怎么避免?

fatal error: concurrent map read and map write  
会发生并发读写问题

1、在1.9之前可以定义一个结构体加上一把读写锁来解决并发安全问题；
2、1.9之后官方在sync包定义了Map结构体，如果遇到并发问题，可以使用sync.Map解决。  

### struct能不能比较

[Go 结构体（struct）是否可以比较？](https://segmentfault.com/a/1190000040099215)  
strruct 是否能比较需要看结构体的元素类型

### go结构体和结构体指针的区别

```golang
type MyStruct struct {
    Name string
}

func (s MyStruct) SetName1(name string) {
    s.Name = name
}

func (s *MyStruct) SetName2(name string) {
    s.Name = name
}

 func SetName1(s MyStruct, name string){
    u.Name = name
 }

 func SetName2(s *MyStruct,name string){
    u.Name = name
 }
```

- 在使用上的考虑：方法是否需要修改接收器？如果需要，接收器必须是一个指针。
- 在效率上的考虑：如果接收器很大，比如：一个大的结构体，使用指针接收器会好很多。
- 在一致性上的考虑：如果类型的某些方法必须有指针接收器，那么其余的方法也应该有指针接收器，所以无论类型如何使用，方法集都是一致的。

回到上面的例子中，从功能使用角度来看：  
如果 SetName2 方法修改了 s 的字段，调用者是可以看到这些字段值变更的，因为其是指针引用，本质上是同一份。  
相对 SetName1 方法来讲，该方法是用调用者参数的副本来调用的，本质上是值传递，它所做的任何字段变更对调用者来说是看不见的。
另外对于基本类型、切片和小结构等类型，值接收器是非常廉价的。  

### go里面interface是什么概念

### go什么场景使用接口

### 函数传递有什么区别

参数都是值传递

### 为什么给变量一个基础类型没有并发安全问题？

### string和byte数组有什么区别？

### go深拷贝，什么时候需要深拷贝

#### 深拷贝（Deep Copy）

拷贝的是数据本身，创造一个样的新对象，新创建的对象与原对象不共享内存，新创建的对象在内存中开辟一个新的内存地址，新对象值修改时不会影响原对象值。既然内存地址不同，释放内存地址时，可分别释放。

值类型的数据，默认全部都是深复制，Array、Int、String、Struct、Float，Bool。

#### 浅拷贝（Shallow Copy）

拷贝的是数据地址，只复制指向的对象的指针，此时新对象和老对象指向的内存地址是一样的，新对象值修改时老对象也会变化。释放内存地址时，同时释放内存地址。

引用类型的数据，默认全部都是浅复制，Slice，Map。

## 常用关键字

### defer是啥？怎么用的？底层原理是啥？

#### defer底层原理

**编译期：**

- 将 defer 关键字被转换 runtime.deferproc；  
- 在调用 defer 关键字的函数返回之前插入 runtime.deferreturn；

**运行时：**

- runtime.deferproc 会将一个新的 runtime._defer 结构体追加到当前 Goroutine 的链表头；  
- runtime.deferreturn 会从 Goroutine 的链表中取出 runtime._defer 结构并依次执行；

#### 使用defer的现象

**defer 关键字的调用时机以及多次调用 defer 时执行顺序是如何确定的；**  

- 后调用的 defer 函数会被追加到 Goroutine _defer 链表的最前面；  
- 运行 runtime._defer 时是从前到后依次执行；

**defer 关键字使用传值的方式传递参数时会进行预计算，导致不符合预期的结果；**  

- 调用 runtime.deferproc 函数创建新的延迟调用时就会立刻拷贝函数的参数，函数的参数不会等到真正执行时计算；

**return 不是原子操作：**

- 执行过程是: 保存返回值(若有)–>执行 defer（若有）–>执行 ret 跳转，申请资源后立即使用 defer 关闭资源是好习惯。

### defer用的多吗？有哪些应用

- 资源释放（数据库资源、锁资源）  
- 连接关闭（TCP、文件句柄）  
- 捕获panic， 进行 recover 防止程序崩溃

### panic 和 recover

- panic 能够改变程序的控制流，调用 panic 后会立刻停止执行当前函数的剩余代码，并在当前 Goroutine 中递归执行调用方的 defer；  
- recover 可以中止 panic 造成的程序崩溃。它是一个只能在 defer 中发挥作用的函数，在其他作用域中调用不会发挥作用；  

#### panic 和 recover使用的现象

- panic 只会触发当前 Goroutine 的 defer；  
- recover 只有在 defer 中调用才会生效；  
- panic 允许在 defer 中嵌套多次调用；  

![image](./image/20221130075701001.png)

### go如何避免panic

通过defer 函数中调用recover() 恢复崩溃情况

### 异常捕获是如何做的？

- 逐级返回error等待顶成捕获
- 通过recover在defer捕获

### select可以用于什么

[select](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-select/)

select 能够让 Goroutine 同时等待多个 Channel 可读或者可写，在多个文件或者 Channel状态改变之前，select 会一直阻塞当前线程或者Goroutine。

- select 能在 Channel 上进行非阻塞的收发操作；
- select 在遇到多个 Channel 同时响应时，会随机执行一种情况；

### make和new

- make 的作用是初始化内置的数据结构，我们常用到的切片、哈希表和 Channel；
- new 的作用是根据传入的类型分配一片内存空间并返回指向这片内存空间的指针；

```golang
slice := make([]int, 0, 100)
hash := make(map[int]bool, 10)
ch := make(chan int, 5)
```

slice 是一个包含 data、cap 和 len 的结构体 reflect.SliceHeader；  
hash 是一个指向 runtime.hmap 结构体的指针；  
ch 是一个指向 runtime.hchan 结构体的指针；  

```golang
i := new(int)

var v int
i := &v
```

上述代码片段中的两种不同初始化方法是等价的，它们都会创建一个指向 int 零值的指针。

### append时的过程

针对切片类型持续进行append可能会导致切片进行扩容

校验是否赋值原有变量

**如果不赋值原有变量:**

1、首先获取当前切片的数组指针、大小、容量  
2、如果追加的元素后，长度大于容量，会触发切片的扩容  
3、最后将追加的元素添加到指定的数组中

**如果赋值原有变量:**  

Go语言编译器已经对这种常见的情况做出了优化  
其他步骤与不赋值原有变量逻辑一致  

**append引发的扩容:**

1、如果期望容量大于当前容量的两倍就会使用期望容量；  
2、如果当前切片的长度小于 1024 就会将容量翻倍；  
3、如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；  

以上仅是针对容量大小的简单计算，还会设计内存对齐

需要注意的是在遇到大切片扩容或者复制时可能会发生大规模的内存拷贝，一定要减少类似操作避免影响程序的性能。

## 并发编程

[Golang并发模型](./concurrent.md)

### context

上下文是Go语言中用来设置截止日期、同步信号，传递请求相关值的接口类型；上下文与Goroutine有比较密切的关系；

```golang
type Context interface {
 Deadline() (deadline time.Time, ok bool)
 Done() <-chan struct{}
 Err() error
 Value(key interface{}) interface{}
}
```

主要作用是在多个Goroutine组成的树中同步取消信号以减少对资源的消耗和占用；还有传值的功能，但是这个功能我们还是很少用到。

### Go语言的互斥锁是怎么实现的？读写锁呢？

Go 语言的 sync.Mutex 由两个字段 state 和 sema 组成。其中 state 表示当前互斥锁的状态，而 sema 是用于控制锁状态的信号量。

```golang
type Mutex struct {
 state int32
 sema  uint32
}
```

#### 互斥锁的加锁过程比较复杂，它涉及自旋、信号量以及调度等概念

- 如果互斥锁处于初始化状态，会通过置位 mutexLocked 加锁；
- 如果互斥锁处于 mutexLocked 状态并且在普通模式下工作，会进入自旋，执行 30 次 PAUSE 指令消耗 CPU 时间等待锁的释放；
- 如果当前 Goroutine 等待锁的时间超过了 1ms，互斥锁就会切换到饥饿模式；
- 互斥锁在正常情况下会通过 runtime.sync_runtime_SemacquireMutex 将尝试获取锁的 Goroutine 切换至休眠状态，等待锁的持有者唤醒；
- 如果当前 Goroutine 是互斥锁上的最后一个等待的协程或者等待的时间小于 1ms，那么它会将互斥锁切换回正常模式；

#### 互斥锁的解锁过程与之相比就比较简单，其代码行数不多、逻辑清晰，也比较容易理解

- 当互斥锁已经被解锁时，调用 sync.Mutex.Unlock 会直接抛出异常；
- 当互斥锁处于饥饿模式时，将锁的所有权交给队列中的下一个等待者，等待者会负责设置 mutexLocked 标志位；
- 当互斥锁处于普通模式时，如果没有 Goroutine 等待锁的释放或者已经有被唤醒的 Goroutine 获得了锁，会直接返回；在其他情况下会通过 sync.r- untime_Semrelease 唤醒对应的 Goroutine；

### Go语言的读写锁是怎么实现的？

```golang
type RWMutex struct {
 w           Mutex
 writerSem   uint32
 readerSem   uint32
 readerCount int32
 readerWait  int32
}
```

w — 复用互斥锁提供的能力；  
writerSem 和 readerSem — 分别用于写等待读和读等待写：  
readerCount 存储了当前正在执行的读操作数量；  
readerWait 表示当写操作被阻塞时等待的读操作个数；  

虽然读写互斥锁 sync.RWMutex 提供的功能比较复杂，但是因为它建立在 sync.Mutex 上，所以实现会简单很多。我们总结一下读锁和写锁的关系：

- 调用 sync.RWMutex.Lock 尝试获取写锁时；每次 sync.RWMutex.RUnlock 都会将 readerCount 其减一，当它归零时该 Goroutine 会获得写锁;将 readerCount 减少 rwmutexMaxReaders 个数以阻塞后续的读操作；
- 调用 sync.RWMutex.Unlock 释放写锁时，会先通知所有的读操作，然后才会释放持有的互斥锁；

读写互斥锁在互斥锁之上提供了额外的更细粒度的控制，能够在读操作远远多于写操作时提升性能。

### go的锁是可重入的吗？

不支持

#### Go 锁设计原则

在工程中使用互斥的根本原因是：为了保护不变量，也可以用于保护内、外部的不变量。  
基于此，Go 在互斥锁设计上会遵守这几个原则。如下：

- 在调用 mutex.Lock 方法时，要保证这些变量的不变性保持，不会在后续的过程中被破坏。
- 在调用 mu.Unlock 方法时，要保证：（1）程序不再需要依赖那些不变量。（2）如果程序在互斥锁加锁期间破坏了它们，则需要确保已经恢复了它们。

#### 可重入锁

简单来讲，可重入互斥锁是互斥锁的一种，同一线程对其多次加锁不会产生死锁，又或是导致阻塞。

锁的场景如下：

- 在加锁上：如果是可重入互斥锁，当前尝试加锁的线程如果就是持有该锁的线程时，加锁操作就会成功。
- 在解锁上：可重入互斥锁一般都会记录被加锁的次数，只有执行相同次数的解锁操作才会真正解锁。

### 获取不到锁会一直等待吗？

不会，如果当前 Goroutine 等待锁的时间超过了 1ms，互斥锁就会切换到饥饿模式；

### 那如何实现一个timeout的锁？

[Go超时锁的设计和实现](https://www.jianshu.com/p/4d85661fba0a)

### golang支持哪些并发机制

- CSP模式
- 共享内存

### sync pool的实现原理

### golang sync.WaitGroup用过吗？有哪些坑？

- 1、Add的协程数量和Done数量一样要相等  
- 2、如果需要通过函数传递WaitGroup，一定要传递指针，Go函数是值传递

### go用共享内存的方式实现并发如何保证安全？

### 怎么理解“不要用共享内存来通信，而是用通信来共享内存”

从架构上来讲，降低共享内存的使用，本来就是解耦和的重要手段之一

## 调度器

- [Golang调度器](./scheduler.md)
- [[Golang三关-典藏版] Golang 调度器 GMP 原理与调度全分析](https://learnku.com/articles/41728)
- [6.5 调度器](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-goroutine/)

### Go语言调度器的发展史说一下？

**单线程调度器 · 0.x：**  
只包含 40 多行代码；  
程序中只能存在一个活跃线程，由 G-M 模型组成；

**多线程调度器 · 1.0：**  
允许运行多线程的程序；  
全局锁导致竞争严重；

**任务窃取调度器 · 1.1：**  
引入了处理器 P，构成了目前的 G-M-P 模型；  
在处理器 P 的基础上实现了基于工作窃取的调度器；  
在某些情况下，Goroutine 不会让出线程，进而造成饥饿问题；  
时间过长的垃圾回收（Stop-the-world，STW）会导致程序长时间无法工作；

**抢占式调度器 · 1.2 ~ 至今：**  
基于协作的抢占式调度器 - 1.2 ~ 1.13  
通过编译器在函数调用时插入抢占检查指令，在函数调用时检查当前 Goroutine 是否发起了抢占请求，实现基于协作的抢占式调度；  
Goroutine 可能会因为垃圾回收和循环长时间占用资源导致程序暂停；  
基于信号的抢占式调度器 - 1.14 ~ 至今  
实现基于信号的真抢占式调度；  
垃圾回收在扫描栈时会触发抢占调度；  
抢占的时间点不够多，还不能覆盖全部的边缘情况；  

**非均匀存储访问调度器 · 提案**  
对运行时的各种资源进行分区；  
实现非常复杂，到今天还没有提上日程；  

### GMP模型是什么？为什么要有P？

Go1.0 的 GM 模型的 Goroutine 调度器限制了用 Go 编写的并发程序的可扩展性，尤其是高吞吐量服务器和并行计算程序。

#### GM实现有如下的问题

**存在单一的全局 mutex（Sched.Lock）和集中状态管理：**  
mutex 需要保护所有与 goroutine 相关的操作（创建、完成、重排等），导致锁竞争严重。

**Goroutine 传递的问题：**  
goroutine（G）交接（G.nextg）：工作者线程（M's）之间会经常交接可运行的 goroutine。  
上述可能会导致延迟增加和额外的开销。每个 M 必须能够执行任何可运行的 G，特别是刚刚创建 G 的 M。

**每个 M 都需要做内存缓存（M.mcache）：**  
会导致资源消耗过大（每个 mcache 可以吸纳到 2M 的内存缓存和其他缓存），数据局部性差。  

**频繁的线程阻塞/解阻塞：**  
在存在 syscalls 的情况下，线程经常被阻塞和解阻塞。这增加了很多额外的性能开销。

#### 加入P带来什么改变

- 每个 P 有自己的本地队列，大幅度的减轻了对全局队列的直接依赖，所带来的效果就是锁竞争的减少。而 GM 模型的性能开销大头就是锁竞争。  
- 每个 P 相对的平衡上，在 GMP 模型中也实现了 Work Stealing 算法，如果 P 的本地队列为空，则会从全局队列或其他 P 的本地队列中窃取可运行的 G 来运行，减少空转，提高了资源利用率。

#### 为什么要有 P

- 一般来讲，M 的数量都会多于 P。像在 Go 中，M 的数量默认是 10000，P 的默认数量的 CPU 核数。另外由于 M 的属性，也就是如果存在系统阻塞调用，阻塞了 M，又不够用的情况下，M 会不断增加。
- M 不断增加的话，如果本地队列挂载在 M 上，那就意味着本地队列也会随之增加。这显然是不合理的，因为本地队列的管理会变得复杂，且 Work Stealing 性能会大幅度下降。  
- M 被系统调用阻塞后，我们是期望把他既有未执行的任务分配给其他继续运行的，而不是一阻塞就导致全部停止。

### 主协程如何等其余协程完再操作

主协程自我阻塞，直到需要的协程完成  
阻塞方法  

**使用sync.WaitGroup()管理其余协程：**  
优点：操作简单  
缺点：不能管控协程的执行完成的顺序  

**利用缓存管道进行协程之间的通信：**  
优点：能够管控一组协程结束  
缺点：不能管控协程的执行完成顺序  

**利用无缓存管道进行协程之间的通信：**  
优点：能够管控协程执行完成的顺序  

### Goroutine是什么？能介绍一下它吗？

Goroutine 是 Go 语言调度器中待执行的任务，它在运行时调度器中的地位与线程在操作系统中差不多，但是它占用了更小的内存空间，也降低了上下文切换的开销。

Goroutine 只存在于 Go 语言的运行时，它是 Go 语言在用户态提供的线程，作为一种粒度更细的资源调度单元，如果使用得当能够在高并发的场景下更高效地利用机器的 CPU。

### goroutine创建数量有限制吗？

理论上没有限制，但是因为协程同样消耗CPU、内存等资源，所以是最好限制一下协程数量

**通过channel限制goroutine：**

```golang
type Limit struct {
 n   int
 c   chan struct{}
 ctx context.Context
}

func NewLimit(ctx context.Context, n int) *Limit {
 return &Limit{
  n:   n,
  c:   make(chan struct{}, n),
  ctx: ctx,
 }
}

// Run f in a new goroutine but with limit.
func (g *Limit) Run(f func(ctx context.Context)) {
 g.c <- struct{}{}
 go func() {
  f(g.ctx)
  <-g.c
 }()
}

```

### 为什么不要频繁创建和停止goroutine

### go里面为什么需要多协程？

### goroutine为什么会存在，为什么不使用线程？

### 同时启了一万个g，如何调度的？

### 一个进程能创建多少线程受哪些因素的限制

### goroutine在项目里面主要承担了什么责任

提供并发性能

### 用go协程的时候也是要走IO的，go是如何处理的？

### 为什么不要大量使用goroutine

在Go语言中，goroutine的创建成本很低，调度效率高，Go语言在设计时就是按以数万个goroutine为规范进行设计的，数十万个并不意外，但是goroutine在内存占用方面确实具有有限的成本，你不能创造无限数量的它们

### 并行goroutine如何实现

可以通过管道或者sync.WaitGroup

### go协程、线程、进程区别

​协程跟线程是有区别的，线程由 CPU 调度是抢占式的，协程由用户态调度是协作式的，一个协程让出 CPU 后，才执行下一个协程。

线程和进程属于内核态

协程属于用户态

**协程（coroutine）：**  
和线程类似，共享堆，不共享栈，协程的切换一般由程序员在代码中显式控制。它避免了上下文切换的额外耗费，兼顾了多线程的优点，简化了高并发程序的复杂。

Goroutine和其他语言的协程（coroutine）在使用方式上类似，但从字面意义上来看不同（一个是Goroutine，一个是coroutine），再就是协程是一种协作任务控制机制，在最简单的意义上，协程不是并发的，而Goroutine支持并发的。因此Goroutine可以理解为一种Go语言的协程。同时它可以运行在一个或多个线程上。

**线程（Thread）：**  
有时被称为轻量级进程(Lightweight Process，LWP），是程序执行流的最小单元。一个标准的线程由线程ID，当前指令指针(PC），寄存器集合和堆栈组成。另外，线程是进程中的一个实体，是被系统独立调度和分派的基本单位，线程自己不拥有系统资源，只拥有一点儿在运行中必不可少的资源，但它可与同属一个进程的其它线程共享进程所拥有的全部资源。

线程拥有自己独立的栈和共享的堆，共享堆，不共享栈，线程的切换一般也由操作系统调度。

### 知道processor大小是多少吗？

因为会消耗大量的内存 (进程虚拟内存会占用 4GB [32 位操作系统], 而线程也要大约 4MB)。

大量的进程 / 线程出现了新的问题

- 高内存占用
- 调度的高消耗 CPU

### 服务能开多少个m由什么决定

- go 语言本身的限制：go 程序启动时，会设置 M 的最大数量，默认 10000. 但是内核很难支持这么多的线程数，所以这个限制可以忽略。
- runtime/debug 中的 SetMaxThreads 函数，设置 M 的最大数量
- 一个 M 阻塞了，会创建新的 M。

### 开多少个p由什么决定

由启动时环境变量 $GOMAXPROCS 或者是由 runtime 的方法 GOMAXPROCS() 决定。这意味着在程序执行的任意时刻都只有 $GOMAXPROCS 个 goroutine 在同时运行。

### m和p是什么样的关系

M 与 P 的数量没有绝对关系，一个 M 阻塞，P 就会去创建或者切换另一个 M，所以，即使 P 的默认数量是 1，也有可能会创建很多个 M 出来。

### P 和 M 何时会被创建

- P 何时创建：在确定了 P 的最大数量 n 后，运行时系统会根据这个数量创建 n 个 P。

- M 何时创建：没有足够的 M 来关联 P 并运行其中的可运行的 G。比如所有的 M 此时都阻塞住了，而 P 中还有很多就绪任务，就会去寻找空闲的 M，而没有空闲的，就会去创建新的 M。

## 内存管理

[Golang内存管理](./memory.md)  
[Go内存分配那些事，就这么简单！](https://mp.weixin.qq.com/s/3gGbJaeuvx4klqcv34hmmw)  

### 知道golang的内存逃逸吗？什么情况下会发生内存逃逸？

**Go 语言的逃逸分析遵循以下两个不变性：**

- 指向栈对象的指针不能存在于堆中；
- 指向栈对象的指针不能在栈对象回收后存活；

程序变量会携带有一组校验数据，用来证明它的整个生命周期是否在运行时完全可知。
如果变量通过了这些校验，它就可以在栈上分配。否则就说它逃逸了，必须在堆上分配。

能引起变量逃逸到堆上的典型情况：

- 在方法内把局部变量指针返回局部变量原本应该在栈中分配，在栈中回收。但是由于返回时被外部引用，因此其生命周期大于栈，则溢出。
- 发送指针或带有指针的值到 channel 中。 在编译时，是没有办法知道哪个 goroutine 会在 channel 上接收数据。所以编译器没法知道变量什么时候才会被释放。
- 在一个切片上存储指针或带指针的值；一个典型的例子就是 []*string ，这会导致切片的内容逃逸。尽管其后面的数组可能是在栈上分配的，但其引用的值一定是在堆上。
- slice 的背后数组被重新分配了，因为 append 时可能会超出其容量( cap )。 slice 初始化的地方在编译时是可以知道的，它最开始会在栈上分配。- 如果切片背后的存储要基于运行时的数据进行扩充，就会在堆上分配。
- 在 interface 类型上调用方法。 在 interface 类型上调用方法都是动态调度的 —— 方法的真正实现只能在运行时知道。想像一个 io.Reader 类型的变量 r , 调用 r.Read(b) 会使得 r 的值和切片b 的背后存储都逃逸掉，所以会在堆上分配。

**案例：**
通过一个例子加深理解，接下来尝试下怎么通过 go build -gcflags=-m 查看逃逸的情况。

```golang

type A struct {
 s string
}

// 这是上面提到的 "在方法内把局部变量指针返回" 的情况
func foo(s string) *A {
 a := new(A)
 a.s = s
 return a //返回局部变量a,在C语言中妥妥野指针，但在go则ok，但a会逃逸到堆
}

func main() {
 a := foo("hello")
 b := a.s + " world"
 c := b + "!"
 fmt.Println(c)
}

PS D:\project\algorithm> go build -gcflags=-m main.go
# command-line-arguments
.\main.go:12:6: can inline foo
.\main.go:19:10: inlining call to foo
.\main.go:22:13: inlining call to fmt.Println
.\main.go:12:10: leaking param: s
.\main.go:13:10: new(A) escapes to heap
.\main.go:19:10: new(A) does not escape
.\main.go:20:11: a.s + " world" does not escape
.\main.go:21:9: b + "!" escapes to heap
.\main.go:22:13: ... argument does not escape
.\main.go:22:13: c escapes to heap
```

- ./main.go:13:10: new(A) escapes to heap 说明 new(A) 逃逸了,符合上述提到的常见情况中的第一种。
- ./main.go:20:11: main a.s + " world" does not escape 说明 b 变量没有逃逸，因为它只在方法内存在，会在方法结束时被回收。
- ./main.go:21:9: b + "!" escapes to heap 说明 c 变量逃逸，通过fmt.Println(a ...interface{})打印的变量，都会发生逃逸

### go内存操作也要处理IO，是如何处理的?

### 项目中如果出现内存泄漏你是怎么排查的？

内存泄露，是指程序在申请内存并且用完这块内存后（对象不再需要了），没有释放已申请的内存空间。
少数偶然的内存泄漏，虽然不太好，但问题不大，我们也不至于对那点内存抠抠搜搜的。
但如果是内存不断泄漏，直到新的对象没有足够的空间生成，就会导致OOM（Out Of Memory）。

一个健康的程序应该有平稳的新陈代谢，内存占用应该维持在一定范围；但如果内存持续飙升，甚至到达了一个危险的值，那么可以怀疑有内存泄漏。

**如何发现：**
肯定要借助监控手段了，针对程序的运行状态加一些内存，CPU，Go协程等关键指标监控，如果发现异常立马追踪

**追踪问题：**
1、如果刚上线引起的可以review代码，查看是否有问题；  
2、借助pprof包 分析内存具体消耗过高的代码位置，反向分析出现问题的位置。

### 内存申请上有什么区别

## 垃圾回收

[深入浅出垃圾回收（三）增量式 GC](https://liujiacai.net/blog/2018/08/04/incremental-gc/)  
[Go 垃圾回收（一）——为什么要学习 GC ?](https://zhuanlan.zhihu.com/p/101132283)  
[Go内存分配那些事，就这么简单！](https://mp.weixin.qq.com/s/3gGbJaeuvx4klqcv34hmmw)  
[Go垃圾回收 1：历史和原理](https://lessisbetter.site/2019/10/20/go-gc-1-history-and-priciple/)

### GC 不回收什么？

为了解释垃圾回收是什么，我们先来说说 GC 不回收什么。在我们程序中会使用到两种内存，分别为堆（Heap）和栈（Stack），而 GC 不负责回收栈中的内存。那么这是为什么呢？

主要原因是栈是一块专用内存，专门为了函数执行而准备的，存储着函数中的局部变量以及调用栈。除此以外，栈中的数据都有一个特点——简单。比如局部变量就不能被函数外访问，所以这块内存用完就可以直接释放。正是因为这个特点，栈中的数据可以通过简单的编译器指令自动清理，也就不需要通过 GC 来回收了。

### 为什么需要垃圾回收？

现在我们知道了垃圾回收只负责回收堆中的数据，那么为什么堆中的数据需要自动垃圾回收呢？

其实早期的语言是没有自动垃圾回收的。比如在 C 语言中就需要使用 malloc/free 来人为地申请或者释放堆内存。这种做法除了增加工作量以外，还容易出现其他问题。

一种可能是并发问题，并发执行的程序容易错误地释放掉还在使用的内存。  
一种可能是重复释放内存，还有可能是直接忘记释放内存，从而导致内存泄露等问题。

而这类问题不管是发现还是排查往往会花费很多时间和精力。
所以现代的语言都有了这样的需求——一个自动内存管理工具。

### 什么是垃圾回收？

当我们说垃圾回收（GC garbage collection）的时候，我们其实说的是自动垃圾回收（Automatic Garbage Collection），一个自动回收堆内存的工具。所以垃圾回收一点也不神奇，它只是一种工具，可以更便捷更高效地帮助程序员管理内存。

### 追踪式垃圾回收（Tracing garbage collection）

追踪式算法的核心思想是判断一个对象是否可达，因为一旦这个对象不可达就可以立刻被 GC 回收了。  
那么我们怎么判断一个对象是否可达呢？  
第一步找出所有的全局变量和当前函数栈里的变量，标记为可达。  
第二步从已经标记的数据开始，进一步标记它们可访问的变量，以此类推，专业术语叫传递闭包。

### GC你了解吗？展开说一下GC的过程

一次完整的垃圾回收会分为四个阶段，分别是标记准备、标记、结束标记以及清理。  
在标记准备和标记结束阶段会需要 STW，标记阶段会减少程序的性能，而清理阶段是不会对程序有影响的

#### 阶段一：Mark Setup 标记准备（STW：Stop the world）

**Write Barrier（写屏障）：**
我们知道三色标记法是一种可以并发执行的算法。  
所以在运行过程中程序的函数栈内可能会有新分配的对象，那么这些对象该怎么通知到 GC，怎么给他们着色呢？这个时候就需要我们的 Write Barrier 出马了。  
Write Barrier 主要做这样一件事情，修改原先的写逻辑，然后在对象新增的同时给它着色，并且着色为”灰色“。  
因此打开了 Write Barrier 可以保证了三色标记法在并发下安全正确地运行。

**Stop The World：**
不过在打开 Write Barrier 前有一个依赖，我们需要先停止所有的 goroutine，也就是所说的 STW（Stop The World）操作。那么接下来问题来了，GC 该怎么通知所有的 goroutine 停止呢 ？

我们知道，在停止 goroutine 的方案中，Go 语言采取的是协助式抢占模式（当前 1.13 及之前版本，1.14开始基于行好事抢占模式）。  
协助模式的做法是在程序编译阶段注入额外的代码，更精确的说法是在每个函数的序言中增加一个协助式抢占点。  
因为一个 goroutine 中通常有无数调用函数的操作，选择在函数序言中增加抢占点可以较好地平衡性能和实时性之间的利弊。  
在通常情况下，一次 Mark Setup 操作会在 10-30 微秒之间。

#### 阶段二：Marking 标记（Concurrent）

在第一阶段打开 Write Barrier 后，就进入第二阶段的标记了。Marking 使用的算法就是我们之前提到的三色标记法。不过我们可以简单了解一下标记阶段的资源分配情况。

在标记开始的时候，收集器会默认抢占 25% 的 CPU 性能，剩下的75%会分配给程序执行。但是一旦收集器认为来不及进行标记任务了，就会改变这个 25% 的性能分配。这个时候收集器会抢占程序额外的 CPU，这部分被抢占 goroutine 有个名字叫 Mark Assist。而且因为抢占 CPU的目的主要是 GC 来不及标记新增的内存，那么抢占正在分配内存的 goroutine 效果会更加好，所以分配内存速度越快的 goroutine 就会被抢占越多的资源。

除此以外 GC 还有一个额外的优化，一旦某次 GC 中用到了 Mark Assist，下次 GC 就会提前开始，目的是尽量减少 Mark Assist 的使用，从而避免影响正常的程序执行。

#### 阶段三：Mark Termination 标记结束（STW）

最重要的 Marking 阶段结束后就会进入 Mark Termination 阶段。这个阶段会关闭掉已经打开了的 Write Barrier，和 Mark Setup 阶段一样这个阶段也需要 STW。

标记结束阶段还需要做的事情是计算下一次清理的目标和计划，比如第二阶段使用了 Mark Assist 就会促使下次 GC 提早进行。如果想人为地减少或者增加 GC 的频率，那么我们可以用 GOGC 这个环境变量设置。Go 的 GC 有且只会有一个参数进行调优，也就是我们所说的 GOGC，目的是为了防止大家在一大堆调优参数中摸不着头脑。

通常情况下，标记结束阶段会耗时 60-90 微秒。

#### 阶段四：Sweeping 清理（Concurrent）

最后一个阶段就是垃圾清理阶段，这个过程是并发进行的。清扫的开销会增加到分配堆内存的过程中，所以这个时间也是无感知不会与垃圾回收的延迟相关联。

### 触发GC时机

- 在分配内存时，会判断当前的Heap内存分配量是否达到了触发一轮GC的阈值（每轮GC完成后，该阈值会被动态设置），如果超过阈值，则启动一轮GC。

- 调用runtime.GC()强制启动一轮GC。

- sysmon是运行时的守护进程，当超过 forcegcperiod (2分钟)没有运行GC会启动一轮GC。

### GC调节参数

Go垃圾回收不像Java垃圾回收那样，有很多参数可供调节，Go为了保证使用GC的简洁性，只提供了一个参数GOGC。

GOGC代表了占用中的内存增长比率，达到该比率时应当触发1次GC，该参数可以通过环境变量设置。

它的单位是百分比，取值范围并不是 [0, 100]，可以是1000，甚至2000，2000时代表2000%，即20倍。

```golang
假如当前heap占用内存为4MB，GOGC = 75，

4 * (1+75%) = 7MB
等heap占用内存大小达到7MB时会触发1轮GC。
```

**GOGC还有2个特殊值：**  

- "off" : 代表关闭GC
- 0 : 代表持续进行垃圾回收，只用于调试

## channel

[一文读懂channel设计](https://mp.weixin.qq.com/s/buwtTCm_szzeusgxHWmKjQ)

**目前的 Channel 收发操作均遵循了先进先出的设计，具体规则如下：**

- 先从 Channel 读取数据的 Goroutine 会先接收到数据；
- 先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利；

### channel底层是用什么实现的？

```golang
type hchan struct {
 qcount   uint // Channel 中的元素个数
 dataqsiz uint // Channel 中的循环队列的长度
 buf      unsafe.Pointer // Channel 的缓冲区数据指针
 closed   uint32
 elemsize uint16 // Channel 能够收发的元素大小
 elemtype *_type // Channel 能够收发的元素类型
 sendx    uint // Channel 的发送操作处理到的位置
 recvx    uint // Channel 的接收操作处理到的位置
 recvq    waitq // Channel 由于缓冲区空间不足而阻塞的 Goroutine 列表
 sendq    waitq // Channel 由于缓冲区空间不足而阻塞的 Goroutine 列表

 lock mutex
}
```

[image](./image/202212082203001.png)

### channel和锁对比一下

并发问题可以用channel解决也可以用Mutex解决，但是它们的擅长解决的问题有一些不同；  
channel关注的是并发问题的数据流动，适用于数据在多个协程中流动的场景；  
而mutex关注的是数据不动，某段时间只给一个协程访问数据的权限，适用于数据位置固定的场景。

### channel的应用场景

- 数据交流：当作并发的 buffer 或者 queue，解决生产者 - 消费者问题。多个 goroutine 可以并发当作生产者（Producer）和消费者（Consumer）。
- 数据传递：一个goroutine将数据交给另一个goroutine，相当于把数据的拥有权托付出去。
- 信号通知：一个goroutine可以将信号(closing，closed，data ready等)传递给另一个或者另一组goroutine。
- 任务编排：可以让一组goroutine按照一定的顺序并发或者串行的执行，这就是编排功能。
- 锁机制：利用channel实现互斥机制。

### 向为nil的channel发送数据会怎么样

```golang
if c == nil {
    if !block {
        return false
    }
    gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
    throw("unreachable")
}
```

往一个nil的channel中发送数据时，调用gopark函数将当前执行的goroutine从running态转入waiting态。

### 同一个协程里面，对无缓冲channel同时发送和接收数据有什么问题

```golang
fatal error: all goroutines are asleep - deadlock!
```

同一个协程里，不能对无缓冲channel同时发送和接收数据，如果这么做会直接报错**死锁**。

对于一个无缓冲的channel而言，只有不同的协程之间一方发送数据一方接受数据才不会阻塞。channel无缓冲时，发送阻塞直到数据被接收，接收阻塞直到读到数据。

### go利用channel通信的方式

### channel有缓冲和无缓冲在使用上有什么区别？

#### 无缓冲是同步的

例如 make(chan int)，就是一个送信人去你家门口送信，你不在家他不走，你一定要接下信，他才会走，无缓冲保证信能到你手上。

#### 有缓冲是异步的

例如 make(chan int, 1)，就是一个送信人去你家仍到你家的信箱，转身就走，除非你的信箱满了，他必须等信箱空下来，有缓冲的保证信能进你家的邮箱。

### 关闭channel有什么作用？

1、资源回收  
2、通知其他协程  

- 在不改变 channel 自身状态的情况下，无法获知一个 channel 是否关闭。
- 关闭一个 closed channel 会导致 panic。所以，如果关闭 channel 的一方在不知道 channel 是否处于关闭状态时就去贸然关闭 channel 是很危险的事情。
- 向一个 closed channel 发送数据会导致 panic。所以，如果向 channel 发送数据的一方不知道 channel 是否处于关闭状态时就去贸然向 channel 发送数据是很危险的事情。

**关于close channel准则：**

- 1)不要在读取端关闭 channel ，因为写入端无法知道 channel 是否已经关闭，往已关闭的 channel 写数据会 panic ；
- 2)有多个写入端时，不要在写入端关闭 channle ，因为其他写入端无法知道 channel 是否已经关闭，关闭已经关闭的 channel 会发生 panic ；
- 3)如果只有一个写入端，可以在这个写入端放心关闭 channel 。

### Channel 可能会引发 goroutine 泄漏

泄漏的原因是 goroutine 操作 channel 后，处于发送或接收阻塞状态，而 channel 处于满或空的状态，一直得不到改变。同时，垃圾回收器也不会回收此类资源，进而导致 gouroutine 会一直处于等待队列中，不见天日。

另外，程序运行过程中，对于一个 channel，如果没有任何 goroutine 引用了，gc 会对其进行回收操作，不会引起内存泄漏。

### go channel实现排序

使用channel进行通信通知，用channel去传递信息，从而控制并发执行顺序。

```golang


var wg sync.WaitGroup

func main() {
 event1 := make(chan struct{}, 1)
 event2 := make(chan struct{}, 1)
 event3 := make(chan struct{}, 1)

 event1 <- struct{}{}
 wg.Add(3)
 start := time.Now().Unix()
 go Handle("event1", event1, event2)
 go Handle("event2", event2, event3)
 go Handle("event3", event3, event1)
 wg.Wait()

 end := time.Now().Unix()
 fmt.Println(end - start)
}

func Handle(event string, inputchan chan struct{}, outputchan chan struct{}) {
 for i := 0; i < 3; i++ {
  time.Sleep(1 * time.Second)
  select {
  case <-inputchan:
   fmt.Println(event)
   outputchan <- struct{}{}
  }
 }
 wg.Done()
}


event1
event2
event3
event1
event2
event3
event1
event2
event3
3
```

### 集群用channel如何实现分布式锁

## 常用框架

### golang用到哪些框架

|框架|优点|缺点|参考|
|-|-|-|-|
|iris|功能相对完善|依赖库较多【127】|[中文参考](https://www.bookstack.cn/read/studyiris-doc/bb3ca1bc8612e3b2.md) [iris框架解析](https://juejin.cn/post/6844903877507055630)|
|gin|代码较简洁，社区活跃||<https://github.com/gin-gonic/gin>|
|beego||性能差，代码冗余|<https://github.com/astaxie/beego> <https://beego.me>|
|echo|||<https://github.com/labstack/echo> <https://echo.labstack.com>|
|revel|||<https://github.com/revel/revel><https://revel.github.io>|

### gin框架的路由是怎么处理的？

[Golang-gin框架路由原理](https://zhuanlan.zhihu.com/p/491337692)  
[gin框架路由理论](https://www.cnblogs.com/randysun/p/15841366.html#gallery-1)

#### 前缀树算法

前缀树的本质就是一棵查找树，有别于普通查找树，它适用于一些特殊场合，比如用于字符串的查找。比如一个在路由的场景中，有1W个路由字符串，每个字符串长度不等，我们可以使用数组来储存，查找的时间复杂度是O(n)，可以用map来储存，查找的复杂度是O(1)，但是都没法解决动态匹配的问题，如果用前缀树时间复杂度是O(logn)，也可以解决动态匹配参数的问题。

下图展示了前缀树的原理，有以下6个字符串，如果要查找cat字符串，步骤如下：

先拿字符c和root的第一个节点a比较，如果不等，再继续和父节点root的第二个节点比较，直到找到c。  
再拿字符a和父节点c的第一个节点a比较，结果相等，则继续往下。  
再拿字符t和父节点a的第一个节点t比较，结果相等，则完成。  
![image](./image/202212102352001.jpg)

同理，在路由中，前缀树可以规划成如下：

![image](./image/202212102352002.jpg)

## 性能

[Golang 大杀器之跟踪剖析 trace](https://juejin.cn/post/6844903887757901831)  
[腾讯 Go 性能优化实战](https://mp.weixin.qq.com/s/Z9DoVGwdAtpbjealQLEMkw)  
[深度解密Go语言之pprof](https://juejin.cn/post/6844903992720359432)  
[graphviz 图形页面分析工具](https://graphviz.gitlab.io/download/)  
[golang pprof 实战](https://blog.wolfogre.com/posts/go-ppof-practice/)  
[实战Go内存泄露](https://segmentfault.com/a/1190000019222661)

### go性能分析工具

### go的profile工具？

pprof库就可以分析程序的运行情况，并且可以提供可视化的功能

**收集数据方式：**  
1、runtime/pprof： 对于只跑一次的程序，调用pprof包提供的函数，手动开启性能数据采集  
2、net/http/pprof： 对于HTTP服务，访问pprof提供的HTTP接口，获取性能数据  
3、go test：使用go test -run=NONE -bench . -memprofile=mem.out  

// 收集cpu数据 默认30s  
go tool pprof <http://hostname/debug/pprof/profile?seconds=120>  
// 收集heap数据  
go tool pprof <http://hostname/debug/pprof/heap>  
// 收集goroutine数据  
go tool pprof <http://hostname/debug/pprof/goroutine>  
// 收集block数据  
go tool pprof <http://hostname/debug/pprof/block>  
// 收集mutex数据  
go tool pprof <http://hostname/debug/pprof/mutex>  
// 收集trace数据  
curl <http://hostname/debug/pprof/trace?seconds=10> >   trace.out

**数据分析使用：**  
生成报告，Web可视化界面、交互式终端三种方式来使用pprof  
go tool pprof options binary  
--text    ：纯文本方式  
--web     ：生成svg并用浏览器打开  
--svg     ：只生成svg  
--list funcname : 筛选出正则匹配的funcname的函数信息  

// 对比两个文件  
go tool pprof -base profile1 profile2  

**数据收集原理：**  
当CPU性能分析启用后，Go runtime会每10ms暂停一下，记录当前运行的goroutine的调用堆栈及相关数据，当性能分析数据保存到硬盘后，我们就可以分析代码运行状态。

内存性能分析则是在堆（Heap）分配的时候，记录一下调用堆栈，默认情况下是每1000次分配取样一次，这个数值可以改变；栈(Stack)分配由于会随时释放，因此不会被内存分析所记录；由于内存分析是取样方式，并且也因为其记录的是“分配内存”，而不是“使用内存”，因此使用内存能分析工具来准确判断程序具体的内存使用是比较困难的。

阻塞分析是一个很独特的分析，它有点儿类似于CPU性能分析，但是它所记录的是goroutine等待资源所花的时间，阻塞分析对分析程序并发瓶颈非常有帮助，阻塞性能分析可以显示出什么时候出现了大批的goroutine被阻塞了；阻塞性能分析是特殊的分析工具，在排除CPU和内存瓶颈前，不应该用它来分析。

**内存分析案例：**

```golang
$ go tool pprof mem.out
Type: alloc_space
Time: May 1, 2020 at 12:16am (CST)
Entering interactive mode (type "help" for  commands, "o" for options)
(pprof) top 5
Showing nodes accounting for 1084.51MB, 100% of  1084.51MB total
Showing top 5 nodes out of 14
      flat  flat%   sum%        cum   cum%
  490.03MB 45.18% 45.18%   490.03MB 45.18%   strings.(*Builder).WriteString
  306.18MB 28.23% 73.42%   306.18MB 28.23%   bytes.makeSlice
  198.50MB 18.30% 91.72%   198.50MB 18.30%   strconv.formatBits
   89.79MB  8.28%   100%    89.79MB  8.28%   bytes.(*Buffer).String
         0     0%   100%   306.18MB 28.23%   bytes.(*Buffer).WriteString


pprof数值说明
flat ：采样时，该函数所占内存或时间(不包含函数等待子函数返回)
flat%  ： flat/总采样值
sum% ：前面所有flat%的累加值
cum  ：采样时，该函数所占内存或时间(包含函数等待子函数返回)
cum% ：cum/总采样值
```

|类型|描述|备注|
|-|-|-|
|allocs|内存分配情况的采集信息||
|blocks|阻塞操作情况的采集信息||
|cmdline|显示程序启动命令及参数||
|goroutine|当前所有协程的堆栈信息||
|heap|堆上内存使用情况的采集信息|与allocs采样信息一致，allocs是所有对象的内存分配，heap则是活跃对象的内存分配|
|mutex|锁争用情况的采样信息||
|profile|CPU占用情况的采样信息||
|threadcreat|系统线程创建情况的采样信息||
|trace|程序运行跟踪信息||

### 用火焰图的优势？

### 火焰图怎么来寻找瓶颈的？

### 说说火焰图？如何分析的？

## 常见功能

### go实现不重启热部署

### protobuf为什么快

### client如何实现长连接

### 怎么检查go问题

### go怎么实现封装继承多态

### 如何拿到多个goroutine的返回值，如何区别他们

### go里面比较成熟的日志框架了解过没有

### go实现一个并发限制爬虫

### 如何通过goclient写代码获取

### 写个channel相关的题，并发模型，爬虫url，控制并发量

### 参数检查中间件核心功能有哪些？

### 对go的中间件和工作机制有了解吗？

### go使用中遇到的问题

### cgo了解过引入的风险点吗？

### 用go实现一个协程池，大概用什么实现

## go语言如何实现服务不重启热部署

## 你觉得java和golang有什么优势劣势？

## 一个二维数组，行遍历快还是列遍历快，为什么？

1、CPU高速缓存：在计算机系统中，CPU高速缓存是用于减少处理器访问内存所需平均时间的部件。在金字塔式存储体系中它位于自顶向下的第二层，仅次于CPU寄存器。其容量远小于内存，但速度却可以接近处理器的频率。当处理器发出内存访问请求时，会先查看缓存内是否有请求数据。如果存在（命中），则不经访问内存直接返回该数据；如果不存在（失效），则要先把内存中的相应数据载入缓存，再将其返回处理器。缓存之所以有效，主要是因为程序运行时对内存的访问呈现局部性（Locality）特征。这种局部性既包括空间局部性（Spatial Locality），也包括时间局部性（Temporal Locality）。有效利用这种局部性，缓存可以达到极高的命中率。  

2、缓存从内存中抓取一般都是整个数据块，所以它的物理内存是连续的，几乎都是同行不同列的，而如果内循环以列的方式进行遍历的话，将会使整个缓存块无法被利用，而不得不从内存中读取数据，而从内存读取速度是远远小于从缓存中读取数据的。随着数组元素越来越多，按列读取速度也会越来越慢。
