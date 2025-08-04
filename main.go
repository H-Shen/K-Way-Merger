package main

// main is the entry point of the K-Way Merger application. It parses command-line
// arguments and invokes the core application logic.
import (
	"KWayMerger/app"
	"fmt"
	"os"
)

const (
	ARGUMENTS = 3 // Minimum number of arguments required: [program] [input files...] [output file]
)

// showUsage displays the correct usage of the application to the user.
func showUsage() {
	fmt.Printf("Usage: %v [file1 file2 ... fileN] [outputFile]\n", os.Args[0])
	fmt.Println("\nDescription:")
	fmt.Println("  K-Way Merger reads integers from multiple input files, sorts each file,")
	fmt.Println("  and merges them into a single sorted output file using a min-heap.")
	fmt.Println("\nArguments:")
	fmt.Println("  file1 file2 ... fileN   One or more input files containing integers")
	fmt.Println("  outputFile             Path to the output file for merged sorted integers")
}

// main is the entry point function of the application.
// It validates command-line arguments and runs the K-Way Merger application.
func main() {
	// Check if incorrect number of arguments are provided
	if len(os.Args) != ARGUMENTS {
		fmt.Println("Incorrect number of arguments!")
		showUsage()
		os.Exit(1)
	}

	// Extract input files and output file from arguments
	inputFiles := os.Args[1 : len(os.Args)-1]
	outputFile := os.Args[len(os.Args)-1]

	// Run the K-Way Merger application
	if err := app.Run(inputFiles, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
