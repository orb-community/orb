package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("halau")
		time.Sleep(time.Duration(5) * time.Second)
	}
}
