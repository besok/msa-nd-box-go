package main

import (
	"log"
	"msa-nd-box-go/server"
	"msa-nd-box-go/storage"
)

func main() {

	doNothingListener := func(event storage.StorageEvent, storageName storage.StorageName,key string, value storage.Line){
		log.Printf("do nothing listener, event:%s, name:%s",event,storageName)
	}

	server.CreateAdminServer("C:\\projects\\msa-nd-box-go\\file_storages",doNothingListener)
}
