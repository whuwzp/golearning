---
title: go学习(6)-基于websocket订阅发布项目kaca源码分析-client端
date: 2018-03-06 18:35:33
tags: [go, 语言]
categories: 科研
---


> kaca的项目源码地址：github.com/scottkiss/kaca，他的项目主要是基于**websokcet**实现**消息订阅与发布**系统。本节主要分析client端工作原理。主要学习了匿名函数的使用等 。

[TOC]

<!-- more -->

上节分析了server端，client端相对容易些（kaca本身就是实现服务架构）。

## 源码分析

### 概要

client端的工作比较简单，分为以下：

1. 创建连接，同server端进行连接，类似server端tcp绑定，client端连接；
2. 订阅感兴趣的topic（上节分析可以看出订阅原理）、发布消息等。

### 连接

NewClient->websocket.DefaultDialer.Dial函数建立连接（具体还待进一步分析）。返回&client。

### 订阅发布

```go
//code1
func (c *client) Sub(topic string) {
   sendMsg := SUB_PREFIX + topic	//上节说的标记方法，便于server端识别
   err := c.conn.WriteMessage(websocket.TextMessage, []byte(sendMsg))	//相当于发送了
}	//发布同理
```

## 经验总结

client相对简单，唯一注意到以下代码，即匿名函数的使用：

```go
//code2
consumer.ConsumeMessage(func(message string) {	//将匿名函数作为传入值
   fmt.Println("consume =>" + message)
})
```

可以看出，将匿名函数作为传入值：

```go
//code3
func (c *client) ConsumeMessage(f func(m string)) {	//定义了一个匿名函数，作为传入值
   go func() {	//启用一个线程，大循环接收消息
      for {
         _, message, err := c.conn.ReadMessage()
         if err != nil {
            log.Println("read:", err)
            break
         }
         log.Printf("recv: %s", message)
         f(string(message))	//调用f打印消息
      }
   }()	//使用了匿名函数，最后的括号内为匿名函数的传入值
}
```

可以看到使用了两个匿名函数：一个是传入值，一个是函数内的匿名函数。

好像显得很麻烦，因为完全可以不要那个传入值，直接import "fmt"，进行打印。可能是为了避免多import，或者就是为了炫技呢，哈哈，厉害了大神！
