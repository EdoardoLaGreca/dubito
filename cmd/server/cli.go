package main

import (
	"os"
	"strconv"
	"strings"
)

func getArgMaxPlayers() (int, error) {
	for i, arg := range os.Args {
		if arg == "-m" || (arg[0] == '-' && arg[1] != '-' && strings.Contains(arg, "m")) {
			maxPlayers, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				return 0, err
			}

			return maxPlayers, nil
		}
	}
}
