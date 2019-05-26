package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex

func main() {
	http.HandleFunc("/",
		func(writer http.ResponseWriter, request *http.Request) {
			mu.Lock()
			c.c++
			mu.Unlock()
		})
	http.HandleFunc("/counter", counter)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

func counter(writer http.ResponseWriter, request *http.Request) {
	mu.Lock()
	_, _ = fmt.Fprintf(writer, "counter = %s\n", c)
	mu.Unlock()
}


type Counter struct {
	c int
	n string
}

var c Counter = Counter{0, "counter"}