package clienttests

import (
	"fmt"
	"time"
)

func SpeedTest() {
	user := "user1"
	table := "table1"

	ProcessRequest(fmt.Sprintf("cu %s;\n", user), leader)
	ProcessRequest(fmt.Sprintf("ct %s %s;\n", user, table), leader)

	//Begin test
	n := 1000

	start := time.Now()
	for i := 1; i <= n; i++ {
		request := fmt.Sprintf("add %s %s test%d value%d;\n", user, table, i, i)
		ProcessRequest(request, leader)
	}
	checkpoint := time.Now()
	for i := 1; i <= n; i++ {
		request := fmt.Sprintf("get %s %s test%d;\n", user, table, i)
		ProcessRequest(request, leader)
	}
	end := time.Now()

	elapsed_write := checkpoint.Sub(start)
	elapsed_read := end.Sub(checkpoint)
	fmt.Println("Time taken for writes:", elapsed_write)
	fmt.Println("Time taken for reads:", elapsed_read)
}
