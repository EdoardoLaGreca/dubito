package main

import (
	"net"
	"strconv"
	"strings"
)

var conn net.Conn

func initConn() error {
	c, err := openConn(serverAddress, serverPort)
	conn = c

	return err
}

func openConn(addr string, port uint16) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr+":"+string(port))
	return conn, err
}

func sendMsg(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg))

	return err
}

func recvMsg(conn net.Conn) (string, error) {
	msg := make([]byte, 0)
	_, err := conn.Read(msg)

	return string(msg), err
}

func closeConn(conn net.Conn) error {
	return conn.Close()
}

func requestPlayers(conn net.Conn) ([]string, error) {
	sendMsg(conn, "get players")

	playersCsv, err := recvMsg(conn)
	if err != nil {
		return nil, err
	}

	players := strings.Split(playersCsv, ",")
	return players, nil
}

func requestMaxPlayers(conn net.Conn) (uint, error) {
	sendMsg(conn, "get max-players")

	maxPlayersStr, err := recvMsg(conn)
	if err != nil {
		return 0, err
	}

	maxPlayers, err := strconv.Atoi(maxPlayersStr)
	if err != nil {
		return 0, err
	}

	return uint(maxPlayers), nil
}
