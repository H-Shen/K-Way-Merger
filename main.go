package main

import (
	"fmt"
	"log"
	"os"
	"untitled/app"
)

func showUsage() {
	fmt.Printf("Usage: %v [file1 file2 ... fileN] [outputFile]\n", os.Args[0])
}

func main() {
	// Tune the log setting
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Check arguments
	if len(os.Args) < 3 {
		fmt.Println("Insufficient arguments!")
		showUsage()
		os.Exit(1)
	}
	app.Run(os.Args[1:(len(os.Args)-1)], os.Args[len(os.Args)-1])
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	return nil
}
