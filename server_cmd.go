package main

import (
	"flag"
	. "msa-nd-box-go/server"
)

func main() {

	name := flag.String(string(SERVER), "server", "to set an unique name for this server.")
	d := flag.String(string(DISCOVERY), "", "to turn on service discovery. The param value is service discovery address")
	lb := flag.String(string(LOAD_BALANCER), "", "to turn on a balance strategy. Possible values are robin or random")
	cb := flag.Bool(string(CIRCUIT_BREAKER), false, "to turn on a circuit breaker. Does not involve any values.")
	p := flag.String(string(PORT), "", "to ser a port to start server. The random port by default.")
	r := flag.String(string(RESTART), "", "to restart application if it goes down or cb goes down. This parameter is a path to file")
	rC := flag.String(string(RESTART_COUNT), "", "count of restarts after that stop restarting")
	rP := flag.Bool(string(RESTART_KEEP_PARAMS), false, "if that flag is set, restart will be with the same params")

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

	if *r != "" {
		srv.AddParam(RESTART, *r)
	}

	if *rC != "" {
		srv.AddParam(RESTART_COUNT, *rC)
	}

	if *rP && *r != "" {
		srv.AddParam(RESTART_KEEP_PARAMS,*r)
	}

	srv.Start()
}
