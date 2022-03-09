package app

import (
	heap2 "KWayMerger/heap"
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
)

func readSortRewrite(file string) {
	// read
	fd, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}
	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanWords)
	var temp int
	var list []int
	for scanner.Scan() {
		temp, err = strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalln(err)
		}
		list = append(list, temp)
	}
	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	err = fd.Close()
	if err != nil {
		log.Fatalln(err)
	}
	// sort
	sort.Ints(list)
	// rewrite
	fd2, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(list); i++ {
		_, err = fmt.Fprintf(fd2, "%v\n", list[i])
		if err != nil {
			log.Fatalln(err)
		}
	}
	err = fd2.Sync()
	if err != nil {
		log.Fatalln(err)
	}
	err = fd2.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

// Merge all numbers in all files and output to 'filename'
func mergeAndWrite(input []string, output string) {
	fd, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	var minHeap heap2.MinHeap
	heap.Init(&minHeap)
	for i := 0; i < len(input); i++ {
		node, err := heap2.NewNode(input[i])
		if err == nil {
			heap.Push(&minHeap, node)
		}
	}
	for !minHeap.Empty() {
		node := heap.Pop(&minHeap).(heap2.Node)
		// write the val to the output file
		_, err = fmt.Fprintf(fd, "%v\n", node.Val)
		if err != nil {
			log.Fatalln(err)
		}
		if node.Scanner.Scan() {
			val, err := strconv.Atoi(node.Scanner.Text())
			if err != nil {
				log.Fatalln(err)
			}
			node.Val = val
			heap.Push(&minHeap, node)
		} else {
			err := node.Fd.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	err = fd.Sync()
	if err != nil {
		log.Fatalln(err)
	}
	err = fd.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func Run(input []string, output string) {
	var wg sync.WaitGroup
	wg.Add(len(input))
	for i := 0; i < len(input); i++ {
		j := i
		// Sort and rewrite integers in files, one goroutine for one file
		go func() {
			defer wg.Done()
			readSortRewrite(input[j])
		}()
	}
	wg.Wait()
	mergeAndWrite(input, output)
}
