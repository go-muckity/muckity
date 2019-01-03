package main

import (
	"fmt"
	"github.com/tsal/muckity/pkg/muckity"
)

func main() {
	muckity.GetMuckityStorage()
	w, err := muckity.NewWorld("Development World")
	if err != nil {
		panic(err)
	}
	fmt.Println("ID\t: ", w.Name())
}
