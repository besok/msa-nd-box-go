package main

import (
	"msa-nd-box-go/server"
)

func main() {
	server.CreateAdminServer(
		":9000",
		"C:\\projects\\msa-nd-box-go\\file_storages")
}
