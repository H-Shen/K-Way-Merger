package test

// Package test contains unit tests for the K-Way Merger application.
// It includes tests for the core functionality of reading, sorting,
// and merging files, as well as helper functions for generating test data
// and verifying results.
import (
	"KWayMerger/app"
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	NUMBER = 1000000 // Number of random integers to generate per test file
)

// generateRandomNumber generates n random integers and writes them to the specified file.
// The integers span the full int32 range.
func generateRandomNumber(filename string, n int) error {
	// create a new, locally seeded RNG
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// open (or create) the file
	fd, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer fd.Close()

	// buffer writes for performance
	buf := bufio.NewWriter(fd)

	for i := 0; i < n; i++ {
		// rng.Int31() ∈ [0, 2³¹)
		// + MinInt32 shifts it to [−2³¹, 2³¹−1]
		num := rng.Int31() + int32(math.MinInt32)

		if _, err := fmt.Fprintln(buf, num); err != nil {
			return fmt.Errorf("failed to write to file %s: %w", filename, err)
		}
	}

	// flush buffered data
	if err := buf.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer to file %s: %w", filename, err)
	}

	// ensure on-disk sync
	if err := fd.Sync(); err != nil {
		return fmt.Errorf("failed to sync file %s: %w", filename, err)
	}

	return nil
}

// verify checks if the integers in the specified file are sorted in ascending order.
// It reads all integers from the file and verifies their order.
//
// Parameters:
//
//	filename - The path to the file to be verified
//
// Returns:
//
//	bool - True if the integers are sorted, false otherwise
//	error - Any error encountered during file operations or parsing
func verify(filename string) (bool, error) {
	// Open file for reading
	fd, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	// Ensure file is closed when function exits
	defer func() {
		closeErr := fd.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", filename, closeErr)
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
			return false, fmt.Errorf("failed to parse number in file %s: %w", filename, err)
		}
		list = append(list, temp)
	}
	if err = scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	// Check if the list is sorted
	return sort.IntsAreSorted(list), nil
}

// createDir creates a directory with the specified name if it doesn't exist.
// If the directory already exists, it verifies that it is indeed a directory.
//
// Parameters:
//
//	dirName - The path to the directory to be created
//
// Returns:
//
//	error - Any error encountered during directory creation or verification
func createDir(dirName string) error {
	err := os.Mkdir(dirName, 0755)
	if err != nil && os.IsExist(err) {
		fileInfo, statErr := os.Stat(dirName)
		if statErr != nil {
			return statErr
		}
		if !fileInfo.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

// TestApp tests the K-Way Merger application with multiple test cases.
// Each test case generates random input files, runs the application,
// and verifies that the output file contains sorted integers.
func TestApp(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Log(err)
	}
	// Define test cases
	tests := []struct {
		name string
		want bool
	}{
		{name: "test1", want: true},
		{name: "test2", want: true},
		{name: "test3", want: true},
		{name: "test4", want: true},
		{name: "test5", want: true},
		{name: "test6", want: true},
		{name: "test7", want: true},
		{name: "test8", want: true},
	}
	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create data and output directories
			dataDir := path + "/data"
			if err = createDir(dataDir); err != nil {
				t.Fatalf("Failed to create data directory: %v\n", err)
			}
			outputDir := path + "/output"
			if err = createDir(outputDir); err != nil {
				t.Fatalf("Failed to create output directory: %v\n", err)
			}
			// Generate input file paths
			inputFilesCount := 4
			var wg sync.WaitGroup
			var inputFiles []string
			for i := 1; i <= inputFilesCount; i++ {
				filename := dataDir + "/" + tt.name + "_" + strconv.Itoa(i) + ".txt"
				inputFiles = append(inputFiles, filename)
			}
			// Generate random numbers in input files (parallel)
			wg.Add(inputFilesCount)
			errCh := make(chan error, inputFilesCount)
			for i := 0; i < inputFilesCount; i++ {
				j := i
				go func() {
					defer wg.Done()
					// Generate random numbers and write to file
					if err := generateRandomNumber(inputFiles[j], NUMBER); err != nil {
						errCh <- err
					}
				}()
			}

			// Close error channel once all goroutines are done
			go func() {
				wg.Wait()
				close(errCh)
			}()

			// Collect all errors from random number generation
			var errs []error
			for err := range errCh {
				errs = append(errs, err)
			}

			// If there are errors, fail the test
			if len(errs) > 0 {
				t.Fatalf("%d errors occurred during random number generation: %v", len(errs), errs)
			}
			// Define output file path
			outputFile := outputDir + "/" + tt.name + "out.txt"
			// Run the K-Way Merger application
			if err := app.Run(inputFiles, outputFile); err != nil {
				t.Fatalf("Failed to run K-Way Merger: %v", err)
			}

			// Verify the output file is sorted
			got, err := verify(outputFile)
			if err != nil {
				t.Fatalf("Failed to verify output file: %v", err)
			}
			if got != tt.want {
				t.Errorf("Are numbers in the output file sorted? %v, want %v", got, tt.want)
			}
			// Clean up test data after the test
			t.Cleanup(func() {
				err = os.RemoveAll(dataDir)
				if err != nil {
					t.Errorf("Failed to remove data directory: %v", err)
				}
				err = os.RemoveAll(outputDir)
				if err != nil {
					t.Errorf("Failed to remove output directory: %v", err)
				}
			})
		})
	}
}
