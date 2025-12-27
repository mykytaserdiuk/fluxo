package main

import (
	"fmt"

	"github.com/mykytaserdiuk/fluxo"
)

const (
	newReport = "hr:new_report"
)

func main() {
	bus := fluxo.NewEventBus()
	handler := func(data string) {
		fmt.Println("New message from HR: ", data)
	}

	err := bus.Subscribe(newReport, handler)
	if err != nil {
		fmt.Println(err.Error())
	}

	bus.Emit(newReport, "we have a new sys.admin")
	bus.Emit(newReport, "we need more computer in room 29")

	_ = bus.Unsubscribe(newReport, handler)

	bus.Emit(newReport, "we have some problem etc.") // we dont receive this message

	err = bus.Unsubscribe(newReport, func(data string) {})
	if err != nil {
		fmt.Println(err.Error())
	}
	// OUTPUT:
	//
	// New message from HR:  we have a new sys.admin
	// New message from HR:  we need more computer in room 29
	// id hasn`t handlers
}
