package test

import (
	"KWayMerger/app"
	"bufio"
	"errors"
	"fmt"
	"log"
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
	NUMBER = 1000000
)

func generateRandomNumber(filename string, n int) {
	rand.Seed(time.Now().UnixNano())
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < n; i++ {
		_, err = fmt.Fprintf(fd, "%v\n", rand.Intn(math.MaxInt32-math.MinInt32+1)+math.MinInt32)
		if err != nil {
			log.Fatalln(err)
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

func verify(filename string) bool {
	// read
	fd, err := os.Open(filename)
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
	return sort.IntsAreSorted(list)
}

func createDir(dirName string) error {
	err := os.Mkdir(dirName, 0755)
	if err != nil && os.IsExist(err) {
		fileInfo, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

func TestApp(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	// Create tests
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
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create directories
			dataDir := path + "/data"
			if err = createDir(dataDir); err != nil {
				log.Printf("Failed to create the directory: %v\n", err)
			}
			outputDir := path + "/output"
			if err = createDir(outputDir); err != nil {
				log.Printf("Failed to create the directory: %v\n", err)
			}
			// generate input files
			inputFilesCount := 4
			var wg sync.WaitGroup
			var inputFiles []string
			for i := 1; i <= inputFilesCount; i++ {
				filename := dataDir + "/" + tt.name + "_" + strconv.Itoa(i) + ".txt"
				inputFiles = append(inputFiles, filename)
			}
			// generate random numbers in input files
			wg.Add(inputFilesCount)
			for i := 0; i < inputFilesCount; i++ {
				j := i
				// Generate random numbers and write to file
				go func() {
					defer wg.Done()
					generateRandomNumber(inputFiles[j], NUMBER)
				}()
			}
			wg.Wait()
			outputFile := outputDir + "/" + tt.name + "out.txt"
			app.Run(inputFiles, outputFile)
			// verify the output
			if got := verify(outputFile); got != tt.want {
				t.Errorf("Are numbers in the output file sorted? %v, want %v", got, tt.want)
			}
			// remove all input and output data
			t.Cleanup(func() {
				err = os.RemoveAll(dataDir)
				if err != nil {
					log.Fatalln(err)
				}
				err = os.RemoveAll(outputDir)
				if err != nil {
					log.Fatalln(err)
				}
			})
		})
	}
}
