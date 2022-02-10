package main

import (
	"fmt"
	"time"
)

func main() {
	d := time.Duration(1328519932000000) / time.Millisecond

	fmt.Println(d)
}
