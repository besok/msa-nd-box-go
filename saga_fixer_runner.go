package main

import (
	"log"
	"msa-nd-box-go/message"
	"msa-nd-box-go/saga"
)

func main() {
	service := saga.NewSagaService("fixer")
	_ = service.New("fix", func(ch message.Chapter) message.ChapterResult {
		input := ch.Input
		log.Println("input:", input)
		input = "fix"
		return message.ChapterResult{State:message.Success,Result:input}
	})
	_ = service.New("back", func(ch message.Chapter) message.ChapterResult {
		return message.ChapterResult{State:message.Rollback,Result:":("}
	})

	service.Start()

}

