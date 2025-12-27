package main

import (
	"fmt"
	"time"

	"github.com/mykytaserdiuk/fluxo"
)

const (
	onExit = "manager:on_exit"
)

func mainOnce() {
	bus := fluxo.New()
	t := time.Date(2025, 12, 31, 12, 12, 12, 12, time.UTC)

	bus.SubscribeOnce(onExit, func(t time.Time) {
		fmt.Println("Once sub: " + t.String())
	})
	bus.Subscribe(onExit, func(t time.Time) {
		fmt.Println("Simple sub: " + t.String())
	})
	bus.Emit(onExit, t)
	bus.Emit(onExit, t)

	// OUTPUT:
	//
	// Simple sub: 2025-12-31 12:12:12.000000012 +0000 UTC
	// Once sub: 2025-12-31 12:12:12.000000012 +0000 UTC
	// Simple sub: 2025-12-31 12:12:12.000000012 +0000 UTC
}
