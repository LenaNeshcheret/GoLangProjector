package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	numbersChannel := make(chan []int)
	minMaxChannel := make(chan [2]int)
	var wg sync.WaitGroup

	wg.Add(2)

	go generateNumbers(5, 300, 10, minMaxChannel, numbersChannel, &wg)
	go findMinMax(numbersChannel, minMaxChannel, &wg)

	wg.Wait()
}

func generateNumbers(a int, b int, amount int, minMaxChan chan [2]int, output chan []int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Generate %d numbers between [ %d, %d ] \n", amount, a, b)
	numbers := make([]int, amount)

	for i := 0; i < amount; i++ {
		numbers[i] = rand.Intn(b-a+1) + a
		fmt.Printf("%d, ", numbers[i])
	}
	output <- numbers

	minMax := <-minMaxChan
	fmt.Printf("\nMAX = %d, \nMIN = %d \n", minMax[1], minMax[0])

	close(output)
	close(minMaxChan)
}

func findMinMax(input chan []int, output chan [2]int, wg *sync.WaitGroup) {
	defer wg.Done()
	numbers := <-input
	if len(numbers) == 0 {
		return
	}
	minVal, maxVal := numbers[0], numbers[0]

	for _, num := range numbers {
		if num < minVal {
			minVal = num
		}
		if num > maxVal {
			maxVal = num
		}
	}
	output <- [2]int{minVal, maxVal}
}
