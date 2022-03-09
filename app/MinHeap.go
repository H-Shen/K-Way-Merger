package app

import (
	"bufio"
	"os"
)

// Node File node
type Node struct {
	val     int
	fd      *os.File
	scanner *bufio.Scanner
}

// MinHeap Implementation of a min-heap that takes Node's val */
type MinHeap []Node

func (minHeap *MinHeap) Len() int {
	return len(*minHeap)
}

func (minHeap *MinHeap) Less(i, j int) bool {
	return (*minHeap)[i].val < (*minHeap)[j].val
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
