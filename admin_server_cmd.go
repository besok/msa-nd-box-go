package main

import (
	"flag"
	"log"
	"msa-nd-box-go/server"
)

func main() {

	name := flag.String("storage", "", "to set a folder to locate file storages")
	p := flag.String("port", "", "to ser a port to start server. It is an appropriate value")

	flag.Parse()

	if *name == "" || *p == "" {
		log.Fatalf(" the appropriate parameters haven't been set. Please use -h to get more information.")
	}

	adminServer := server.CreateAdminServer(*name)
	adminServer.Start(*p)

}
