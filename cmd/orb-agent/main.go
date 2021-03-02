package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("orb-agent-checkpoint")
		time.Sleep(time.Duration(30) * time.Second)
	}
}
