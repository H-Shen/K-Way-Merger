package heap

import (
	"bufio"
	"errors"
	"os"
	"strconv"
)

// Node File node
type Node struct {
	Val     int
	Fd      *os.File
	Scanner *bufio.Scanner
}

// MinHeap Implementation of a min-heap that takes Node's Val */
type MinHeap []Node

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
			Val:     val,
			Fd:      fd,
			Scanner: scanner,
		}, nil
	}
	if err = scanner.Err(); err != nil {
		return Node{}, err
	}
	return Node{}, errors.New("failed to scan the first integer")
}

func (minHeap *MinHeap) Len() int {
	return len(*minHeap)
}

func (minHeap *MinHeap) Less(i, j int) bool {
	return (*minHeap)[i].Val < (*minHeap)[j].Val
}

func (minHeap *MinHeap) Swap(i, j int) {
	(*minHeap)[i], (*minHeap)[j] = (*minHeap)[j], (*minHeap)[i]
}

func (minHeap *MinHeap) Push(val interface{}) {
	*minHeap = append(*minHeap, val.(Node))
}

func (minHeap *MinHeap) Pop() interface{} {
	oldLength := len(*minHeap)
	val := (*minHeap)[oldLength-1]
	*minHeap = (*minHeap)[:(oldLength - 1)]
	return val
}

func (minHeap *MinHeap) Empty() bool {
	return len(*minHeap) == 0
}
