package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()

	fmt.Println("time -------", t)
	micro := t.UnixMicro()

	fmt.Println(micro)

	fmt.Println("time from micro", time.UnixMicro(micro))
}
