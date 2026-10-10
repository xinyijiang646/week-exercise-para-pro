package main

import (
	"fmt"
	"sync"
        "time"
        "os"
        "strconv"
        "path/filepath"
        "math"
)

func dotSequential(a []int32, b []int32) int32{
    var sum int32
    numElements := len(a)
    for i:=0; i<numElements; i++ {
        sum += int32(math.Abs(math.Sin(float64(a[i]))*float64(b[i])))
    }
    return sum
}

func dotChunk(a []int32, b []int32, sum *int32, wg *sync.WaitGroup) {
    numElements := len(a)
    for i:=0; i<numElements; i++ {
        *sum += int32(math.Abs(math.Sin(float64(a[i]))*float64(b[i])))
    }
    wg.Done()
}

func dotParallel(a []int32, b []int32, numThreads int) int32{
    sum := make([]int32, numThreads*16)
    numElements := len(a)
    numElementsPerThread := int(math.Ceil(float64(numElements) / float64(numThreads)))
    var wg sync.WaitGroup
    for i:=0; i<numThreads; i++ {
        chunkStart := numElementsPerThread * i
        if(chunkStart > numElements) {
            chunkStart = numElements
        }
        chunkEnd := chunkStart + numElementsPerThread
        if(chunkEnd > numElements) {
            chunkEnd = numElements
        }
        wg.Add(1)
        go dotChunk(a[chunkStart:chunkEnd], b[chunkStart:chunkEnd], &sum[i*16], &wg)
    }
    wg.Wait()
    for i:=1; i<numThreads; i++ {
        sum[0] += sum[i*16];
    }
    return sum[0]
}

func main() {
        if(len(os.Args) < 3) {
            panic(fmt.Sprintf("Usage: go run %s.go <problemsize> <threadcount>\n", filepath.Base(os.Args[0])))
        }
        numElements, err := strconv.Atoi(os.Args[1])
        if(err!=nil) {
            panic(err)
        }
        numThreads, err := strconv.Atoi(os.Args[2])
        if(err!=nil) {
            panic(err)
        }

        a := make([]int32, numElements)
        b := make([]int32, numElements)


        for i:=0; i<numElements; i++ {
            a[i] = int32(i)
            b[i] = int32(numElements-i)
        }

        start := time.Now()
        dot_sequential := dotSequential(a[:], b[:])
        duration := time.Since(start)
        fmt.Printf("computed %d, sequential, time: %s\n", dot_sequential, duration)

        start = time.Now()
        dot_parallel := dotParallel(a[:], b[:], numThreads)
        duration = time.Since(start)
        fmt.Printf("computed %d, %d threads, time: %s\n", dot_parallel, numThreads, duration)
}
