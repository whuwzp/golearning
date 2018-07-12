---
title: go学习(4)-微信公众号项目改进
date: 2018-02-28 18:35:33
tags: [go, 语言]
categories: 科研
---

> 大神的项目源码地址：github.com\leeeboo\wechat，他的项目主要是微信公众号后端实现。本文是试图进行改进。

[TOC]

<!-- more -->

##  改进思路

当时觉得用两个function不如试试用接口，但其实最后没觉得简单多少。（尴尬）

## 实现代码

### 方法1（可行）

因为go中的数组和c系列一样，必须是同类型的变量（Python无限制），所以我想让workers自动遍历的时候就会很麻烦，因此只能采用这种比较死的方法实现（具体见经验总结）：

```go
package main

import (
   "io"
   "net/http"
   "time"
   "log"
   "github.com/leeeboo/wechat-new/wx"
   "regexp"
)
type Jobber interface { //jod的接口
   job(w http.ResponseWriter, r *http.Request)
}
type Base struct { //这个基本成员一致
   Method   string
   Pattern  string
}
type Get struct {
   Base
}
type Post struct {
   Base
}
type Workers struct {//两个实现接口
   num int
   Get
   Post
}

var GET = Get{Base{"GET", "^/"}}
var POST = Post{Base{"POST", "^/"}}
var workers Workers

func init() {
   workers = Workers{2,GET, POST}//加入两个实现
}

type httpHandler struct {
}
func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

   t := time.Now()
	//复杂就在这里，尝试了很多方法，希望可以利用类似索引的东西去自动遍历、匹配，但是由于get、post分别为不同类型，无法使用数组、slice这样的（它们都是同类型集合），只能使用结构体，但是结构体没法遍历，只能使用`.`符号
   if workers.Get.Method == r.Method{
      if m, _ := regexp.MatchString(workers.Get.Pattern, r.URL.Path); m{
         workers.Get.job(w, r)
         go writeLog(r, t, "unmatch", "")
      }
   } else {
      if workers.Post.Method == r.Method{
         if m, _ := regexp.MatchString(workers.Post.Pattern, r.URL.Path); m{
            workers.Post.job(w, r)
            go writeLog(r, t, "unmatch", "")
         }
      }
   }
   io.WriteString(w, "")
   return
}

//get实现接口，只是换成了接口，内部代码不变
func (g *Get) job(w http.ResponseWriter, r *http.Request)  {
   client, err := wx.NewClient(r, w, token)

   if err != nil {
      log.Println(err)
      w.WriteHeader(403)
      return
   }

   if len(client.Query.Echostr) > 0 {
      w.Write([]byte(client.Query.Echostr))
      return
   }

   w.WriteHeader(403)
   return
}

//POST
func (p *Post) job(w http.ResponseWriter, r *http.Request)  {
   client, err := wx.NewClient(r, w, token)
   if err != nil {
      log.Println(err)
      w.WriteHeader(403)
      return
   }
   client.Run()
   return
}
```

另外关于定义类型系列可以按照如下简化：（这个原本是方法2想尝试的思路，不过也好用于简化）

```go
type Base struct {
   Method   string
   Pattern  string
}
type Get Base
type Post Base

type Workers struct {
   num int
   Get
   Post
}
```

### 方法2（不可行）

就是之前的数组的问题，因为数组是要求同类型的，因此想到将Get和Post结构体设为相同类型结构体的子类型，如下：

```go
type Base struct {
   Method   string
   Pattern  string
}
type Get Base
type Post Base

var get = Get{"GET", "^/"}
var post = Post{"POST", "^/"}


func init() {
   fmt.Print(GET, POST,"\n")
    var m = [...]Base{get, post} //error:cannot use get (type Get) as type Base in array or slice literal
```

还是出现了错误，显示Get仍然还是Get，而不是Base，因此不能算了同类。

### 方法3（应该可行，但是还不成功）

想利用map，m := map[string]struct{}的形式将两个method放在一起，这样就可以利用如下方法访问了。然而，遇到了问题，待解决，还不太熟悉结构体的原理。

```go
for a, b := range m {} //形式遍历访问了
```

### 

## 经验总结

1. 相同底层类型但是不同类型间不可直接比较；

   ```go
   type Get string
   var GET Get = "GET"
   if r.Method == GET {}//报错，不同类型
   ```

2. 结构体的声明，要么字面值声明如下：

   ```go
   type Base struct {
   	Method   string
   	Pattern  string
   }
   type Get struct {
   	Base
   }
   var GET = Get{Base{"GET", "^/"}}
   ```

   在函数外，只可以定义声明放在一起，不可以单独赋值（也就是声明变量时即赋值，不可以声明之后再赋值）；

3. 还有一个和Python不太一样的地方，go中的数组和c系列一样，必须是同类型的变量（Python无限制），所以我想让workers自动遍历的时候就会很麻烦：

   + 必须新建type，因为string等不能用来实现接口；
   + 如果采用type Get string和type POST string这样的话，数组就不能包含这两个，因为不同类型，当然可以用string(GET)、string(POST)来转化，但是这样一来也会无法使用我们为这两个类型添加的方法了（job），那样数组中的仅仅就是字符串了，丢失了原有的特性方法；

4. 代码如下：

   ```go
   type Get struct {
   }
   type Post struct {
   	Pattern  string	//把这行注释掉就ok 
   }

   var m = map[string]struct{}{
   	"GET": Get{}, //ok
   	"POST": Post{}, //error：cannot use a literal (type a) as type struct {} in assignment
   }
   ```

   ​































































