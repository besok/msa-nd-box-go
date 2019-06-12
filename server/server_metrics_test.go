package server

import (
	"fmt"
	"msa-nd-box-go/message"
	"testing"
)

func TestGaugeStore_Take(t *testing.T) {
	store := CreateGaugeStore(message.Service{
		Service: "Test", Address: "1",
	})

	store.AddGauge(Pulse{})
	mes := store.Take()
	fmt.Println(mes)
}