package main

import (
	"net"
	"strconv"
	"strings"

	"github.com/EdoardoLaGreca/dubito/internal/cardutils"
	"github.com/EdoardoLaGreca/dubito/internal/netutils"
)

var serverAddress string = "localhost"
var serverPort uint16 = 9876

var conn net.Conn

// do not use openConn if you call initConn
func openConn(addr string, port uint16) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr+":"+strconv.Itoa(int(port)))
	return conn, err
}

func initConn() error {
	c, err := openConn(serverAddress, serverPort)
	if err != nil {
		return err
	}

	conn = c

	return nil
}

func requestJoin(conn net.Conn) error {
	err := netutils.SendMsg(conn, "join "+username)
	if err != nil {
		return err
	}

	return nil
}

func requestPlayers(conn net.Conn) ([]string, error) {
	err := netutils.SendMsg(conn, "get players")
	if err != nil {
		return nil, err
	}

	playersCsv, err := netutils.RecvMsg(conn)
	if err != nil {
		return nil, err
	}

	players := strings.Split(playersCsv, ",")
	return players, nil
}

func requestMaxPlayers(conn net.Conn) (uint, error) {
	err := netutils.SendMsg(conn, "get max-players")
	if err != nil {
		return 0, err
	}

	maxPlayersStr, err := netutils.RecvMsg(conn)
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
	err := netutils.SendMsg(conn, "get cards")
	if err != nil {
		return nil, err
	}

	cardsStr, err := netutils.RecvMsg(conn)
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

func requestLeave() error {
	return netutils.SendMsg(conn, "leave")
}
