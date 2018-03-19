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