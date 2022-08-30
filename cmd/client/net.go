package main

import (
	"net"
	"strconv"
	"strings"

	"github.com/EdoardoLaGreca/dubito/internal/cardutils"
)

var conn net.Conn

func openConn(addr string, port uint16) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr+":"+string(port))
	return conn, err
}

func initConn() error {
	c, err := openConn(serverAddress, serverPort)
	conn = c

	return err
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
	err := sendMsg(conn, "get players")
	if err != nil {
		return nil, err
	}

	playersCsv, err := recvMsg(conn)
	if err != nil {
		return nil, err
	}

	players := strings.Split(playersCsv, ",")
	return players, nil
}

func requestMaxPlayers(conn net.Conn) (uint, error) {
	err := sendMsg(conn, "get max-players")
	if err != nil {
		return 0, err
	}

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

func requestCards(conn net.Conn) ([]cardutils.Card, error) {
	err := sendMsg(conn, "get cards")
	if err != nil {
		return nil, err
	}

	cardsStr, err := recvMsg(conn)
	if err != nil {
		return nil, err
	}

	cardsSp := strings.Split(cardsStr, ",")

	cards := make([]cardutils.Card, len(cardsStr))

	for i := range cardsSp {
		cards[i], err = cardutils.CardByName(cardsSp[i])
		if err != nil {
			return nil, err
		}
	}

	return cards, nil
}
