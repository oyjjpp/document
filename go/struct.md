# 数据结构

## 数组

### 创建方式

```golang
data := [3]int{1,2,3}
data := [...]int{1,2,3}
```

[...]T 的方式声明数组时，会在编译时进行推到，最终都会调用NewArray()函数创建数组  

类型、大小会在编译阶段确认  
当前数组是否应该在堆栈中初始化也在编译期就确定的

### 语句转换

当元素数量小于或者等于4个时，会直接将数组中的元素放置在栈上(不考虑逃逸分析)；  
当元素数量大于 4 个时，会将数组中的元素放置到静态区并在运行时取出；  

## 切片

### 创建切片

```golang
data := array[0:3]
data := []int{1, 2, 3}
data := make([]int, 3)
```

1、使用下标方式获得数组或切片的一部分  
2、使用字面量初始化切片  
3、使用关键字创建切片

切片内元素的类型都是在编译期间确

### 切片的数据结构

```golang
type SliceHeader struct {
 Data uintptr
 Len  int
 Cap  int
}
```

![image](./image/golang-slice-struct.png)
Data 是指向数组的指针;  
Len 是当前切片的长度；  
Cap 是当前切片的容量，即Data数组的大小：

### 追加和扩容

如果期望容量大于当前容量的两倍就会使用期望容量；  
如果当前切片的长度小于 1024 就会将容量翻倍；  
如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；  

### 切片总结

切片的很多功能都是由运行时实现的，无论是初始化切片，还是对切片进行追加或扩容都需要运行时的支持，需要注意的是在遇到大切片扩容或者复制时可能会发生大规模的内存拷贝，一定要减少类似操作避免影响程序的性能。

## 哈希表

golang中的哈希表就是我们经常使用的map

### 设计原理

哈希函数、哈希冲突

#### 哈希冲突如何解决？

##### 链路方法

##### 开放寻址法

### map数据结构

```golang

// A header for a Go map.
type hmap struct {
 // Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
 // Make sure this stays in sync with the compiler's definition.
 count     int // # live cells == size of map.  Must be first (used by len() builtin)
 flags     uint8
 B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
 noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
 hash0     uint32 // hash seed

 buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
 oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
 nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

 extra *mapextra // optional fields
}

// mapextra holds fields that are not present on all maps.
type mapextra struct {
 // If both key and elem do not contain pointers and are inline, then we mark bucket
 // type as containing no pointers. This avoids scanning such maps.
 // However, bmap.overflow is a pointer. In order to keep overflow buckets
 // alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
 // overflow and oldoverflow are only used if key and elem do not contain pointers.
 // overflow contains overflow buckets for hmap.buckets.
 // oldoverflow contains overflow buckets for hmap.oldbuckets.
 // The indirection allows to store a pointer to the slice in hiter.
 overflow    *[]*bmap
 oldoverflow *[]*bmap

 // nextOverflow holds a pointer to a free overflow bucket.
 nextOverflow *bmap
}

// A bucket for a Go map.
type bmap struct {
 // tophash generally contains the top byte of the hash value
 // for each key in this bucket. If tophash[0] < minTopHash,
 // tophash[0] is a bucket evacuation state instead.
 tophash [bucketCnt]uint8
 // Followed by bucketCnt keys and then bucketCnt elems.
 // NOTE: packing all the keys together and then all the elems together makes the
 // code a bit more complicated than alternating key/elem/key/elem/... but it allows
 // us to eliminate padding which would be needed for, e.g., map[int64]int8.
 // Followed by an overflow pointer.
}

type bmap struct {
    topbits  [8]uint8
    keys     [8]keytype
    values   [8]valuetype
    pad      uintptr
    overflow uintptr
}
```

![image](./image/hmap-and-buckets.png)

count ： 代表哈希表中元素数量  
B ： 代表哈希桶的数量 len(buckets)  = 2^B  
hash0 ： 哈希系数，创建哈希表时确定，使用哈希函数时作为参数使用  
buckets ： 哈希桶  
oldbuckets ： 溢出桶，一般在扩容时保存buckets时使用，时buckets一半的大小  

### 为什么需要oldbuckets？

当单个桶的容量超过8(bucketCnt)个时，就可能启用溢出桶保存数据

### golang中如何解决哈希冲入？

Go 语言使用拉链法来解决哈希碰撞的问题实现了哈希表

溢出时如何确定时哪个值？

### 什么情况会出现扩容？

装载因子已经超过 6.5；  
哈希使用了太多溢出桶；  

### map总结

Go 语言使用拉链法来解决哈希碰撞的问题实现了哈希表，它的访问、写入和删除等操作都在编译期间转换成了运行时的函数或者方法。哈希在每一个桶中存储键对应哈希的前 8 位，当对哈希进行操作时，这些 tophash 就成为可以帮助哈希快速遍历桶中元素的缓存。

哈希表的每个桶都只能存储 8 个键值对，一旦当前哈希的某个桶超出 8 个，新的键值对就会存储到哈希的溢出桶中。随着键值对数量的增加，溢出桶的数量和哈希的装载因子也会逐渐升高，超过一定范围就会触发扩容，扩容会将桶的数量翻倍，元素再分配的过程也是在调用写操作时增量进行的，不会造成性能的瞬时巨大抖动。

## 接口

Go语言只会在传递参数、返回参数以及变量赋值时才会对某个类型是否实现接口进行检查  
Go语言中有两种略微不同的接口，一种是带有一组方法的接口，另一种是不带任何方法的interface{}  
使用runtime.iface表示第一种接口，使用runtime.eface表示第二种不包含任何方法的接口interface{}
