### 介绍
reduceByKey,相较于普通的 shuffle 操作,它的一个特点就是会进行 map 端的本地聚合
对 map 端给下个 stage 每个 task 创建的输出文件中,写数据之前,就会进行本地的 combiner 操作,也就是说对每一个 key,对应的values,都会执行算子
函数
### 对性能的提升
1. 在进行本地聚合以后,在 map 端的数据量就变少了,减少磁盘 IO. 而且可以减少磁盘空间的占用
2. 下一个 stage ,拉取数据的量,也就变少了,减少网络的数据传输的性能消耗
3. 在 reduce 端进行数据缓存的内存占用就变少了
4. reduce 端,要进行聚合的数据量也变少了

### 使用场景
1. 对于非常普通的,比如说,就是实现类似于 wordcount 程序一样,对每个 key 对应的值,进行某种数据公式或者算法的计算(累加,累乘)
2. 对于一些类似于要对每个 key 进行一些字符串拼接的这种较为复杂的操作,可以自己衡量一下,最好用 reduceByKey 实现.
(shuffle 基本上占了整个 spark 作业 90% 的性能消耗,只要能对 shuffle 进行一定的调优,都是有价值的)
