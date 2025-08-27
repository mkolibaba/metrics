package main

import "os"

func main() {
	println("Hello!")
	os.Exit(1) // want "os.Exit is used in main"
}
