package game

import (
	"fmt"
	"github.com/tsal/muckity/pkg/muckity"
)

func main() {
	w, err := muckity.NewWorld("A Brand New World")
	if err != nil {
		panic(err)
	}
	fmt.Println("ID\t: ", w.Name())
}
