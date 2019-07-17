package processing

import (
	"msa-nd-box-go/server"
	"net/http"
	"sync"
)

const PATH server.Param = "path"

var (
	isBusy = false
	mutex  = sync.Mutex{}
)

func InitWorker() {
	srv := server.CreateServer("worker-service")
	srv.AddGauge(server.Pulse)
	srv.AddGauge(State)

	srv.AddParam(server.DISCOVERY, "localhost:9001")
	srv.AddParam(PATH, "path")

	srv.AddHandler("task", taskHandler)

}

func taskHandler(w http.ResponseWriter, r *http.Request) {

}

func State() (k string, v string, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	val := "false"
	if isBusy {
		val = "true"
	}

	return "state", val, nil
}
