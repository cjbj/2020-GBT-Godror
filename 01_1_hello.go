package main

/*

  go run 01_1_hello.go

  go build 01_1_hello.go
  ./hello

  GOOS=windows go build 01_1_hello.go

*/

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}
