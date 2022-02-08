package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	fmt.Println("What's your name?")

	var name string
	fmt.Fscan(os.Stdin, &name)

	fmt.Println("What's your age?")

	var age int
	_, err := fmt.Scanf("%d", &age)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello, %s %d years old!\n", name, age)

}
