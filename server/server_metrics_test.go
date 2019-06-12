package server

import (
	"fmt"
	"msa-nd-box-go/message"
	"testing"
)

func TestGaugeStore_Take(t *testing.T) {
	store := CreateGaugeStore()

	store.AddGauge(Pulse{})
	mes := store.Take(message.Service{
		Service: "Test", Address: "1",
	})
	fmt.Println(mes)
}

