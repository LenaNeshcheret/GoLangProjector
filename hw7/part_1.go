package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	numChan := make(chan int)
	avgChan := make(chan float64)
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		generateRandomNumbers(numChan)
		close(numChan)
	}()
	go func() {
		defer wg.Done()
		calculateAverage(numChan, avgChan)
		close(avgChan)
	}()
	go func() {
		defer wg.Done()
		printAverage(avgChan)
	}()

	wg.Wait()
}

func generateRandomNumbers(intCh chan int) {

	for i := 0; i < 10; i++ {
		num := rand.Intn(100)
		fmt.Printf("%d, ", num)
		intCh <- num
		time.Sleep(2 * time.Second)
	}
}

func calculateAverage(input chan int, output chan float64) {
	var numbers []int
	for num := range input {
		numbers = append(numbers, num)
		fmt.Printf("Random numbers: %d,", numbers)
		total := 0
		for _, n := range numbers {
			total += n
		}
		average := float64(total) / float64(len(numbers))
		output <- average
	}
}

func printAverage(output chan float64) {
	for avg := range output {
		fmt.Printf(" Average: %.2f\n", avg)
	}
}
