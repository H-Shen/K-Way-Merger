package heap

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Node holds the current integer value read from a file,
// along with the file descriptor and a Scanner for further reads.
type Node struct {
	Val     int
	Fd      *os.File
	Scanner *bufio.Scanner
}

// NewNode opens the given file, reads its first integer, and returns a Node.
// If any error occurs, it closes the file before returning.
func NewNode(filename string) (Node, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return Node{}, fmt.Errorf("open file %s: %w", filename, err)
	}

	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		// close on failure to read
		fd.Close()
		if scanErr := scanner.Err(); scanErr != nil {
			return Node{}, fmt.Errorf("scan file %s: %w", filename, scanErr)
		}
		return Node{}, fmt.Errorf("no integer found in file %s", filename)
	}

	val, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fd.Close()
		return Node{}, fmt.Errorf("parse integer in file %s: %w", filename, err)
	}

	return Node{Val: val, Fd: fd, Scanner: scanner}, nil
}

// MinHeap is a slice of Nodes that implements a min-heap ordered by Node.Val.
// It is compatible with container/heap.
type MinHeap []Node

// NewMinHeap returns a pointer to an empty heap with the given initial capacity.
func NewMinHeap(capacity int) *MinHeap {
	h := make(MinHeap, 0, capacity)
	return &h
}

// Len returns the number of elements in the heap.
func (h *MinHeap) Len() int {
	return len(*h)
}

// Less reports whether the element at index i is less than the one at j.
func (h *MinHeap) Less(i, j int) bool {
	return (*h)[i].Val < (*h)[j].Val
}

// Swap swaps the elements at indices i and j.
func (h *MinHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// Push inserts x into the heap. x must be a Node.
func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(Node))
}

// Pop removes and returns the smallest element from the heap.
func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Empty reports whether the heap contains no elements.
func (h *MinHeap) Empty() bool {
	return h.Len() == 0
}
