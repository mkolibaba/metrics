package main

import (
	. "os"
)

func main() {
	Exit(1) // want "os.Exit is used in main"
}
