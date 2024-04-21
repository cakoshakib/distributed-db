package main

import (
	"flag"
	"fmt"

	ct "github.com/cakoshakib/distributed-db/client/clienttests"
)

const (
	defaultValue = "value"
)

var (
	testType string
)

func main() {
	flag.StringVar(&testType, "t", defaultValue, "specify test type using go run . -t [speed/value]")
	flag.Parse()

	switch testType {
	case "speed":
		fmt.Println("Running speed test...")
		ct.SpeedTest()
	case "value":
		fmt.Println("Running value test...")
		ct.ValueTest()
	default:
		fmt.Println("client error: no test specified.\nspecify test type using go run . -t [speed/value]")
	}
}
