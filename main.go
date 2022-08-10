package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	t2 := t.Add(-5 * time.Hour)

	fmt.Println(t, t2)
}
