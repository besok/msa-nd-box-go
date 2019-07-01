package main

import (
	"flag"
	. "msa-nd-box-go/server"
)

func main() {

	name := flag.String("server", "server", "to set an unique name for this server.")
	d := flag.String("discovery", "", "to turn on service discovery. The param value is service discovery address")
	lb := flag.String("balance", "", "to turn on a balance strategy. Possible values are robin or random")
	cb := flag.Bool("cb", false, "to turn on a circuit breaker. Does not involve any values.")
	p := flag.String("port", "", "to ser a port to start server. The random port by default.")

	flag.Parse()

	srv := CreateServer(*name)
	srv.AddGauge(Pulse)
	srv.AddGauge(CircuitBreaker)

	if *d != "" {
		srv.AddParam(DISCOVERY, *d)
	}
	if *p != "" {
		srv.AddParam(PORT, *p)
	}

	if *cb {
		srv.AddParam(CIRCUIT_BREAKER, "true")
	}
	if *lb != "" {
		srv.AddParam(LOAD_BALANCER, *lb)
	}

	srv.Start()
}
