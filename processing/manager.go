package processing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/server"
	"net/http"
	"os"
	"path/filepath"
)

type Manager struct {
	a *server.AdminServer
	batchSize int
}
type taskQ []string

func (q *taskQ) empty() bool {
	return len(*q) == 0
}

func (q *taskQ) push(p string) {
	*q = append(*q, p)
}

func (q taskQ) pop() string {
	el := q[len(q)-1]
	q = q[:len(q)-1]
	return el
}

type queue interface {
	empty() bool
	push(p string)
	pop() string
}

var workerPath string
var resList = make([]Result, 0)
var qTask = taskQ(make([]string, 0))

func InitManager(batchSize int, wPath string) {
	workerPath = wPath
	adm := Manager{server.CreateAdminServer("file_storages_procesing"),batchSize}

	server.AddParamHandler(processPath)
	server.NewMetricHandler(processFreeWorkers)

	adm.a.AddHandler("/task", processResult)
	adm.a.AddHandler("/res", showResult)
	adm.a.AddHandler("/start", adm.startProcess)

	adm.a.Start("9001")
}

func processFreeWorkers(a *server.AdminServer, m message.MetricsMessage) error {
	metric, ok := m.Metrics["state"]
	address := m.From.Address
	if !ok {
		_, _ = http.Get(fmt.Sprint("http://", address, "/close"))
		return nil
	}

	if metric.Value == "false" {
		if !qTask.empty() {
			p := qTask.pop()
			_, err := http.Post(fmt.Sprint("http://", address, "/close"), " text/plain",
				bytes.NewReader([]byte(p)))
			if err != nil {
				fmt.Println("error to send a new task")
			}

		} else {
			log.Println("the q is empty, it needs to close service")
			_, _ = http.Get(fmt.Sprint("http://", address, "/close"))
		}
	}

	return nil
}

func processResult(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var res Result
	err := decoder.Decode(&res)
	if err != nil {
		log.Fatalf(" error %s while parsing json %s \n", err, r.Body)
		return
	}

	resList = append(resList, res)

}

func (m *Manager) startProcess(w http.ResponseWriter, r *http.Request) {
	p, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(" error while reading a body:", e)
	}

	if e = filepath.Walk(string(p), addFileToQ); e != nil {
		log.Println(" error while going to folder", e)
	}

	lines, ok := m.a.Storage(server.RegistryStorage).Get("worker-service")

	startNum := 0
	if !ok {
		startNum = m.batchSize
	} else {
		startNum = m.batchSize - lines.Size()
	}

	for i := 0; i < startNum; i++ {
		if err := m.a.StartNewCommand(workerPath); err != nil {
			log.Println(" error to run worker, error: ", err)
		}
	}

}

func addFileToQ(p string, i os.FileInfo, err error) error {
	if i.IsDir() {
		return nil
	}
	qTask.push(p)
	log.Println("add file to q : ", p)
	return nil
}

func showResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	js, e := json.Marshal(resList)
	if e != nil {
		log.Fatalf("can't convert to json, %s", e)
	}
	_, e = w.Write(js)
	if e != nil {
		log.Fatalf("can't send, %s", e)
	}
}

func processPath(_ *server.AdminServer, _ message.Service, k string, v string) error {
	if k == string(PATH) {
		workerPath = v
		log.Println("set new worker path: ", v)
	}
	return nil
}
