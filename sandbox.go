package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)
	s := 0
	for i := 0; i < 20; i++ {
		go func(ch chan int) {
			//time.Sleep(time.Second*1)
			ch <- 10
		}(ch)
	}

	for i:=0; i< 20; i++{
		el := <-ch
		s+=el
	}
	fmt.Println(s)
}
