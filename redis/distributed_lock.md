# 分布式锁

为有效控制多个线程同时访问、调整一个有状态的值时，我们需要一种技术来控制多个线程的安全操作，即为分布式锁

## redis如何实现分布式锁

### 使用set nx实现分布式

#### 注意事项

set命令要用 set key value px milliseconds nx；  
value要具有唯一性；  
释放锁时要验证value值，不能误解锁；  
防止事务执行时间大于设置锁的超时时间，可以设置定时进行检查；  

原理：setnx 如果不存在则设置成功，如果失败则设置失败  可以用来实现上述的分布式锁  
setnx key value  
set key value ex nx  

使用redis做分布式锁，要考虑三个核心要素：加锁、解锁、锁超时等带来的一系列问题  

#### 一、例如最基本的使用方式

```redis
if(setnx(lock_id, 1)){
    ...
    del(lock_id)
}
```

上述过程，如果线程A在获取到锁之后，开始执行业务代码，但是未执行完，就挂了，也就是锁未释放，导致其他的线程再也拿不到锁了，就变成了死锁。

#### 二、添加一个expire解决上面死锁问题

```redis
if(setnx(lock_id, 1)){
    expire(lock_id,30)
    ...
    del(lock_id)
}
```

这样一来，获取到锁的同时，给锁设置一个过期时间，即使线程挂了，超时之后锁也会释放；  
但是同时也会有新的问题  

1、setnx，expire是两个操作，非原子性的，如果执行完setnx，执行expire之前线程挂了还是会出现试过  
解决方案：set(lock_id,1,30, nx) 将两个操作换成一个原子操作即可解决

2、极端情况，线程A执行的时间超过了设置锁的过期时间，这样就会导致其他线程B会提现获取到锁，在A执行完成之后释放锁就会将线程B的锁释放掉  
解决方案：set(lock_id,lock_id,30, nx)，设置值的时候可以设置当前线程ID，在锁删除的时候做一下判断是否相等在执行del操作

### redlock(分布式锁)

## 参考

[利用 Redis 实现分布式锁](https://www.cnblogs.com/jojop/p/14008824.html)  
[七种方案！探讨Redis分布式锁的正确使用姿势](https://z.itpub.net/article/detail/0A3DCC6FF8BD96C478FF1D7644DBFA57)
