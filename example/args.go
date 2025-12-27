package main

import (
	"fmt"

	"github.com/mykytaserdiuk/fluxo/event"
)

const (
	newFollowerE = "main:new_follower"
)

func mainArgs() {
	bus := event.New()

	// With args in func
	err := bus.Subscribe(newFollowerE, func(name string) {
		fmt.Println("NEW FOLLOWER: ", name)
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Without args
	err = bus.Subscribe(newFollowerE, func() {
		fmt.Println("You have new follower!")
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// With an abundance of arguments
	err = bus.Subscribe(newFollowerE, func(name, lastName string) {
		fmt.Println("Follower name : ", name, lastName)
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	bus.Emit(newFollowerE, "hello")

	// OUTPUT:
	//
	// Follower name :  hello
	// You have new follower!
	// NEW FOLLOWER:  hello
}
