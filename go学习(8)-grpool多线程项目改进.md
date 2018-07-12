---
title: go学习(8)-grpool多线程项目改进
date: 2018-03-12 18:35:33
tags: [go, 语言]
categories: 科研
---

> 大神源码：github.com/ivpusic/grpool。他的代码实现的是任务多线程化。主要学习了goroutine和channel的使用。上节源码进行了分析，本节想试着改进一下。改进代码地址：github.com/whuwzp/goland-learning/grpool-test/grpool.go。

[TOC]

<!-- more -->



## 改进目标

就是觉得workpool，jobqueue，dispatcher等太繁杂了，希望能简化这些。

## 改进方法

### 改机思路1

>  将pool改为全局变量，这样就无需每次都传递这个变量或者是地址。

代码如下：

```go
package grpool

import (

)
var pool Pool	//变为全局变量

type Job func()

//无需加入workerpool的成员变量
//不需要NewWorker中传入pool
//每个worker也不用都携带pool的地址
type Worker struct {	
   JobChannel chan Job
   Stop chan struct{}
}

type Pool struct {
   Workers chan *Worker
   JobQueue chan Job
}

func NewPool (NumWorkers int,  JobQueueLen int) *Pool {
   workers := make(chan *Worker, NumWorkers)
   jobqueue := make(chan Job, JobQueueLen)

   for i := 0; i <= cap(workers); i++{
      worker := NewWorker()
      worker.start()
   }
   pool = Pool{
      Workers: workers,
      JobQueue: jobqueue,
   }
   go dispatch()

   return &pool
}

func dispatch(){
   for {
      select {
      case job := <-pool.JobQueue:
         w := <-pool.Workers
         w.JobChannel<-job
      }

   }
}

func NewWorker() *Worker{
   return &Worker{
      JobChannel: make(chan Job),
   }
}

func (w *Worker) start()  {
   go func() {
      for {
         pool.Workers<- w
         select {
         case job := <-w.JobChannel:
            job()
         }
      }
   }()
}
```



### 改进思路2

> 上面的方法中pool是全局变量，作用域太大，做了如下改动。在函数中加入&pool的传入值或者变为pool的方法。

代码如下：

```go
package grpool

import (

)


type Job func()

type Worker struct {
   JobChannel chan Job
   Stop chan struct{}
}

type Pool struct {
   Workers chan *Worker
   JobQueue chan Job
}

func NewPool (NumWorkers int,  JobQueueLen int) *Pool {
   workers := make(chan *Worker, NumWorkers)
   jobqueue := make(chan Job, JobQueueLen)

   pool := &Pool{
      Workers: workers,
      JobQueue: jobqueue,
   }

   for i := 0; i <= cap(workers); i++{
      worker := NewWorker()
      worker.start(pool)
   }

   go pool.dispatch()

   return pool
}

//这里也是使得函数可以调用pool
func (pool *Pool)dispatch(){
   for {
      select {
      case job := <-pool.JobQueue:
         w := <-pool.Workers
         w.JobChannel<-job
      }

   }
}

func NewWorker() *Worker{
   return &Worker{
      JobChannel: make(chan Job),
   }

}

//不得不将pool作为参数传递进函数中
func (w *Worker) start(pool *Pool)  {
   go func() {
      for {
         pool.Workers<- w
         select {
         case job := <-w.JobChannel:
            job()
         }
      }
   }()
}
```











































