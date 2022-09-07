package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/EdoardoLaGreca/dubito/internal/cardutils"
	"github.com/EdoardoLaGreca/dubito/internal/netutils"
)

type player struct {
	conn  net.Conn
	name  string
	cards []cardutils.Card
}

var joinedPlayers []player = make([]player, 0)

func getPlayerByConn(conn net.Conn) (*player, error) {
	for _, p := range joinedPlayers {
		if p.conn == conn {
			return &p, nil
		}
	}

	return &player{}, fmt.Errorf("player not found")
}

func fmtPlayerName(p *player) string {
	return p.name + " (" + p.conn.RemoteAddr().String() + ")"
}

func handler(conn net.Conn, maxPlayers int, players chan<- player) {
	log.Println("a player connected (IP: " + conn.RemoteAddr().String() + ")")
	hasJoined := false
	var p *player

	for {
		msg, err := netutils.RecvMsg(conn)
		if err != nil {
			log.Println(err.Error())
		}

		fields := strings.Fields(msg)

		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "join":
			if !hasJoined {
				players <- player{conn: conn, name: fields[1]}

				// wait for the player to be added to joinedPlayers
				for {
					foundPlayer, err := getPlayerByConn(conn)
					if err != nil {
						continue
					}

					p = foundPlayer
					hasJoined = true
					break
				}
				log.Println("player " + fmtPlayerName(p) + " joined")
			}
		case "get":
			if hasJoined {
				switch fields[1] {
				case "players":
					var jpStr string
					for _, p := range joinedPlayers {
						jpStr += p.name + ","
					}
					jpStr = strings.TrimSuffix(jpStr, ",")
					netutils.SendMsg(conn, jpStr)

				case "max-players":
					netutils.SendMsg(conn, strconv.Itoa(maxPlayers))

				case "cards":
					var cardsStr string
					for _, c := range p.cards {
						cardsStr += cardutils.CardToString(c) + ","
					}
					cardsStr = strings.TrimSuffix(cardsStr, ",")
					netutils.SendMsg(conn, cardsStr)

				default:
					log.Println("invalid request from " + fmtPlayerName(p) + ": \"" + msg + "\"")
				}
			}
		default:
			log.Println("invalid request from " + fmtPlayerName(p) + ": \"" + msg + "\"")
		}
	}
}

func giveCards(playersCount int) [][]cardutils.Card {
	cardsEach := 52 / playersCount
	givenCards := make(map[cardutils.Card]bool)
	cards := make([][]cardutils.Card, playersCount)

	for i := range cards {
		playerCards := make([]cardutils.Card, cardsEach)

		for j := range playerCards {
			for {
				suit := rand.Intn(4)
				rank := rand.Intn(13) + 1
				card := cardutils.Card{Suit: cardutils.Suit(suit), Rank: cardutils.Rank(rank)}

				if !givenCards[card] {
					playerCards[j] = card
					givenCards[card] = true
				}
			}
		}

		cards[i] = playerCards
	}

	return cards
}

func main() {
	lisAddr, err := getListenAddress()
	if err != nil {
		panic(err.Error())
	}

	lisPort, err := getListenPort()
	if err != nil {
		panic(err.Error())
	}

	lis, err := net.Listen("tcp", lisAddr+":"+strconv.Itoa(int(lisPort)))
	if err != nil {
		panic(err.Error())
	}

	maxPlayers, err := getArgMaxPlayers()
	if err != nil {
		panic(err.Error())
	}

	playersChan := make(chan player)

	var wg sync.WaitGroup

	// create goroutine to append new players to joinedPlayers
	go func(players <-chan player) {
		// add joined players
		for i := 0; i < maxPlayers; i++ {
			joinedPlayers = append(joinedPlayers, <-players)
		}
	}(playersChan)

	// let players connect
	for i := 0; i < maxPlayers; i++ {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err.Error())
		}

		wg.Add(1)
		go func() {
			handler(conn, maxPlayers, playersChan)
			wg.Done()
		}()
	}

	// give cards
	cards := giveCards(len(joinedPlayers))
	for i := range joinedPlayers {
		joinedPlayers[i].cards = cards[i]
	}

	wg.Wait()
}
