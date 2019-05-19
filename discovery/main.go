package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		return
	}
	world := os.Args[1]
	file := os.Args[2]
	f, e := os.Open(file)
	if e != nil {
		fmt.Printf("%s\n", e)
		return
	}
	res := findLine(f, &world)

	for _, t := range res {
		if len(t) > 0 {
			fmt.Printf("%s \n", t)
		}
	}

}

func findLine(file *os.File, world *string) []string {
	scanner := bufio.NewScanner(file)
	res := make([]string, 10)
	i := 0
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, *world) {
			res[i] = text
			i++
		}
	}
	return res
}
