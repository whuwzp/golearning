---
title: go学习(7)-grpool多线程项目源码分析
date: 2018-03-12 18:35:33
tags: [go, 语言]
categories: 科研
---

> 大神源码：github.com/ivpusic/grpool。他的代码实现的是任务多线程化。主要学习了goroutine和channel的使用。

[TOC]

<!-- more -->

## 概要

总的来说，整体系统由dispatch和worker（开始大循环工作）组成。比较简单，不一一分析了。

1. dispatch：相当于路由，接收JobQ中的新任务，然后找到一个worker，扔进他的channel中，循环；
2. worker：从各自channel中接收分配的任务，完成，**又重新加入到workerpool中**，循环。

## 源码分析

1. channel的使用

   ```go
   type worker struct {
      workerPool chan *worker	//可以看出使用指针作为channel的传输类型
      jobChannel chan Job		//可以看出使用函数作为channel的传输类型
      stop       chan struct{}
   }

   type Job func()

   ```

2. 任务完成后将worker又加入到原pool中，这点很重要

   ```go
   func (w *worker) start() {
      go func() {
         var job Job
         for {
            // worker free, add it to pool
            w.workerPool <- w

            select {
            case job = <-w.jobChannel:
               job()
            case <-w.stop:
               w.stop <- struct{}{}
               return
            }
         }
      }()
   }
   ```





























