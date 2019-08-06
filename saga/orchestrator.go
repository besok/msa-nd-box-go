package saga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/server"
	"net/http"
)

type Orch struct {
	a *server.AdminServer
}

func RunSagaOrch() {
	o := Orch{server.CreateAdminServer("file_storages_saga")}

	o.a.AddHandler("/saga", o.handlerSaga)

	o.a.Start("9002")
}

func (o *Orch) handlerSaga(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var saga message.SagaMessage
	err := decoder.Decode(&saga)
	if err != nil {
		_, _ = w.Write([]byte("can not process"))
		log.Fatalf(" error %s while parsing json %s \n", err, r.Body)
	}

	log.Printf("got a saga message:%+v", saga)

	chapters := saga.Chapters
	state := message.Start
	idx := 0
	if len(chapters) == 0 {
		log.Println("saga does not have chapters, terminate.")
		_, _ = w.Write([]byte("there are no chapters"))
		return
	}

	for ok := true; ok; ok = checkSagaRes(state) {
		if idx > len(chapters) {
			state = message.Finish
			_, _ = w.Write([]byte(state))
			return
		}
		if idx < 0 {
			state = message.Abort
			_, _ = w.Write([]byte(state))
			return
		}

		ch := chapters[idx]

		s, ok := o.findService(ch.Service)
		if !ok {
			idx--
			state = message.Rollback
		} else {
			var url string
			switch state {
			case message.Success:
			case message.Start:
				url = fmt.Sprintf("http://%s/%s", s, ch.Chapter)
			case message.Rollback:
				url = fmt.Sprintf("http://%s/%s", s, ch.Rollback)
			default:
				_, _ = w.Write([]byte(state))
				return
			}

			buffer := new(bytes.Buffer)
			_ = json.NewEncoder(buffer).Encode(sagaResultStart(ch))
			resp, err := http.Post(url, "application/json; charset=utf-8", buffer)
			if err != nil {
				log.Println("got error: ", err)
				idx--
				state = message.Rollback
				continue
			}

			decoder := json.NewDecoder(resp.Body)
			var res message.ChapterResult
			_ = decoder.Decode(&res)

			switch res.State {
			case message.Success:
				state = message.Success
				idx++
			case message.Rollback:
				state = message.Rollback
				idx--
			default:
				_, _ = w.Write([]byte(state))
				return
			}
		}
	}

}

func checkSagaRes(state message.ChapterState) bool {
	return state != message.Abort || state != message.Finish
}

func (o *Orch) findService(service string) (string, bool) {
	addr := o.a.GetServiceByName(service)
	if addr == "" {
		return addr, false
	}
	return addr, true
}

func sagaResultStart(ch message.Chapter) *message.ChapterResult {
	return message.NewChapterResult(ch, message.Start)
}
