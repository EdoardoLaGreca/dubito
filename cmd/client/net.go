package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/EdoardoLaGreca/dubito/internal/cardutils"
	"github.com/EdoardoLaGreca/dubito/internal/netutils"
)

type netResponse struct {
	msg string
	err error
}

var serverAddress string = "localhost"
var serverPort uint16 = 9876

var conn net.Conn

var recvChan chan netResponse = make(chan netResponse) // receive messages
var closeChan chan struct{} = make(chan struct{})      // the connection closed
var stopCheckChan chan struct{} = make(chan struct{})  // stop receiving for messages

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

	// goroutine to check for crashes and receive messages
	go func(conn net.Conn) {
		for {
			select {
			case <-stopCheckChan:
				return
			default:
				resp, err := netutils.RecvMsg(conn)
				if err != nil && err == io.EOF {
					closeChan <- struct{}{}
					return
				} else {
					recvChan <- netResponse{msg: resp, err: err}
				}
			}
		}
	}(conn)

	return nil
}

func requestJoin() error {
	err := netutils.SendMsg(conn, "join "+username)
	if err != nil {
		return err
	}

	resp := <-recvChan
	if resp.err != nil {
		return resp.err
	}

	if resp.msg != "ok" {
		return fmt.Errorf(resp.msg)
	}

	return nil
}

func requestPlayers() ([]string, error) {
	err := netutils.SendMsg(conn, "get players")
	if err != nil {
		return nil, err
	}

	playersCsv := <-recvChan
	if playersCsv.err != nil {
		return nil, playersCsv.err
	}

	players := strings.Split(playersCsv.msg, ",")
	return players, nil
}

func requestMaxPlayers() (uint, error) {
	err := netutils.SendMsg(conn, "get max-players")
	if err != nil {
		return 0, err
	}

	maxPlayersStr := <-recvChan
	if maxPlayersStr.err != nil {
		return 0, maxPlayersStr.err
	}

	maxPlayers, err := strconv.Atoi(maxPlayersStr.msg)
	if err != nil {
		return 0, err
	}

	return uint(maxPlayers), nil
}

func requestCards() ([]cardutils.Card, error) {
	err := netutils.SendMsg(conn, "get cards")
	if err != nil {
		return nil, err
	}

	cardsStr := <-recvChan
	if cardsStr.err != nil {
		return nil, cardsStr.err
	}

	cardsSp := strings.Split(cardsStr.msg, ",")

	cards := make([]cardutils.Card, len(cardsSp))

	for i := range cardsSp {
		cards[i], err = cardutils.CardByName(cardsSp[i])
		if err != nil {
			return nil, err
		}
	}

	return cards, nil
}

func requestTurn() (bool, error) {
	err := netutils.SendMsg(conn, "get my-turn")
	if err != nil {
		return false, err
	}

	resp := <-recvChan
	if resp.err != nil {
		return false, resp.err
	}

	if resp.msg != "yes" {
		return false, nil
	}

	return true, nil
}

func requestPlaceCards(cards []cardutils.Card) (bool, error) {
	var cardsStr string

	for _, c := range cards {
		cardsStr += cardutils.CardToString(c) + ","
	}

	strings.TrimRight(cardsStr, ",")

	err := netutils.SendMsg(conn, "place "+cardsStr)
	if err != nil {
		return false, err
	}

	resp := <-recvChan
	if resp.err != nil {
		return false, resp.err
	}

	if resp.msg != "ok" {
		return false, nil
	}

	return true, nil
}

func requestLeave() error {
	err := netutils.SendMsg(conn, "leave")
	conn.Close()

	return err
}
