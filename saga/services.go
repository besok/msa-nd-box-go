package saga

import (
	"encoding/json"
	"fmt"
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/server"
	"net/http"
)

var adminServer = "localhost:9002"

type Server struct {
	srv *server.Server
}

func NewSagaService(service string) Server {
	srv := Server{server.CreateServer(service)}
	srv.srv.AddGauge(server.Pulse)
	srv.srv.AddParam(server.DISCOVERY, adminServer)

	return srv
}

func (s *Server) New(chapter string, f func(message.Chapter) message.ChapterResult) error {

	h := func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		w.Header().Set("Content-Type", "application/json")
		var ch message.Chapter
		err := decoder.Decode(&ch)
		if err != nil {
			resp, _ := json.Marshal(message.ChapterResult{State: message.Rollback})
			log.Println("wrong request")
			_, _ = w.Write(resp)
			return
		}

		f(ch)
		resp, _ := json.Marshal(f(ch))
		_, _ = w.Write(resp)
	}

	s.srv.AddHandler(fmt.Sprintf("/%s",chapter), h)
	return nil
}

func (s *Server) Start(){
	s.srv.Start()
}