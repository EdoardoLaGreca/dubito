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

	return
}

func getArgMaxPlayers() (int, error) {
	pos := getArgPos("-m")

	if pos == -1 {
		return 0, fmt.Errorf("no -m arg found")
	}

	maxPlayers, err := strconv.Atoi(os.Args[pos+1])
	if err != nil {
		return 0, err
	}

	return maxPlayers, nil
}

func getListenAddress() (string, error) {
	pos := getArgPos("-a")

	if pos == -1 {
		return "", fmt.Errorf("no -a arg found")
	}

	return os.Args[pos+1], nil
}

func getListenPort() (uint16, error) {
	pos := getArgPos("-p")

	if pos == -1 {
		return 0, fmt.Errorf("no -p arg found")
	}

	port, err := strconv.Atoi(os.Args[pos+1])
	if err != nil {
		return 0, err
	}

	return uint16(port), nil
}
