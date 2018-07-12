---
title: go学习(1)-安装配置
date: 2018-02-21 18:35:33
tags: [go, 语言]
categories: 科研
---

>go的安装和环境配置


## 安装

1. 下载（国内不行，需要镜像）
2. 安装
```
tar -C /usr/local -xzf go1.9.tar.gz //解压至/usr/local目录下
```

## 配置
### go配置说明

总共涉及两个环境变量<sup>**见附录1**</sup>：`GOROOT`和`GOPATH`。其中，GOROOT为go的安装目录（即/usr/local/go），GOPATH为工作目录（开发源码，中间文件等，可以添加多个）。

```
sudo gedit /etc/profile

//下面编辑该文件，在末尾添加如下代码（需要去掉注释部分）
//配置GOROOT
export GOROOT=/usr/local/go 
export PATH=$PATH:$GOROOT/bin   
//配置GOPATH
export GOPATH=$HOME/mygo
export PATH=$PATH:$GOPATH/bin

//保存文件，并使其生效
source ~/.profile
//此时，可查看版本等
go version
//不行则重启

```

注释：
1. 配置GOROOT后，即可在系统任何位置执行`go`相关的命令<sup>**见附录1**</sup>，如`go build`、`go install`等。
2. 配置GOPATH是工作区间，我是在用户目录下建了`mygo`文件夹，其中包含了三个子目录：
+ bin目录包含可执行命令，存放`go install`命令生成的可执行文件（由*export PATH=$PATH:$GOPATH/bin*导致）
+ pkg目录包含包对象，编译相关，`go build`
+ src目录包含go的源文件，它们被组织成包（每个目录都对应一个包），使用`go install`时，会自动在该目录下找。

### GoLand配置
原谅我第一次配置这些，如图：（要设置为文件的完整路径）
![](/img/Sci/Go/learning/1-1.png)

## hello world

### 码代码
```
mkdir $GOPATH/src/test/
mkdir $GOPATH/src/test/hello
gedit $GOPATH/src/test/hello/hello.go

//编辑hello.go文件，如下：
package main  
import "fmt"    
func main() {  
    fmt.Printf("Hello, world.\n")  
}  
```

### 编译
```
go install test/hello
//go install 会去 $GOPATH/src中去找，test/hello中的，.go文件
```
这时，在`$GOPATH/bin`中生成可执行文件`hello`。

### 执行

```
hello
//这时，输出：Hello, world.
```
因为`export PATH=$PATH:$GOPATH/bin`，所以在任何输入`hello`，即相当于在那个路径执行<sup>**见附录1**</sup>。


## 附录
### 环境变量

> 添加环境变量：添加系统环境变量是为了更方便的使用命令，这样的话，在系统任何路径下输入某个命令`a`，即相当于执行该路径下的该命令`a`,而不用每次都写绝对路径。


下面以配置GOROOT为例进行详解：
```
export GOROOT=/usr/local/go 
```
GOROOT为go的安装目录，在系统中输入路径时可以简写`/usr/local/go`为`$GOROOT`。类似 的还有，如默认的`$HOME`指代的就是`当前用户的工作目录`。

```
export PATH=$PATH:$GOROOT/bin   
```
这是添加系统环境变量，这样的话，在系统任何路径下输入命令`go`即相当于执行`/$GOROOT/bin/.go`(将`$GOROOT`替换为`/usr/local/go`)，因此完整执行为：`/usr/local/go/bin/.go`。


## go路径

在读docker源码时，发现有些函数标红，无法使用goland跳转，估计就是路径不对，于是查看`import`如下：


```go
import (
   "fmt"
   "os"
   "path/filepath"
   "runtime"

   "github.com/docker/docker/cli"
   "github.com/docker/docker/daemon/config"
   "github.com/docker/docker/dockerversion"
   "github.com/docker/docker/pkg/reexec"
   "github.com/docker/docker/pkg/term"
   "github.com/sirupsen/logrus"
   "github.com/spf13/cobra"
)
```

于是，将docker文件夹，放在了`D:\code\go\src\github.com\docker`目录下，这样就可以了



## sublime text3 go 配置
参考网址：https://studygolang.com/articles/4938
