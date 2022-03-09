package app

import (
	"bufio"
	"container/heap"
	"errors"
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

func NewNode(filename string) (Node, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return Node{}, err
	}
	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanWords)
	if scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return Node{}, err
		}
		return Node{
			val:     val,
			fd:      fd,
			scanner: scanner,
		}, nil
	}
	if err = scanner.Err(); err != nil {
		return Node{}, err
	}
	return Node{}, errors.New("failed to scan the first integer")
}

// Merge all numbers in all files and output to 'filename'
func mergeAndWrite(input []string, output string) {
	fd, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	var minHeap MinHeap
	heap.Init(&minHeap)
	for i := 0; i < len(input); i++ {
		node, err := NewNode(input[i])
		if err == nil {
			heap.Push(&minHeap, node)
		}
	}
	for !minHeap.Empty() {
		node := heap.Pop(&minHeap).(Node)
		// write the val to the output file
		_, err = fmt.Fprintf(fd, "%v\n", node.val)
		if err != nil {
			log.Fatalln(err)
		}
		if node.scanner.Scan() {
			val, err := strconv.Atoi(node.scanner.Text())
			if err != nil {
				log.Fatalln(err)
			}
			node.val = val
			heap.Push(&minHeap, node)
		} else {
			err := node.fd.Close()
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
