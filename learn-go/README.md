### 运行时
#### Context 是做什么用的 ？
- 多个 Goroutine 组成的树中同步`取消信号`以减少对资源的消耗和占用，timeout，cancel
- 在一个调用链中传递上下文

#### 同步
- 并发编程同步原语  sync.Mutex、sync.RWMutex、sync.WaitGroup、sync.Once 和 sync.Cond 以及扩展原语 errgroup.Group、semaphore.Weighted 和 singleflight.Group
- sync.Mutex 有两种模式 — 正常模式和饥饿模式 1ms 转变
- 相比于饥饿模式，正常模式下的互斥锁能够提供更好地性能，饥饿模式的能避免 Goroutine 由于陷入等待无法获取锁而造成的高尾延时。
- sync.RWMutex  不限制并发读，性能更好
- sync.WaitGroup 可以等待一组 Goroutine 的返回，一个比较常见的使用场景是批量发出 RPC 或者 HTTP 请求
- sync.Once.Do 方法中传入的函数只会被执行一次，哪怕函数中发生了 panic；
- 两次调用 sync.Once.Do 方法传入不同的函数也只会执行第一次调用的函数；
- sync.Cond.Wait 方法会将当前 Goroutine 陷入休眠状态
- sync.Cond.Signal 方法会唤醒队列最前面的 Goroutine；
- sync.Cond.Broadcast 方法会唤醒队列中全部的 Goroutine；
- sync.Cond 不是一个常用的同步机制，在遇到长时间条件无法满足时，与使用 for {} 进行忙碌等待相比，sync.Cond 能够让出处理器的使用权

#### 定时器

#### channel CSP
- 通信顺序进程（Communicating sequential processes，CSP）1。Goroutine 和 Channel 分别对应 CSP 中的实体和传递信息的媒介，
- channel 无锁管道，其实是乐观锁，依赖 cas 原语
- Channel 是一个用于同步和通信的有锁队列


#### 调度器
Go调度器的发展史，单线程调度器，多线程调度器，任务窃取调度器，抢占式调度器（基于协作的和基于信号的），非均匀内存访问调度器   

- 单线程很单纯
- 多线程有全局锁
- 任务窃取引入调度器P，作为中间层，当前处理器本地队列中的goroutine运行完后，会从其他处理器的队列中随机获取一些goroutine执行；这种方式有两个问题：1.某些goroutine可以长时间运行，造成其他goroutine饥饿 2.垃圾回收需要暂停整个程序（Stop-the-world，STW），最长可能需要几分钟的时间，导致整个程序无法工作.
- 基于协作的抢占式调度，所有 Goroutine 在函数调用时都有机会进入运行时检查是否需要执行抢占，Go 语言运行时会在垃圾回收暂停程序、系统监控发现 Goroutine 运
行超过 10ms 时发出抢占请求，因为这里的抢占是通过编译器插入函数实现的，还是需要函数调用作为入口才能
触发抢占，所以这是一种协作式的抢占式调度。目前的抢占式调度也只会在垃圾回收扫描任务时触发。

- Go 语言在 1.14 版本中实现了基于信号的抢占式调度。

##### 数据结构
- G — 表示 Goroutine，它是一个待执行的任务；
- M — 表示操作系统的线程，它由操作系统的调度器调度和管理；
- P — 表示`处理器`，它可以被看做运行在线程上的本地调度器；调度器在启动时就会创建 GOMAXPROCS 个处理器

说明：  
- Goroutine 只存在于Go运行时，是Go在用户态提供的线程。占用内存小，降低了上下文切换。  
- 调度器最多可以创建 10000 个线程，但是其中大多数的线程都不会执行用户代码（可能陷入系统调用），最多只会有 GOMAXPROCS 个活跃线程能够正常运行。  
- 在默认情况下，运行时会将 GOMAXPROCS 设置成当前机器的核数，我们也可以使用 runtime.GOMAXPROCS 来改变程序中最大的线程数。
- 在这种情况下不会触发操作系统的线程调度和上下文切换，所有的调度都会发生在用户态，由 Go 语言调度器触发，能够减少非常多的额外开销。
- 有两个运行队列：一个是`处理器P`本地队列，另一个是`调度器`持有的全局队列。只有在本地队列没有剩余空间时才会使用全局队列。

#### 网络轮询器
Go语言运行时网络轮询器通过多模块设计，可以在不同的操作系统上实现IO多路复用

- 阻塞IO
- 非阻塞IO：需要循环直到数据准备完成。可以在等待的过程中，执行其他任务，增加CPU利用率
- IO多路复用：select（最多1024个文件描述符）， poll(链表结构，无限个文件描述符)

select 限制：
- 监听能力有限：1024
- 内存拷贝开销大，需要维护较大的数据结构存储 fd， 且该结构需要拷贝的内核
- 时间复杂度O(n)，需要遍历 fd 列表

go 会根据不同平台选择不同的模型。epoll， kqueue，evport 。。。

- epoll 常驻内核，select 需要从用户态复制到内核态
- epoll 返回的是可用的fd子集，select 返回的是全部，需要用户自己去判断是否可用
- fd 较少或 fd 都比较繁忙的时候， select 更占优势

#### golang 运行时系统监控
Go 语言的系统监控也起到了很重要的作用，它在内部启动了一个不会中止的循环，在循环的内部会：
- 检查死锁
- 触发GC
- 触发线程抢占
- 触发网络轮询

### 内存管理

#### 内存分配器（管理的是堆中的对象）
- 数据和变量都会被分配到程序所在的虚拟内存中，
- 内存空间包含两个重要区域 — 栈区（Stack）和堆区（Heap）。
- 栈区：函数调用的参数、返回值以及局部变量大都会被分配到栈上，这部分内存会由`编译器`进行管理；
- 堆区：不同编程语言使用不同的方法管理堆区的内存，C++ 等编程语言会由工程师主动申请和释放内存，Go 以及 Java 等编程语言会由工程师和编译器共同管理，堆中的对象由`内存分配器分配`并由`垃圾收集器回收`。

内存分配的三个组件：
- 用户程序（Mutator）
- 分配器（Allocator）
- 收集器（Collector）
用户程序通过通过分配器从堆中初始化相应的内存区域来完成申请内存。

##### 分配方法
###### 线性分配器（Sequential Allocator，Bump Allocator）
只需要在内存中维护一个指向内存特定位置的指针，当用户程序申请内存时，分配器只需要检查剩余的空闲内存、返回分配的内存区域并修改指针在内存中的位置。

- 高效，执行速度快
- 无法在内存被释放时重用内存

需要合适的垃圾回收算法配合使用：  
- 标记压缩（Mark-Compact）
- 复制回收（Copying GC）
- 分代回收（Generational GC）  
这些算法可以通过拷贝的方式整理存活对象的碎片，将空闲内存定期合并，这样就能利用线性分配器的效率提升内存分配器的性能了。

###### 空闲链表分配器（Free-List Allocator）
它可以重用已经被释放的内存，它在内部会维护一个类似链表的数据结构。当用户程序申请内存时，空闲链表分配器会依次遍历空闲的内存块，找到足够大的内存，
然后申请新的资源并修改链表。选择策略：
- 首次适应（First-Fit）：从头开始遍历，选择第一个大于申请内存的内存块；
- 循环首次适应（Next-Fit）：从上次遍历结束位置开始遍历，选择第一个大小大于申请内存的内存块；
- 最优适应（Best-Fit）： 从链表头遍历整个链表，选择最合适的内存块；
- **隔离适应（Segregated-Fit）(与Go语言类似)**： 将内存分割成多个链表，每个链表中的内存块大小相同，申请内存时先找到满足条件的链表，再从链表中选择合适的内存块；

#### 内存分配
- 线程缓存分配（Thread-Caching Malloc，TCMalloc）是用于分配内存的的机制，它比 glibc 中的 malloc 函数还要快很多。
- Go 语言的内存分配器就借鉴了 TCMalloc 的设计实现高速的内存分配，它的`核心理念是使用多级缓存根据将对象根据大小分类，并按照类别实施不同的分配策略`。

大小对象分类：
- 微对象：(0, 16b)
- 小对象：(16b, 32kb)
- 大对象：(32kb,+OO)

内存多级缓存：
- 线程缓存（Thread Cache）
- 中心缓存（Central Cache）
- 页堆（Page Heap）  
分别对应上述不同大小的对象，发现资源不足时，就从上一级组件中获取更多的内存资源 

#### 虚拟内存布局（堆区）
- Go 1.10 以前，堆区的内存空间都是连续的
- 1.11 版本开始，Go 使用稀疏的堆内存空间替代了连续的内存

##### 线性内存
- 实现简单
- 需要预留大块内存
- 不预留会在一些情况下程序崩溃

##### 稀疏内存


### 垃圾回收
#### 常用方法
##### 标记清除
- 标记阶段：从根对象出发查找并标记堆中所有存活的对象；
- 清除阶段：遍历堆中的全部对象，回收未被标记的垃圾对象并将回收的内存加入空闲链表；  

存在 STW(stop the world) 的问题

##### 三色抽象

##### 屏障技术

##### 增量和并发

#### 演进过程
- SWT 标记清除
- 并发垃圾收集
- 混合写屏障

