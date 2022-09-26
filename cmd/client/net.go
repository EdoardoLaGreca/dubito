package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/EdoardoLaGreca/dubito/internal/cardutils"
	"github.com/EdoardoLaGreca/dubito/internal/netutils"
)

type netResponse struct {
	msg string
	err error
}

// response to "get update" request
type update struct {
	gameOver    bool
	playerWon   bool           // not relevant if gameOver = false
	playerTurn  bool           // not relevant if gameOver = true
	cardsAmount int            // not relevant if gameOver = true
	cardRank    cardutils.Rank // not relevant if gameOver = true
}

var errWinner error = fmt.Errorf("winner")
var errLoser error = fmt.Errorf("loser")

var serverAddress string = "localhost"
var serverPort uint16 = 9876

var conn net.Conn

var recvChan chan netResponse = make(chan netResponse) // receive messages
var closeChan chan struct{} = make(chan struct{})      // the connection closed
var stopCheckChan chan struct{} = make(chan struct{})  // stop receiving for messages

var netMutex sync.Mutex

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
	netMutex.Lock()
	defer netMutex.Unlock()

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
	netMutex.Lock()
	defer netMutex.Unlock()

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
	netMutex.Lock()
	defer netMutex.Unlock()

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
	netMutex.Lock()
	defer netMutex.Unlock()

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

func StrToUpdate(response string) (update, error) {
	respLines := strings.Split(response, " ")

	ud := update{}

	// first line
	switch respLines[0] {
	case "y":
		ud.gameOver = true
		ud.playerWon = true
	case "n":
		ud.gameOver = true
		ud.playerWon = false
	case "u":
		ud.gameOver = false
	default:
		return update{}, fmt.Errorf("invalid line 0")
	}

	if ud.gameOver {
		return ud, nil
	}

	// second line
	switch respLines[1] {
	case "y":
		ud.playerTurn = true
	case "n":
		ud.playerTurn = false
	default:
		return update{}, fmt.Errorf("invalid line 1")
	}

	// third line
	cards := strings.Fields(respLines[2])

	cardsAmount, err := strconv.Atoi(cards[0])
	if err != nil {
		return update{}, err
	}
	ud.cardsAmount = cardsAmount

	cardRank, err := cardutils.RankByName(cards[1])
	if err != nil {
		return update{}, nil
	}
	ud.cardRank = cardRank

	return ud, nil
}

// the response message is structured as follows:
// [y/n/u if the player won/lost or the game is not over yet]\n
// [y/n if the current turn is the player's turn]\n
// [cards which the last player said to have placed (<N> <card rank>, e.g. "3 seven")]
//
// e.g.
// "u\n
// n\n
// 2 ace"
func requestUpdate() (update, error) {
	netMutex.Lock()
	defer netMutex.Unlock()

	err := netutils.SendMsg(conn, "get update")
	if err != nil {
		return update{}, err
	}

	resp := <-recvChan
	if resp.err != nil {
		return update{}, resp.err
	}

	return StrToUpdate(resp.msg)
}

func requestPlaceCards(cards []cardutils.Card) (bool, error) {
	netMutex.Lock()
	defer netMutex.Unlock()

	var cardsStr string

	for _, c := range cards {
		cardsStr += cardutils.CardToString(c) + ","
	}

	cardsStr = strings.TrimRight(cardsStr, ",")

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

// return true if the doubt was correct (last player lied)
func requestDubito() (bool, error) {
	netMutex.Lock()
	defer netMutex.Unlock()

	err := netutils.SendMsg(conn, "dubito")
	if err != nil {
		return false, err
	}

	resp := <-recvChan
	if resp.err != nil {
		return false, resp.err
	}

	if resp.msg == "right" {
		return true, nil
	}

	return false, nil
}

func requestLeave() error {
	netMutex.Lock()
	defer netMutex.Unlock()

	err := netutils.SendMsg(conn, "leave")
	conn.Close()

	return err
}
