package main

import (
	"msa-nd-box-go/server"
)

func main() {
	adminServer := server.CreateAdminServer("C:\\projects\\msa-nd-box-go\\file_storages")
	adminServer.Start("9000")
}
