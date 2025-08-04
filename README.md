### K-Way Merger

A Go implementation of the K-way merge algorithm for efficiently merging multiple sorted files into a single sorted output file. This implementation uses a min-heap to achieve optimal performance.

## Overview

The K-Way Merger application:

1. Reads multiple input files containing integers
2. Sorts each input file in parallel using goroutines
3. Merges the sorted files into a single output file using a min-heap
4. Handles errors gracefully throughout the process

This implementation is particularly useful for scenarios where you need to sort very large datasets that don't fit into memory (external sorting), by splitting the data into smaller chunks, sorting each chunk, and then merging them.

## Architecture

The application consists of the following main components:

1. **app package**: Contains the core logic for reading, sorting, and merging files
   - `readSortRewrite`: Reads a file, sorts its integers, and rewrites the sorted data
   - `mergeAndWrite`: Merges multiple sorted files into one using a min-heap
   - `Run`: Orchestrates the parallel sorting and merging process

2. **heap package**: Implements a min-heap data structure for efficient merging
   - `Node`: Represents a file and its current value
   - `MinHeap`: Implements the min-heap interface for sorting nodes by value

3. **main.go**: Entry point that parses command-line arguments and invokes the core logic

## Usage

### Command Line

To use the K-Way Merger, run:

```shell
# Build the application
go build -o kwaymerger

# Run with input files and output file
./kwaymerger [file1 file2 ... fileN] [outputFile]
```

Example:

```shell
./kwaymerger input1.txt input2.txt input3.txt output.txt
```

### Docker

To build and run using Docker:

```shell
# Build the Docker image
 docker build --no-cache -t kwaymerger .

# Run the container with input files mounted
 docker run -v /path/to/input/files:/data kwaymerger /data/file1.txt /data/file2.txt /data/output.txt
```

## Running Tests

To run the unit tests:

```shell
# Run all tests
 go test -v ./test

# Run with Docker
 docker run kwaymerger go test -v ./test
```

## Performance Considerations

- The algorithm efficiently merges K sorted files with a time complexity of O(N log K), where N is the total number of elements
- Parallel sorting of input files improves performance on multi-core systems
- Memory usage is optimized by processing files sequentially and using a heap of size K

## Requirements

* Go 1.22 or higher

## References

* [K-way merge algorithm](https://en.wikipedia.org/wiki/K-way_merge_algorithm)
* [External merge sort](https://en.wikipedia.org/wiki/External_sorting)
* [Merge k Sorted Lists](https://leetcode.com/problems/merge-k-sorted-lists)
* [Kth Smallest Element in a Sorted Matrix](https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix)
