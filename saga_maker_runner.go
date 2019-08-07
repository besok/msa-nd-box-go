package main

import (
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/saga"
)

func main() {
	service := saga.NewSagaService("maker")
	_ = service.New("make", func(ch message.Chapter) message.ChapterResult {
		input := ch.Input
		log.Println("input:", input)
		input = "good"
		return message.ChapterResult{State:message.Success,Result:input}
	})
	_ = service.New("back", func(ch message.Chapter) message.ChapterResult {
		return message.ChapterResult{State:message.Rollback,Result:":("}
	})

	service.Start()
}

