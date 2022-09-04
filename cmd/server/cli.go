package main

import (
	"fmt"
	"os"
	"strconv"
)

// return the specified argument position, -1 if it could not be found
func getArgPos(argname string) (pos int) {
	pos = -1

	for i, a := range os.Args {
		if a == argname {
			pos = i
		}
	}

	return 0, fmt.Errorf("no -m arg found")
}
