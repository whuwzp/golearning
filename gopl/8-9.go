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
	defer 	n.Done()
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