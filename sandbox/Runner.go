package main

import (
	"msa-nd-box-go/server"
	"net/http"
)

func main() {
	server.RegisterService()
	_  = http.ListenAndServe("localhost:9000", nil)
}
