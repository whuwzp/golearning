---
title: go学习(9)-cache2go内存管理项目源码分析
date: 2018-03-19 18:35:33
tags: [go, 语言]
categories: 科研
---

> 大神源码：github.com/muesli/cache2go。他的代码实现的是实现内存管理。主要学习了time、空接口等的使用。

[TOC]

<!-- more -->



## 源码分析

### 概要

缓存管理cache主要就是时间的问题，即到期删除，因此我主要是去搞清楚其监测原理，其核心方法expirationCheck 。

另外，由于涉及访问问题，因此**访问锁Lock**也需要注意需要使用。
[清晰版PDF](/doc/Sci/Go/learning/9-1.pdf)

![](/img/Sci/Go/learning/9-1.PNG)

### 核心

核心就是cache的到期管理问题，该项目在为cachetable增加新的item条目时会进行监测。

```go
func (table *CacheTable) addInternal(item *CacheItem) {
 	...
   if item.lifeSpan > 0 && (expDur == 0 || item.lifeSpan < expDur) {
      table.expirationCheck()
   }
}
```

```go
func (table *CacheTable) expirationCheck() {
   table.Lock()
   if table.cleanupTimer != nil {
      table.cleanupTimer.Stop()
   }
   if table.cleanupInterval > 0 {
      table.log("Expiration check triggered after", table.cleanupInterval, "for table", table.name)
   } else {
      table.log("Expiration check installed for table", table.name)
   }
   now := time.Now()
   smallestDuration := 0 * time.Second
   for key, item := range table.items {
      item.RLock()
      lifeSpan := item.lifeSpan
      accessedOn := item.accessedOn
      item.RUnlock()

      if lifeSpan == 0 {
         continue
      }
       //距离上一次访问空闲的时间是否超过失效时间，如果是删除该项条目
      if now.Sub(accessedOn) >= lifeSpan {	
         // Item has excessed its lifespan.
         table.deleteInternal(key)
      } else {
          //这里类似冒泡算法，找出最近就要到期的item的剩余时间
         if smallestDuration == 0 || lifeSpan-now.Sub(accessedOn) < smallestDuration {
            smallestDuration = lifeSpan - now.Sub(accessedOn)
         }
      }
   }
	//得到了最近失效时间，在这个时间之后将会启动一个goroutine，到点后再监测，预备删除到期失效条目
   table.cleanupInterval = smallestDuration
   if smallestDuration > 0 {
      table.cleanupTimer = time.AfterFunc(smallestDuration, func() {
         go table.expirationCheck()	//afterfunc配合goroutine用得很妙
      })
   }
   table.Unlock()
}
```



## 复写

为了加强记忆，自己把代码的核心功能写了一下：

### 简易cache2go

只含有核心功能，代码如下：

```go
package my_cache2go

import (
   "time"
   "sync"
   "fmt"
)

type item struct {
   sync.RWMutex
   key interface{}
   data interface{}
   lifespan time.Duration
   accesson time.Time
}

type cachetable struct {
   sync.RWMutex
   name string
   items map[interface{}]*item
}

var (
   cache = make(map[string]*cachetable)
   mutex  sync.RWMutex
)

func Cache(name string) *cachetable  {
   mutex.Lock()
   defer mutex.Unlock()
   t, ok := cache[name]
   if !ok{
      t = &cachetable{
         name: name,
         items: make(map[interface{}]*item),
      }
      cache[name] = t
   }
   return t
}

func (t *cachetable) Add(key interface{}, lifespan time.Duration, data interface{})  {
   item := &item{
      key: key,
      data: data,
      accesson: time.Now(),
      lifespan: lifespan,
   }
   t.Lock()
   t.Addinternal(item)
}

func (t *cachetable) Addinternal(item *item)  {

   _, ok := t.items[item.key]
   if !ok{
      t.items[item.key] = item
      t.Unlock()
      t.expireCheck()
   }
}

func (t *cachetable) expireCheck()  {
   smallest := 0 * time.Second
   for k, i := range t.items{
      Now := time.Now()
      if Now.Sub(i.accesson) >= i.lifespan{
         fmt.Println("over time, stating to delete")
         t.delete(k)
      }
      if smallest == 0 || i.lifespan - Now.Sub(i.accesson) < smallest{
         smallest = i.lifespan - Now.Sub(i.accesson)
      }
   }

   if smallest > 0{
      time.AfterFunc(smallest, func() {
         go t.expireCheck()
      })
   }
}

func (t *cachetable) delete(key interface{})  {
   t.Lock()
   defer t.Unlock()
   delete(t.items, key)
   fmt.Println("deleting...")
}

func (t *cachetable) Value(key interface{}) interface{} {
   t.Lock()
   defer t.Unlock()
   r, ok := t.items[key]
   if !ok {
      fmt.Println("not in cache!")
      return nil
   }
   return r.data
}
```

main函数

```go
package main

import (
   "github.com/whuwzp/my_cache2go"
   "time"
   "fmt"
)

func main() {
   table := my_cache2go.Cache("test_table")
   table.Add("test_item", 5 * time.Second, "just for test")
   v1 := table.Value("test_item")
   fmt.Println("the value is ", v1)
   time.Sleep(5500 * time.Millisecond)
   v2 := table.Value("test_item")
   fmt.Println("the value is ", v2)
}
```

### 经验总结

### time函数使用

#### time.sleep

等待并阻塞一定时间

#### time.after

1. 将返回一个channel，channel类型为time；

2. 执行后，将在一定时间后向channel发送**当前发送**的时间；

3. <-channel接收的time的时间是发送的时间，而不是接收时间。

   ```go
   func main()  {
      t1 := time.Now()
      tc1 := time.After(1 * time.Second)
      fmt.Println("wait1...")
      time.Sleep(2 * time.Second)
      t2 := <-tc1	//t2并不是3.几秒，而是发送时的1.几秒
      t3 := time.Now()
      fmt.Println("wait2...")
      fmt.Println((t2).Sub(t1))
      fmt.Println((t3).Sub(t1))
   }

   //执行结果
   wait1...
   wait2...
   1.0000572s
   2.0011145s
   ```

#### time.tick

和time.after差不多，只是这个是**周期性**的向channel发送时间，after是一次性的。

#### time.afterfunc

如下例：函数返回的t是*time.Timer类型，在t1时间后，执行传入的函数func()。

```go
t := time.AfterFunc(t1, func() {
   go table.expirationCheck()
})
```

如果后悔，可以在执行函数func()前，使用stop结束而不执行。

#### time.after/before

判断时间点前后，返回true/false

```go
func main()  {
   t1 := time.Now()
   time.Sleep(1 * time.Second)
   t2 := time.Now()
   a := t1.Before(t2)
   b := t1.After(t2)
   fmt.Println(a, b)
}
//执行结束
true false
```

#### time.sub

计算时间差

```go
func main()  {
   t1 := time.Now()
   time.Sleep(1 * time.Second)
   t2 := time.Now()
   fmt.Println(t2.Sub(t1))
}
//执行结果
1.0000572s
```

### 结构体匿名函数使用

> 结构体嵌套：一个结构体包含其他结构体，避免重复写相同的成员

```go
type Point struct {
    X, Y int
} 
type Circle struct {
    Center Point
    Radius int
}
//访问方式：
var c Circle
c.Center.X = 8
c.Center.Y = 8
c.Radius = 5
```

访问太繁琐，可以采用**匿名成员**的方法：

```go
type Circle struct {
    Point
    Radius int
} 
//访问方式：
var w Wheel
c.X = 8 // equivalent to c.Center.X = 8
c.Y = 8 // equivalent to c.Center.Y = 8
c.Radius = 5 // equivalent to c.Radius = 5
```

在cache2go的项目中，很多结构体都嵌入了匿名成员RWMutex结构体，因此可以直接访问：

```go
type CacheTable struct {
   sync.RWMutex
}
type CacheItem struct {
	sync.RWMutex
}
```

### 空接口

空接口的使用也是亮点，如下，可以使用**任何类型的变量作为key**，从而算了间接地实现类似Python字典不同类型搁一起的感觉。（之前的WeChat改进可以尝试）

```go
items map[interface{}]*CacheItem
```

### Lock

一个函数a获取了然后调用了其他的函数b，a中使用了Lock，则b中不需要再Lock了，否则阻塞。

```go
t.Lock()
```

































