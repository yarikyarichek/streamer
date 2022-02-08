package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Println("What's your name?")

	var name string
	fmt.Fscan(os.Stdin, &name)

	fmt.Println("Hello,", name)
}
