### K-Way Merger

A toy implementation of K-way Merger in Go. The app: 

* Receives a list of input files and an output file
* Sort integers in each input file by increasing order and rewriting the result to the input file
* Merge all integers in all input files and write the result to the output file using a binary heap.

A possible application scenario: Sorting billions of integers after splitting them into multiple files.

#### Usage

To build the image and execute the unit tests in a docker container, run:

```shell
docker build --no-cache -t kwaymerger-image01 .
docker run --name kwaymerger-container01 kwaymerger-image01
```


#### Requirements

* go 1.17

#### References
* [Direct k-way merge by heap](https://en.wikipedia.org/wiki/K-way_merge_algorithm#Heap)
* [Merge k Sorted Lists](https://leetcode.com/problems/merge-k-sorted-lists)
* [Kth Smallest Element in a Sorted Matrix](https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix)
