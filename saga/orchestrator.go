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
		_, _ = w.Write([]byte(fmt.Sprintf(" error %s while parsing json %s \n", err, r.Body)))
		return
	}

	log.Printf("got a saga message:%+v", saga)

	chapters := saga.Chapters
	state := message.Start
	idx := 0
	if len(chapters) == 0 {
		_, _ = w.Write([]byte("there are no chapters"))
		return
	}
	input := chapters[0].Input
	for {
		if idx > len(chapters) - 1  {
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
		log.Printf("next chapter:%+v \n",ch)
		s, ok := o.findService(ch.Service)
		if !ok {
			idx--
			state = message.Rollback
		} else {
			var url string
			switch state {
			case message.Success, message.Start:
				url = fmt.Sprintf("http://%s/%s", s, ch.Chapter)
			case message.Rollback:
				url = fmt.Sprintf("http://%s/%s", s, ch.Rollback)
			default:
				_, _ = w.Write([]byte(state))
				return
			}
			ch.Input = input
			buffer := new(bytes.Buffer)
			_ = json.NewEncoder(buffer).Encode(ch)
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

			input = res.Result
		}
	}

}

func (o *Orch) findService(service string) (string, bool) {
	addr := o.a.GetServiceByName(service)
	if addr == "" {
		return addr, false
	}
	return addr, true
}

