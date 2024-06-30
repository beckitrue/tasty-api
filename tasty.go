package main

import (
	"cmd/tastymenu"
	"fmt"
	"os"
)

func main() {

	fmt.Println("Welcome to Tasty!")
	tastymenu.Run(os.Args)
}
