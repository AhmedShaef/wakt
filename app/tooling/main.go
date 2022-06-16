package main

import (
	"fmt"
	"strings"
)

func main() {
	projctID := "1"
	projectIDs := strings.Split(projctID, ",")
	fmt.Println(projectIDs)
}
