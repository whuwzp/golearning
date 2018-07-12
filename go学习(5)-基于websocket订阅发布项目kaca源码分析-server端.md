---
title: go学习(5)-基于websocket订阅发布项目kaca源码分析-server端
date: 2018-03-03 18:35:33
tags: [go, 语言]
categories: 科研
---

> kaca的项目源码地址：github.com/scottkiss/kaca，他的项目主要是基于**websokcet**实现**消息订阅与发布**系统。本节主要分析server端工作原理。主要学习了goroutine（多线程）、channel、基于select的多路复用、defer函数的使用等 。

[TOC]

<!-- more -->

## 消息订阅发布原理

1. **sub——>系统**（订阅主题）
2. **pub——>系统**（发布消息）
3. **系统——>sub**（推送消息）

suber向系统订阅自己感兴趣的内容（topic），puber向系统发布消息，系统将消息发给订阅了相应主题的suber。可以类比之前的Ryu控制器的工作原理：应用层APP向opf.handler（Ryu控制器内核，会将openflow消息转化为相应事件）订阅各自需要接收的事件（事件类型），ofp.handler就会将底层发来的事件（原始为消息）发给订阅者。

## 源码分析

### 概要

kaca项目就是利用web的C/S模型来实现的消息订阅发布，其中suber和puber为client，系统为server。本节主要分析server端。由上原理可知，server端主要有以下任务：

1. 与net.http的接口`对接`，即利用websocket模型监听消息、处理消息；
2. 具体的处理函数，包含多client下的高效服务（多线程）、sub/pub模型实现（中心调度模块）。

server端源码：

```go
//code1
go disp.run() 	//负责调度
if checkOrigin {
   http.HandleFunc("/ws", serveWsCheckOrigin)	
    //注册serveWsCheckOrigin为pattern的处理函数（实现对接）
} else {
   http.HandleFunc("/ws", serveWs)
}
err := http.ListenAndServe(addr, nil)	//开始监听
```

### 对接

net.http的注册和监听处理整体如下图所示：[PDF清晰版](/doc/Sci/Go/learning/5-1.pdf)

![kaca_server](/img/Sci/Go/learning/5-1.PNG)

其实整体步骤和微信项目差不太多，只是在ListenAndServe阶段的处理函数不一样，如下：

```go
//code2
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux	
	//和之前不同，这里采用的是默认处理函数
	}
	handler.ServeHTTP(rw, req)
}
```

采用的是默认的处理函数，当然这个处理函数就是在code1中注册的serveWsCheckOrigin，具体注册方法类似，也是实现了接口ServeHTTP的方式完成的。这里稍稍提一下以下代码：

```go
//code3
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	mux.Handle(pattern, HandlerFunc(handler)) 	
    //handler即为serveWsCheckOrigin，HandlerFunc()是类型转换操作
    //即将一个函数转化另一个函数类型
    //因此具备了该函数类型的方法，从而实现接口
}

type HandlerFunc func(ResponseWriter, *Request)	
//转换为此类型，因此转换后具有ServeHTTP的方法
// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)	//实际就是执行函数serveWsCheckOrigin
    //ServeHTTP中的handler.ServeHTTP(rw, req)
    //handler即为HandlerFunc(handler)
}
```

+ 函数类型的也可以具有方法，从而实现某接口；
+ 虽然最后调用f(w, r)和调用函数本身没差别，但是实现了统一的接口，便与开发；
+ 可以用类似HandlerFunc(handler)将函数转化为其他类型函数。

### 处理函数

#### 概要

1. 首先为了实现sub/pub模型：
   + 肯定需要建立一个map，也就是每个client感兴趣的topic表，利用sub消息建立；
   + 依据这个map，将client发布的pub消息发送给特定topic的client。
2. 多线程服务：利用goroutine实现多个连接的同时服务。

#### 多线程

1. **dispatch->run**	//每个连接的dispatch线程处理后交由run（中心调度模块）处理；
	. **run->deliver** 	      //run处理后（涉及map、topic等，以确定发给谁、添加至map等操作），如有需要发送处理的则由deliver发送。

```go
//code4
func serveWsCheckOrigin(w http.ResponseWriter, r *http.Request) {
   disp.register <- c	//新连接注册（新client）
   go c.dispatch()	//为每个新的连接启用处理线程
   c.deliver()	//信息发送的线程
}

```

主要由三个函数实现：

1. run

   为了实现多线程的高效服务，大神首先启用了**调度器**的线程，该线程无限循环，负责`中心调度`；

   ```go
   //code5
   go disp.run()
   func (d *dispatcher) run() {
   	for {
   		select {
               //新连接的注册，添加至connections map[*connection]bool表
               case c := <-d.register:	
               //注销，删去
               case c := <-d.unregister:
               //组播消息，m中带有c.id信息，便于分辨
               case m := <-d.broadcast:
               //订阅topic
               case m := <-d.sub:
               //发布消息
               case m := <-d.pub:
   		}
   	}
   }
   ```

2. dispatch：每个连接（也就是client）都会启用一个线程，一个大循环处理pub/sub等消息，经过；

   ```go
   //code6
   func (c *connection) dispatch() {
   	defer func() {
   		disp.unregister <- c
   		c.ws.Close()
   	}()	//使用了defer函数，该函数会在本函数结束时自动执行（完成连接的注销等）
   	for {	//大循环，循环处理消息
   		_, message, err := c.ws.ReadMessage()
   		msg := string(message)
   		if strings.Contains(msg, SUB_PREFIX) {
   			topic := strings.Split(msg, SUB_PREFIX)[1]
   			disp.sub <- strconv.Itoa(int(c.cid)) + SPLIT_LINE + topic//向通道中传入的消息中带入了c.id信息，便于分辨，map表信息
   		} else if strings.Contains(msg, PUB_PREFIX) {
   			topic_msg := strings.Split(msg, PUB_PREFIX)[1]
   			disp.pub <- topic_msg
   		} else {
   			disp.broadcast <- message
   		}
   	}
   }
   ```

3. deliver：run处理完后由deliver负责发送信息：

   ```go
   //code7
   func (c *connection) deliver() {
   	ticker := time.NewTicker(pingPeriod)
   	defer func() {
   		ticker.Stop()
   		c.ws.Close()
   	}()
   	for {
   		select {
   		case message, ok := <-c.send:
   			if !ok {
   				c.sendMsg(websocket.CloseMessage, []byte{})
   				return
   			}
   			if err := c.sendMsg(websocket.TextMessage, message); err != nil {
   				return
   			}
   		case <-ticker.C:
   			if err := c.sendMsg(websocket.PingMessage, []byte{}); err != nil {
   				return
   			}
   		}
   	}
   }
   ```

#### pub/sub实现

其实pub/sub模型主要就是client和topic的映射，也就是**订阅表**，之后依据特定topic消息发送给特定订阅者。

1. 订阅实现：约定消息中**携带前缀**的形式来标记和辨识sub或者pub。

   ```go
   //code8
   //dispatch()
   msg := string(message)
   if strings.Contains(msg, SUB_PREFIX) {
      topic := strings.Split(msg, SUB_PREFIX)[1]	//分离出topic
      disp.sub <- strconv.Itoa(int(c.cid)) + SPLIT_LINE + topic	
       //将连接c与订阅的topic传入channel（用c.id作为c的唯一标识）
   }    

   //run()
   case m := <-d.sub:
   	//subscribe message
   	msp := strings.Split(m, SPLIT_LINE)
   	log.Println("sub->" + m, msp)
   	for c := range d.connections {
   		if msp[0] == strconv.Itoa(int(c.cid)) {	//为c的topic订阅表中增加topic
   			c.topics = append(c.topics, msp[1])
   		}
   	}
   ```

2. 发布实现：

   ```go
   //code9
   //disptch()
   if strings.Contains(msg, PUB_PREFIX) {
   	topic_msg := strings.Split(msg, PUB_PREFIX)[1]
   	disp.pub <- topic_msg
   }

   //run()
   case m := <-d.pub:
      //publish message
      msp := strings.Split(m, SPLIT_LINE)
      log.Println("pub->" + m)
      for c := range d.connections {	//遍历所有连接
         for _, t := range c.topics {
            if t == msp[0] {	//如果连接（client）订阅了该topic
               select {
               	case c.send <- []byte(msp[1]):	//则扔给deliver的channel去发送
               default:
                  close(c.send)
                  delete(d.connections, c)
               }
               break
            }
         }
      }
   }
   ```


