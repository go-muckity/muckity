package main

import (
	"fmt"
	"github.com/go-muckity/muckity/pkg/muckity"
	"os"
	"time"
)

func main() {
	var w muckity.WorldSystem
	w = muckity.GetWorld()
	err := w.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("Created GenericWorld:", w.String())
	var runner = func() {
		t, err := w.Run()
		if err != nil {
			fmt.Println("crashed at tick:", t, "-", err)
			return
		}
		fmt.Println("recorded ticks:", t)
		return
	}
	go runner()
	time.Sleep(time.Second * 5)
	w.Shutdown()
	os.Exit(0)
}
