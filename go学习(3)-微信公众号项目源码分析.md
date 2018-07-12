---
title: go学习(3)-微信公众号项目源码分析
date: 2018-02-24 18:35:33
tags: [go, 语言]
categories: 科研
---



> 大神的项目源码地址：github.com\leeeboo\wechat，他的项目主要是微信公众号后端实现。主要学习了接口的使用、net/http的相关原理、调用关系。

[TOC]

<!-- more -->

## 微信公众号后台原理

**关注用户<——>微信平台<——>公众号后台服务器**

1. 在微信注册时会让设置URL（服务器IP、域名等）和**token**（用于微信和后台间通信的验证）;
2. 用户发消息，首先到达微信平台，由微信转发（加token等验证）给公众号拥有者自己的后台服务器（回复为逆过程）；
3. 该项目就是在后台服务器部分。

## 源码分析

### 概要

由于go语言本身提供了良好的net.http接口，所以开发相对较为容易。大体上，http提供了`ServeHTTP接口`，我们只需要设置处理函数handler对接，并且`实现该接口`即可。因此重点有两个：`对接`和`实现`。

### 对接

go提供了接口，所以第一步我们需要找到`对接方法`，然后再去实现。

如图所示，是main对接http接口的实现概要。[PDF文件清晰版](/doc/Sci/Go/learning/3-1.pdf)

![](/img/Sci/Go/learning/3-1.PNG)

#### main

主要代码如下code1：（不重要略）

```go
//code1
server := http.Server{
   Addr:           fmt.Sprintf(":%d", port),
   Handler:        &httpHandler{},
   ...
}
log.Fatal(server.ListenAndServe())
```

可以查看http.Server，实际对接已经完成了，也就是`&httpHandler{}`（如下code2可以看出结构体httpHandler的指针具有ServeHTTP方法，所以需要取地址符`&`）。因而`*httpHandler`实现了http.Handler的ServeHTTP接口。

```go
//code2
type httpHandler struct {}
func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {...}
```

一般来说，我们只需要对接上，其余的只需要ListenAndServe监听，然后开发处理函数就好，没必要继续。但是这里我稍微看了一下http中的包。

#### ListenAndServe

简要代码如下code3。可以看出\*(http.Server)具有ListenAndServe的方法（监听addr），最终调用(\*(http.Server)).Serve方法。

```go
//code3
func (srv *Server) ListenAndServe() error {
    addr := srv.Addr
    ...
    ln, err := net.Listen("tcp", addr)
    ...
    return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```

#### Serve

这里Serve进入了大循环，接收新连接。新建了连接c，并调用了c.serve。http.conn的解释见下。

```go
//code4
func (srv *Server) Serve(l net.Listener) error {
    ...
    for {
        ...
        rw, e := l.Accept()
        ...
        c := srv.newConn(rw)
        ...
        go c.serve(ctx)
    }
};
```

#### coon

```go
//code5
type conn struct {
    server *Server
}
//(*conn)的方法serve 
func (c *conn) serve(ctx context.Context) {
    ...
    serverHandler{c.server}.ServeHTTP(w, w.req) //c结构体具有成员server，其类型为(*Server)，而code6可以看出，serverHandler结构体具有(*Server)类型成员。
    //感觉有点绕，按理说如果conn实现了ServeHTTP接口，那就可以直接c.server.ServeHTTP(w, w.req)了。当然go自有它的道理吧，不深究。
    ...
}
```

#### serverHandler

```go
//code6
type serverHandler struct {
	srv *Server
}
//serverHandler的方法ServeHTTP
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler //sh.srv即为Server类型，其具有Handler的成员，由main函数可知，我们已经预设了Handler成员为&httpHandler{}。所以最终到达了我们需要自己去实现的方法这里。
	...
	handler.ServeHTTP(rw, req)//调用实现了ServeHTTP接口的方法
}
```

#### httpHandler

```go
//code7
type httpHandler struct {
}
//实现Handler的ServeHTTP接口
func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ...//自己开发部分
}
```

### 实现

上节找到了需要我们自己实现的部分（完成了对接）。

用户请求方法主要为GET和POST两种方法，因此大神写了两个方法，具体调用哪一个则交由(*httpHandler).ServeHTTP进行判断选择。（所以文件名为route也是有原因的，负责将请求`路由`到具体处理函数上，本身不做具体处理）。而路由依据就是r *http.Request，也就是输入。

#### 处理函数组

首先说明POST和GET处理函数。利用init把这两种处理函数加入数组。

```go
//code8
type WebController struct {
	Function func(http.ResponseWriter, *http.Request)
	Method   string
	Pattern  string
}
var mux []WebController
func init() {
   mux = append(mux, WebController{post, "POST", "^/"})
   mux = append(mux, WebController{get, "GET", "^/"})
}
```

#### 路由函数ServeHTTP

其实就是看r.Method匹配是POST还是GET。然后调用。

```go
//code9
func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t := time.Now()
   for _, webController := range mux {
      if m, _ := regexp.MatchString(webController.Pattern, r.URL.Path); m {
         if r.Method == webController.Method {
            webController.Function(w, r)
            go writeLog(r, t, "match", webController.Pattern)
            return
         }
      }
   }
   ...
}
```

#### 处理函数GET

一般先GET，后POST。GET就是请求，并且参数直接放在URL中，POST复杂些，相对来说参数安全些。

简要代码如下：主要是验证请求和返回结果。

```go
//code10
func get(w http.ResponseWriter, r *http.Request) {
    client, err := wx.NewClient(r, w, token)
    ...
    w.Write([]byte(client.Query.Echostr))//输出部分
    ...
}
```

##### NewClient

这里代码就不贴了，主要完成的是提取请求的相关信息并验证签名（也就是自己在微信公众平台上设置的），微信与后台服务器通信的验证token（当然还有其他信息合成的签名）。

#### 处理函数POST

前面的处理同GET，多了一个`client.run()`.

```go
//code11
func (this *WeixinClient) Run() {
	err := this.initMessage()
	MsgType, ok := this.Message["MsgType"].(string)
	switch MsgType {
	case "text":
		this.text()
	...
	}
	return
}
//initMessage
//在initMessage中读取了请求，并且赋值（因为了指针参数，可以直接修改）。并且还有之前的map中检查是否存在key的验证应用。
func (this *WeixinClient) initMessage() error {
	body, err := ioutil.ReadAll(this.Request.Body)
    m, err := mxj.NewMapXml(body)
	message, ok := m["xml"].(map[string]interface{})
	if !ok {
		return errors.New("Invalid Field `xml` Type.")
	}
}
//text
func (this *WeixinClient) text() {
	inMsg, ok := this.Message["Content"].(string)
	var reply TextMessage
	reply.Content = value2CDATA(fmt.Sprintf("我收到的是：%s", inMsg))
	replyXml, err := xml.Marshal(reply)//常用语xml
	this.ResponseWriter.Header().Set("Content-Type", "text/xml")//设置header参数
	this.ResponseWriter.Write(replyXml)//最后写入reply
}
```

关于文本、格式相关的没有深究。

## 测试调试

真正的调试方法将程序应该是部署在自己的服务器上，但这样调试成本太高（而且我只是代码分析，没必要专门部署）。因此有以下简洁的调试方法：

### ngrok

将外网地址映射到本机。申请一个账号，完成外网url到本机localhost端口的映射，然后在微信平台设置URL为该URL，即可。

### 模拟微信平台

在了解了通信原理之后，其实调试只需要伪装成微信平台即可，我之前想过，但是POST包构造不会，这里大神提供了测试代码。这个就不多讲了。

源码地址：github.com\leeeboo\wechat-mp-debugger

## 经验总结

1. 熟悉了goland的单步调试，很好分析运行步骤和变量情况；
2. 运用了很多指针方法、接口等概念，学习了相关知识；
3. 自己绘制了流程图，这样确实很好看整个架构，尤其是http名称太多相似的时候，借助流程图很方便。
4. goland的使用：
![](/img/Sci/Go/learning/3-2.PNG)


















































































