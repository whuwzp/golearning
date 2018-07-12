---
title: go学习(2)-《GO语言圣经》读书笔记以及练习题
date: 2018-02-24 18:35:33
tags: [go, 语言]
categories: 科研

---

> 《GO语言圣经》（gitbook上的开源书籍）读书笔记以及练习题。

<!-- more -->



## 程序结构

### 声明和赋值

1. “:=”是一个变量声明语句，而“=‘是一个变量赋值操作

   ```go
   i, j = j, i // 交换 i 和 j 的值
   i, j := 0, 1 //声明并赋值
   ```

2. 简短变量声明语句中必须至少要声明一个新的变量，下面的代码将不能编译通过：

   ```go
   f, err := os.Open(infile)
   f, err := os.Create(outfile) // compile error: no new variables 
   _, err ：= io.Copy(os.Stdout, resp.Body) // 也错误，因为只有err一个变量，且前面已经声明过了
   ```

   解决的方法是第二个简短变量声明语句改用普通的多重赋值语言。 

### 命令行标志参数

使用示例如下：

```go
$ flag -s / a ba c
//其中flag为运行程序（可执行文件）
//s为标志参数（前面的"-"应该相当于标识?）
//再后面的"a ba c"就是flag.Args中的参数了
```

其他系统或命令行中常见的如： `-h`或者`-help`,是一个道理。这里的`-s /`表示用"/"分隔各个参数，`-n `用于忽略行尾的换行符。

 ```go
$ flag -h //命令行输入
Usage of D:\code\go\bin\flag.exe:
  -n    using enter or not
  -s string
        septer (default " ")
 ```

书中实现代码如下：

```go
package main
import (
	"flag"
	"fmt"
	"strings"
)
var n = flag.Bool("n", false, "using enter or not") 
// 第一个n是指向标志参数n（第二个）的指针，flag.xxx()返回的是指针
//"n"为参数的名字，false为该参数的默认值，"using enter or not"为说明介绍
var s = flag.String("s", " ", "septer")
func main() {
	flag.Parse() //必须使用这个使其生效
	fmt.Print(strings.Join(flag.Args(), *s)) // *s，取该地址上参数的值，默认为" "（空格）
	if !*n { //*n，同理，默认为false
		fmt.Println()
	}
}
```



### 自定义数据类型

```go
type Celsius float64 // 摄氏温度
type Fahrenheit float64 // 华氏温度
```

这样Celsius，Fahrenheit就和int、float等具有相同的地位（当然底层还有的自然会不一样）。

只有相同类型的才具有"=="比较等特性。这样可以避免出现某些单位不一致等错误。

### 包和文件 

1. 同一个包（如tempconv）的多个go文件（假设有tempconc.go，conv.go两个源文件），则两个源文件都为`package tempconv` ，文件夹名为tempconv；
2. 在一个文件声明的类型和常量，在**同一个包**的其他源文件也是可以直接访问的；
3. 当包被导入的时候，**包内**的成员将通过类似tempconv.CToF的形式访问 ;
4. 由于CTOF为**大写**，因此**包外**引用后也可以使用（所有大写的变量和函数都可以供外部使用）；
5. 按照惯例，一个包的名字和包的导入路径的最后一个字段相同，例如gopl.io/ch2/tempconv包的名字一般是tempconv 

### 作用域

1. 语法块是由花括弧所包含的一系列语句，就像函数体或循环体花括弧对应的语法块那样。语法块内部声明的名字是无法被外部语法块访问的 ；
2. 任何在在函数外部（也就是包级语法域） 声明的名字可以在同一个包的任何源文件中访问的 ；
3. 对于导入的包，例如tempconv导入的fmt包，则是对应源文件级的作用域，因此**只能在当前**的文件中访问导入的fmt包，当前包的其它源文件无法访问在当前源文件导入的包 ;

## 基础数据类型

### 整数

1. **无符号和无符号：**其中有符号整数采用2的补码形式表示，也就是最高bit位用作表示符号位，一个n-bit的有符号数的值域是从-2 到2 - 1。无符号整数的所有bit位都用于表示非负数，值域是0到2 - 1。 

2. int8类型整数的值域是从-128到127，而uint8类型整数的值域是从0到255 .

3. 通常Printf格式化字符串包含多个%参数时将会包含对应相同数量的额外操作数，但是%之后的 [1] 副词告诉Printf函数再次使用第一个操作数。第二，%后的 # 副词告诉Printf在用%o、%x或%X输出时生成0、0x或0X前缀 

   ```go
   o := 0666
   fmt.Printf("%d %[1]o %#[1]o\n", o) // "438 666 0666"
   ```

4. 早期计算机世界只有一个ASCII字符集，美国信息交换标准代码。ASCII，更准确地说是美国的ASCII，使用7bit来表示128个字符；但是其他语言字符无法被表示。Unicode则是包含了所有的世界字符120000种，需要32bit，但是占用资源较多；utf-8是Unicode的一个标准，使用的是变长编码，提高了利用率。

## 复合型

### 字符串

字符串是不可修改的

### 数组

一般形式：

```go
var q [3]int = [3]int{1, 2, 3}
```

但是也可以指定一个**索引**和对应值列表的方式初始化 ：（实际USD等就是int类型的具体值）

```go
type Currency int
const (
    USD Currency = iota // 美元
    EUR // 欧元
    GBP // 英镑
    RMB // 人民币
)
symbol := [...]string{USD: "$", EUR: "€", GBP: "￡", RMB: "￥"}
fmt.Println(RMB, symbol[RMB]) 
```

1. 不同于很多语言，go中数组是不会作为隐式的指针传递的，也就是传入的是数组的复制品，并不会改变数组本身，这种复制本身效率低，可以用指针传入；

### slice（动态数组）

一般形式：

```go
months := [...]string{1: "January", /* ... */, 12: "December"} 
```

1. 一个slice由三个部分构成：指针、长度和容量。 内置的len和cap函数分别返回slice的长度和容量；
2. 指针指向第一个slice元素对应的底层数组元素的地址；
3. 内置的append函数用于向slice追加元素；
4. 和数组不同的是，slice之间不能比较 ,slice唯一合法的比较操作是和nil比较 

### map（词典）

1.  map中所有的key都有相同的类型，所有的value也有着相同的类型 

2. 如果map元素类型是一个数字，你可以需要区分一个已经存在的0，和不存在而返回零值的0，可以像下面这样测试：（一个不存在的key，如a[k]将会返回0，所以要检验）

   ```go
   age, ok := ages["bob"]if !ok { /* "bob" is not a key in this map; age == 0. */ } 
   ```

3. map的key一定是要可以比较的，但是如果是slice则无法比较，可以采用下面方法将slice转化为string进行比较：

   ```go
   func k(list []string) string { return fmt.Sprintf("%q", list) } 
   ```

4. map的key需要是可以比较的，map本身不可比较，因此直接不行；结构体可以比较，所以可以作为key；

### 结构体

1. 结构体的指针和结构体相同“功用”，不知道怎么表述，详见下例：

   ```go
   package main

   import (
   	"fmt"
   )
   func main ()  {
   	type student struct {
   		name string
   	}
   	var a student
   	a.name = "wang"
   	
   	var b *student = &a
   	var c *student =&a
   	(*b).name = "luo"
   	fmt.Printf("%T\t%T\n", a, b)
   	fmt.Println(a.name,b.name);

   	c.name = "wen"
   	fmt.Printf("%T\t%T\n", a, c)
   	fmt.Println(a.name,c.name);
   }

   //output
   main.student	*main.student
   luo luo
   main.student	*main.student
   wen wen
   ```

2. 如果考虑效率的话，较大的结构体通常会用指针的方式传入和返回,返回值一般为*student这样的；

3. 结构体中成员小写则不可导出：

   ```go
   package p
   type T struct{ a, b int } // a and b are not exported
   package q
   import "p"
   var _ = p.T{a: 1, b: 2} // compile error: can't reference a, b
   var _ = p.T{1, 2} // compile error: can't reference a, b
   ```

   #### 嵌入结构体

   一个结构体包含其他结构体，避免重复写相同的成员

   ```go
   type Point struct {
       X, Y int
   } 
   type Circle struct {
       Center Point
       Radius int
   }
   type Wheel struct {
       Circle Circle
       Spokes int
   }

   //访问方式：
   var w Wheel
   w.Circle.Center.X = 8
   w.Circle.Center.Y = 8
   w.Circle.Radius = 5
   w.Spokes = 20
   ```

   赋值太繁琐，可以采用**匿名成员**的方法：

   ```go
   type Circle struct {
       Point
       Radius int
   } 
   type Wheel struct {
       Circle
       Spokes int
   }

   //访问方式：
   var w Wheel
   w.X = 8 // equivalent to w.Circle.Point.X = 8
   w.Y = 8 // equivalent to w.Circle.Point.Y = 8
   w.Radius = 5 // equivalent to w.Circle.Radius = 5
   w.Spokes = 20
   ```

   但是字面值却无法简洁：

   ```go
   //以下是错误的：
   w = Wheel{8, 8, 5, 20} // compile error: unknown fields
   w = Wheel{X: 8, Y: 8, Radius: 5, Spokes: 20} // compile error: unknown fields

   //以下是正确的方式
   w = Wheel{Circle{Point{8, 8}, 5}, 20}
   w = Wheel{
   	Circle: Circle{
       Point: Point{X: 8, Y: 8},
       Radius: 5,
   },
   Spokes: 20, // NOTE: trailing comma necessary here (and at Radius)
   }
   ```

   匿名成员的应用在cache2go项目可以看到。

## 函数

函数像其他值一样，拥有类型，可以被赋值给其他变量，传递给函数，从函数返回。对函数值（function value） 的调用类似函数调用。例子如下：

```go
func square(n int) int { return n * n }
func negative(n int) int { return -n }
func product(m, n int) int { return m * n }
f := square
fmt.Println(f(3)) // "9"
f = negative  //用函数进行赋值
fmt.Println(f(3)) // "-3"
fmt.Printf("%T\n", f) // "func(int) int"
f = product // compile error: can't assign func(int, int) int to func(int) int 
```

### defer

当defer语句被执行时，跟在defer后面的函数会被延迟执行。直到包含该defer语句的函数执行完毕时，defer后的函数才会被执行，不论包含defer语句的函数是通过return正常结束，还是由于panic导致的异常结束。 

### 匿名函数

函数squares返回另一个类型为 func() int 的函数。对squares的一次调用会生成一个局部变量x并返回一个匿名函数。每次调用时匿名函数时，该函数都会先使x的值加1，再返回x的平方。第二次调用squares时，会生成第二个x变量，并返回一个新的匿名函数。新匿名函数操作的是第二个x变量。 

```go
func squares() func() int {
    var x int
    return func() int {
        x++
        return x * x
	}
}
func main() {
    f := squares()
    fmt.Println(f()) // "1"
    fmt.Println(f()) // "4"
    fmt.Println(f()) // "9"
    fmt.Println(f()) // "16"
}
```

我们看到变量的生命周期不由它的作用域决定：squares返回后，**变量x仍然隐式的存在于f中**。 

```go
func (i int) int {
    fmt.print(i)
}(2) //2是匿名函数的传入值，没有则空着
```



## 方法

### 方法和函数的对比

```go
package geometry
import "math"
type Point struct{ X, Y float64 }
// traditional function
func Distance(p, q Point) float64 {
return math.Hypot(q.X-p.X, q.Y-p.Y)
} 
// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
return math.Hypot(q.X-p.X, q.Y-p.Y)
}
```

这里我们已经看到了方法比之函数的一些好处：方法名可以简短。当我们在包外调用的时候这种好处就会被放大，因为我们可以使用这个短名字，而可以**省略掉包的名字**，下面是例子：

```go
import "gopl.io/ch6/geometry"
perim := geometry.Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
fmt.Println(geometry.Path.Distance(perim)) // "12", standalone function,采用函数，需要外加包名
fmt.Println(perim.Distance()) // "12", method of geometry.Path ，采用方法，则无需
```

### 基于指针对象的方法

#### 一般声明方法

```go
func (p *Point) ScaleBy(factor float64) {
p.X *= factor
p.Y *= factor
}
```

这个方法的名字是 (*Point).ScaleBy 。

#### 为什么要用指针？

+ 用途1：当这个接受者变量本身比较大时，我们就可以用其指针而不是对象来声明方法 ；
+ 用途2： 见下例：

```go
package main

import (
   "fmt"
)
func main()  {
   var jack student
   jack.name = "zhang"
   jack.changeName1("wang")
   fmt.Printf("the name is %s\n", jack.name)
   jack.changeName2("wang")
   fmt.Printf("the name is %s\n", jack.name)
}
type student struct {
   name string
}
func (li student) changeName1(s string)  {
   li.name = s
}
func (li *student) changeName2(s string)  {
   li.name = s
}
//output
the name is zhang
the name is wang
```

> 注：
>
> + 一般来说，第二个方法的名字是：*student.changeName2，因此需要用一个指针，但是我们并没有用(&jack).changeName2（就是先取下地址），这是因为编译器会**隐式地为我们取变量的地址**。另一方面，，也可以可以用一个 *Point 这样的接收器来调用Point的方法，因为编译器在这里也会给我们**隐式地插入 * 这个操作符** 。
> + 引用作者原文：不管你的method的receiver是指针类型还是非指针类型，都是可以通过指针/非指针类型
>   进行调用的，编译器会帮你做类型转换 
> + 由此分析，我们不能通过一个无法取到地址的接收器来调用指针方法，比如临时变量的内存地址就无法获取得到：`Point{1, 2}.ScaleBy(2) // compile error: can't take address of Point literal `

## 接口

> 接口是约定

在微信项目源码分析和改进中可以熟悉接口的使用。

## 线程

### goroutine

当一个程序启动时，其主函数即在一个单独的goroutine中运行，我们叫它main goroutine。新的goroutine会用go语句来创建。在语法上，go语句是一个普通的函数或方法调用前加上关键字go。go语句会使其语句中的函数在一个新创建的goroutine中运行。而go语句本身会迅速地完成。 

## channel

### 无缓存（同步）channel 

一个基于无缓存Channels的发送操作将导致发送者goroutine阻塞，直到另一个goroutine在相同的Channels上执行接收操作，当发送的值通过Channels成功传输之后，两个goroutine可以继续执行后面的语句。反之，如果接收操作先发生，那么接收者goroutine也将阻塞，直到有另一个goroutine在相同的Channels上执行发送操作。 

关闭channel：

```go
//使用下式关闭channel：
close(naturals)

//用这种方法验证close与否
x, ok := <-naturals
if !ok {
    break // channel was closed and drained
} 
```

### 带缓存channel

一般形式：

```go
ch = make(chan string, 3)
```

向缓存Channel的发送操作就是向内部缓存队列的尾部插入元素，接收操作则是从队列的头部删除元素。如果内部缓存队列是满的，那么发送操作将阻塞直到因另一个goroutine执行接收操作而释放了新的队列空间。相反，如果channel是空的，接收操作将阻塞直到有另一个goroutine执行发送操作而向队列插入元素。 

###  串联的Channels（Pipeline） 

Channels也可以用于将多个goroutine链接在一起，一个Channels的输出作为下一个Channels的输入。这种串联的Channels就是所谓的管道（pipeline） 。 

书里的例子直接搬过来，比较好懂：

```go
func main() {
    naturals := make(chan int)
    squares := make(chan int)
    
    // Counter
    go func() {
        for x := 0; ; x++ {
        	naturals <- x
        }
    }()
    
    // Squarer
    go func() {
        for {
            x := <-naturals
            squares <- x * x
        }
    }()
    
    // Printer (in main goroutine)
    for {
    	fmt.Println(<-squares)
    }
}
```



### 单方向channel

类型 chan<- int 表示一个只发送int的channel，只能发送不能接收。相反，类型 <-chan int 表示一个只接收int的channel，只能接收不能发送。

### 基于select的多路复用

1. select会等待case中有能够执行的case时去执行。当条件满足时，select才会去通信并执行case之后的语句；这时候其它通信是不会执行的；
2. 一个没有任何case的select语句写作select{}，会永远地等待下去；
3. 如果多个case同时就绪时，select会随机地选择一个执行，这样来保证每一个channel都有平等的被select的机会；

## 其他

### 标签break

这里的break语句用到了标签break，这样可以**同时终结select和for两个循环**；如果没有用标签就break的话只会退出内层的select循环，而外层的for循环会使之进入下一轮select循环。 

```go
loop:
    for {
    	select {
    		case size, ok := <-fileSizes:
            if !ok {
           		break loop // fileSizes was closed
    		} 
        }
    }
```



## 练习题

练习题的代码见我的github账户，地址：github.com/whuwzp/goland-learning/gopl

### 第八章：Goroutines和Channels 

#### 练习8.1

我的方法如下：

```go
//server
package main

import (
   "io"
   "log"
   "net"
   "time"
   "flag"
)

func main() {
   var port = flag.String("port", "", "connect port")	//创建flag参数
   flag.Parse()
   listener, err := net.Listen("tcp", "localhost:" + *port)	//根据flag参数确定端口
   if err != nil {
      log.Fatal(err)
   }
   for {
      conn, err := listener.Accept()
      if err != nil {
         log.Print(err) 
         continue
      }
      go handleConn(conn) 
   }
}
func handleConn(c net.Conn) {
   defer c.Close()
   for {
      _, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
      if err != nil {
         return // e.g., client disconnected
      }
      time.Sleep(1 * time.Second)
   }
}
```

```go
//clockwall
package main

import (
   "io"
   "log"
   "net"
   "os"
)

func main() {
   var ch = make(chan bool)
   ports := [...]string{"8010", "8020", "8030"}	//便于添加新的服务器端口
   for _, port := range ports {
      go conn(port)	//采用多线程的方式
   }
   <-ch	//如果不加这句，main运行到这儿就会结束，其他线程将停止，也就没有效果了
}

func mustCopy(dst io.Writer, src io.Reader) {	//io.Writer、Reader是个接口，只要实现了该接口，具有方法就可以作为参数
   if _, err := io.Copy(dst, src); err != nil {
      log.Fatal(err)
   }
}
func conn(port string)  {
   conn, err := net.Dial("tcp", "localhost:"+ port)
   if err != nil {
      log.Fatal(err)
   }
   defer conn.Close()
   mustCopy(os.Stdout, conn)	//conn具有该方法，即实现了io.Reader的接口
}
```

更新：clockwall也可以使用**匿名函数**进行简化：

```go
func main() {
   var ch = make(chan bool)
   ports := [...]string{"8010", "8020", "8030"}	//便于添加新的服务器端口
   for _, port := range ports {
       go func(p string){
           conn, err := net.Dial("tcp", "localhost:"+ p)
           if err != nil {
              log.Fatal(err)
           }
           defer conn.Close()
           mustCopy(os.Stdout, conn)
       }(port)	//传入值
   }
   <-ch	
}
```

#### 练习8.3

试了一下，好像*net.TCPConn并没有CloseRead和CloseWrite 两个方法（也找源码中找了，确实没有）

```go
//非答案，但是可以满足要求，也就是关闭了stdin，仍然会将最后的信息传回打印
go func() {
   io.Copy(os.Stdout, conn)
   log.Println("done")
   done<-true
}()
mustCopy(conn, os.Stdin)
<-done	//交换了位置，不过这样程序无法停止
conn.Close()	//交换了位置
```



#### 练习8.8

其他不变，只需改handleconn如下：

```go
func handleConn(c net.Conn) {
   input := bufio.NewScanner(c)
   timeout := time.After(10 *time.Second)	//初始化一个超时时间
   for {
      select {
      case <-timeout:	
         c.Close()	//超时后关闭
         return
      default:
         if input.Scan() {
            go echo(c, input.Text(), 1*time.Second)
            timeout := time.After(10 *time.Second)	//因为有了新的输入，所以又重新计时
         }
      }
   }
   c.Close()
}
```
#### 练习8.9

这里就是把原先的main函数改为了函数，由main中路径的子路径依次调用

```go
package main

import (
   "flag"
   "fmt"
   "io/ioutil"
   "os"
   "path/filepath"
   "sync"
)

func main() {
   var n_main sync.WaitGroup
   // Determine the initial directories.
   flag.Parse()
   roots := flag.Args()
   if len(roots) == 0 {
      roots = []string{"D:/test/"}
   }
   // Traverse the file tree.
   for _, entry := range dirents(roots[0]){
      n_main.Add(1)
      go du([]string{filepath.Join(roots[0], entry.Name())}, &n_main)
       //因为之前没有加filepath.Join(roots[0], entry.Name())，而是直接用entry.name
       //出现 the system cannnot find the file specified的错误
       //因为只是文件名，没有加进路径
   }
   n_main.Wait()
}


func du(roots []string, n_m *sync.WaitGroup) {
   defer n_m.Done()
   fileSizes := make(chan int64)
   var n sync.WaitGroup

   n.Add(1)
   go walkDir(roots[0], &n, fileSizes)

   var nfiles, nbytes int64
   go func() {
      n.Wait()
      close(fileSizes)
   }()
   for f := range fileSizes{
      nfiles++
      nbytes += f
   }
   printDiskUsage(nfiles, &roots, nbytes) // final totals
}



func printDiskUsage(nfiles int64, r *[]string, nbytes int64) {
   fmt.Printf("%s: %d files %.3f GB\n", (*r)[0], nfiles, float64(nbytes)/1e9)
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
   defer  n.Done()
   for _, entry := range dirents(dir) {
      if entry.IsDir() {
         subdir := filepath.Join(dir, entry.Name())
         n.Add(1)
         go walkDir(subdir, n, fileSizes)
      } else {
         fileSizes <- entry.Size()
      }
   }
}
func dirents(dir string) []os.FileInfo {
   entries, err := ioutil.ReadDir(dir)
   if err != nil {
      fmt.Fprintf(os.Stderr, "du1: %v\n", err)
      return nil
   }
   return entries
}
```