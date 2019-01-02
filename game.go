package main

import (
	"fmt"
	"github.com/tsal/muckity/pkg/muckity"
)

func main() {
	url := "mongodb://muckity:muckity@mongo:27017/muckity"
	fmt.Println("Attempting to grab url: ", url)
	storage := muckity.NewMuckityStorage(url)
	fmt.Println(storage.Client.Database("muckity"))
	w, err := muckity.NewWorld("A Brand New World")
	if err != nil {
		panic(err)
	}
	fmt.Println("ID\t: ", w.Name())
}
