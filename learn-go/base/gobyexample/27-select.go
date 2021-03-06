package main

import (
	"fmt"
	"time"
)

// go 的 select 可以等待多个 channel 操作。
// 将 goroutines 和 channel 与 select 相结合是 go 的一个强大功能

/*
select()
	早在Unix时代, select 机制就已经被引入。
	通过调用 select() 函数来监控一系列的文件句柄,一旦其中一个文件句柄发生了IO动作,该 select() 调用就会被返回。
	后来该机制也被用于实现高并发的Socket服务器程序。
	Go语言直接在语言级别支持 select 关键字,用于处理异步IO问题

	select默认是阻塞的，只有当监听的 channel 中有发送或接收可以进行时才会运行，当多个channel都准备好的时候，select是随机的选择一个执行的。

select 用法类似 switch
	select 限制每个 case 语句里必须是一个IO操作
	default就是当监听的channel都没有准备好的时候，默认执行的（select不再阻塞等待channel）。

select{} 可以用于阻塞一个 goroutine

*/

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i:=0; i<2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}