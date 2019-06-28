package main

import (
	"fmt"
	"github.com/go-muckity/muckity/pkg/muckity"
	"os"
)

func main() {
	var w muckity.WorldSystem
	w = muckity.GetWorld()
	err := w.Init(100)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println("Created GenericWorld:", w.String())
	i, err := w.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(i)
}
