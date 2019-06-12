package main

import . "msa-nd-box-go/server"

func main() {
	serv := CreateServer("test-server", Config{})
	serv.AddGauge(Pulse{})
	serv.AddParam(DISCOVERY, "localhost:9000")
	AddInitOperator()
	serv.Start()
}
