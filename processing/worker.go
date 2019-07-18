package processing

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"msa-nd-box-go/server"
	"net/http"
	"os"
	"sync"
)

const PATH server.Param = "path"

type Result struct {
	Count int
	Error error
}

var (
	isBusy = false
	mutex  = sync.Mutex{}
	adminServer = "localhost:9000"
)

func InitWorker() {
	srv := server.CreateServer("worker-service")
	srv.AddGauge(server.Pulse)
	srv.AddGauge(State)

	srv.AddParam(server.DISCOVERY, adminServer)
	srv.AddParam(PATH, "path")

	srv.AddHandler("/task", taskHandler)

	srv.Start()
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	bts, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(" error while reading a body:",e)
	}
	go ProcessTask(string(bts))
}
func ProcessTask(p string) {
	isBusy = true

	i, e := Task(p)
	result := Result{i, e}

	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(result)
	_, e = http.Post(fmt.Sprint("http://",adminServer, "/task"), "application/json; charset=utf-8", buffer)
	log.Printf(" send a request to admin server: result:%s, result error:%s",result.Count, e)
}

func Task(pathFile string) (int, error) {
	f, err := os.Open(pathFile)
	defer f.Close()

	if err != nil {
		return 0, err
	}

	count := 0
	r := bufio.NewReader(f)
	b := make([]byte, 1024)

	for {
		n, err := r.Read(b)
		if err != nil && err != io.EOF {
			return count, err
		}
		if err != nil && err == io.EOF {
			return count, nil
		}

		if n == 0 {
			return count, nil
		}

		count += n
	}
}

func State() (k string, v string, err error) {
	val := "false"
	if isBusy {
		val = "true"
	}

	return "state", val, nil
}
