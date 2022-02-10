# go-mongoqb


How to Use

```go

package main

import (
    "github.com/gokultp/go-mongoqb"
    "fmt"
)

func main() {
    qb := mongoqb.NewQueryBuilder().
            Eq("name", "gokul").
            Gt("age", 27)
    fmt.Println(qb.Build())
}

```