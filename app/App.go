package app

// Package app provides the core functionality for the K-Way Merger application.
// It handles reading input files, sorting their contents, and merging them
// using a min-heap to produce a single sorted output file.
import (
	myHeap "KWayMerger/heap"
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
)

// readSortRewrite reads integers from a file, sorts them, and rewrites the sorted
// integers back to the same file. Each integer is read as a separate word.
//
// Parameters:
//
//	file - The path to the file to be read, sorted, and rewritten
//
// Returns:
//
//	error - Any error encountered during reading, sorting, or writing
func readSortRewrite(file string) error {
	// Open file for reading
	fd, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file, err)
	}
	// Ensure file is closed when function exits
	defer func() {
		closeErr := fd.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", file, closeErr)
		}
	}()

	// Read integers from file
	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanWords) // Split on whitespace
	var temp int
	var list []int
	for scanner.Scan() {
		temp, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("failed to parse number in file %s: %w", file, err)
		}
		list = append(list, temp)
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", file, err)
	}

	// Sort the integers
	sort.Ints(list)

	// Open file for writing (truncate existing content)
	fd2, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s for writing: %w", file, err)
	}
	// Ensure file is closed when function exits
	defer func() {
		closeErr := fd2.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", file, closeErr)
		}
	}()

	// Write sorted integers back to file
	for i := 0; i < len(list); i++ {
		_, err = fmt.Fprintf(fd2, "%v\n", list[i])
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", file, err)
		}
	}

	// Sync file to ensure data is written to disk
	err = fd2.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file %s: %w", file, err)
	}

	return nil
}

// mergeAndWrite merges integers from multiple sorted input files into a single
// sorted output file using a min-heap. It reads the smallest available value
// from each input file, adds it to the heap, and then extracts the minimum
// value to write to the output file.
//
// Parameters:
//
//	input - Slice of paths to the input files containing sorted integers
//	output - Path to the output file where merged sorted integers will be written
//
// Returns:
//
//	error - Any error encountered during merging or writing
func mergeAndWrite(input []string, output string) error {
	// Open output file for writing
	fd, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file %s: %w", output, err)
	}
	// Ensure output file is closed when function exits
	defer func() {
		closeErr := fd.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close output file %s: %w", output, closeErr)
		}
	}()

	// Initialize min-heap
	var minHeap myHeap.MinHeap
	heap.Init(&minHeap)

	// Track all open files for proper cleanup
	var openFiles []*os.File
	defer func() {
		// Close all remaining open files on error
		for _, f := range openFiles {
			f.Close()
		}
	}()

	// Create nodes for each input file and add to heap
	for i := 0; i < len(input); i++ {
		node, err := myHeap.NewNode(input[i])
		if err != nil {
			return fmt.Errorf("failed to create node for file %s: %w", input[i], err)
		}
		openFiles = append(openFiles, node.Fd)
		heap.Push(&minHeap, node)
	}

	// Merge process: extract minimum value from heap and write to output
	for !minHeap.Empty() {
		node := heap.Pop(&minHeap).(myHeap.Node)
		// Write the smallest value to output file
		_, err = fmt.Fprintf(fd, "%v\n", node.Val)
		if err != nil {
			return fmt.Errorf("failed to write to output file %s: %w", output, err)
		}

		// Read next value from the same file if available
		if node.Scanner.Scan() {
			val, parseErr := strconv.Atoi(node.Scanner.Text())
			if parseErr != nil {
				return fmt.Errorf("failed to parse number in file %s: %w", node.Fd.Name(), parseErr)
			}
			node.Val = val
			heap.Push(&minHeap, node) // Reinsert node with new value
		} else {
			// File is exhausted, remove from open files list
			for i, f := range openFiles {
				if f == node.Fd {
					openFiles = append(openFiles[:i], openFiles[i+1:]...)
					break
				}
			}

			// Close the file
			closeErr := node.Fd.Close()
			if closeErr != nil {
				return fmt.Errorf("failed to close file: %w", closeErr)
			}

			// Check for scanner errors
			if scanErr := node.Scanner.Err(); scanErr != nil {
				return fmt.Errorf("error reading file: %w", scanErr)
			}
		}
	}

	// Sync output file to ensure data is written to disk
	err = fd.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync output file %s: %w", output, err)
	}

	return nil
}

// Run executes the K-Way Merger application. It processes input files in parallel,
// sorting each file, then merges the sorted files into a single output file.
//
// Parameters:
//
//	input - Slice of paths to the input files containing integers
//	output - Path to the output file where merged sorted integers will be written
//
// Returns:
//
//	error - Any error encountered during processing or merging
func Run(input []string, output string) error {
	// Channel to collect errors from goroutines
	errCh := make(chan error, len(input))
	var wg sync.WaitGroup

	wg.Add(len(input))
	// Process each input file in parallel
	for i := 0; i < len(input); i++ {
		j := i
		go func() {
			defer wg.Done()
			// Read, sort, and rewrite each file
			if err := readSortRewrite(input[j]); err != nil {
				errCh <- fmt.Errorf("error processing file %s: %w", input[j], err)
			}
		}()
	}

	// Close error channel once all goroutines are done
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Collect all errors from goroutines
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	// If there are errors from any goroutine, return them
	if len(errs) > 0 {
		return fmt.Errorf("%d errors occurred during processing: %v", len(errs), errs)
	}

	// Merge the sorted files into the output file
	if err := mergeAndWrite(input, output); err != nil {
		return fmt.Errorf("failed to merge files: %w", err)
	}

	return nil
}
