package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	size := int64(1024 * 1024 * 1024 * 16) // 16GB
	chunkSize := 1024 * 1024 * 256         // 256MB chunks

	file, err := os.OpenFile("random_file", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	start := time.Now()

	var wg sync.WaitGroup
	numWorkers := 4
	ch := make(chan int64, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, chunkSize)
			for offset := range ch {
				rand.Read(buf)
				file.WriteAt(buf, offset)

				progress := float64(offset) / float64(size) * 100
				fmt.Printf("\rProgress: %.2f%%", progress)
			}
		}()
	}

	for offset := int64(0); offset < size; offset += int64(chunkSize) {
		ch <- offset
	}
	close(ch)

	wg.Wait()

	duration := time.Since(start)
	speed := float64(size) / duration.Seconds() / 1024 / 1024 // MB/s
	fmt.Printf("\nCompleted in %v, Average speed: %.2f MB/s\n", duration, speed)
}
