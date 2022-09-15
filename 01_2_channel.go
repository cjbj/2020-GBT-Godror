package main

import (
	"fmt"
)

func main() {
	c := make(chan int)

  // Goroutines run in parallel
	go calc_sum(c, []int{1, 2, 3, 4, 5})
	go calc_sum(c, []int{8, -1, 0, 15, -25, 6, 2})
	go calc_sum(c, []int{1})
  
	for i := 0; i < 3; i++ {
		sum := <- c                                  // read from channel
		fmt.Printf("Calculated sum is %v\n", sum)
	}
}

func calc_sum(c chan int, arr []int) {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	c <- sum                                       // write into channel
}
