package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/EdoardoLaGreca/dubito/internal/cardutils"
	"github.com/EdoardoLaGreca/dubito/internal/netutils"
)

type player struct {
	conn  net.Conn
	name  string
	cards []cardutils.Card
}

var joinedPlayers []*player = make([]*player, 0)

func getPlayerByConn(conn net.Conn) (*player, error) {
	for _, p := range joinedPlayers {
		if p.conn == conn {
			return p, nil
		}
	}

	return &player{}, fmt.Errorf("player not found")
}

func fmtPlayerName(p *player) string {
	return p.name + " (" + p.conn.RemoteAddr().String() + ")"
}

func handler(conn net.Conn, maxPlayers int, players chan<- *player) {
	log.Println("a player connected (IP: " + conn.RemoteAddr().String() + ")")
	hasJoined := false
	var p *player // read this only if the player has joined

msgLoop:
	for {
		msg, err := netutils.RecvMsg(conn)
		if err != nil {
			if err == io.EOF {
				// connection closed
				log.Println("the connection to " + conn.RemoteAddr().String() + " has been closed")
			} else {
				log.Println("an error occurred while reading a message: " + err.Error())
			}
			break
		}

		log.Println(conn.RemoteAddr().String() + " made a request: \"" + msg + "\"")

		fields := strings.Fields(msg)

		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "join":
			if !hasJoined {
				if len(joinedPlayers) == maxPlayers {
					netutils.SendMsg(conn, "the lobby is full")
					break msgLoop
				} else {
					tmpPlayer := new(player)
					tmpPlayer.conn = conn
					tmpPlayer.name = fields[1]
					players <- tmpPlayer

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
					netutils.SendMsg(conn, "ok")

					log.Println("player " + fmtPlayerName(p) + " joined")
				}
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
		case "leave":
			if hasJoined {
				log.Println("player " + fmtPlayerName(p) + " left")
				break msgLoop
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
					break
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

	playersChan := make(chan *player)
	maxPlayersJoinedChan := make(chan struct{})

	var wg sync.WaitGroup

	// create goroutine to append new players to joinedPlayers
	go func(players <-chan *player) {
		// add joined players
		for i := 0; i < maxPlayers; i++ {
			joinedPlayers = append(joinedPlayers, <-players)
		}
	}(playersChan)

	log.Println("waiting for all the players to join...")

	// let players connect
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Println(err.Error())
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				handler(conn, maxPlayers, playersChan)
			}()
		}
	}()

	// check if all the players joined
	go func() {
		for len(joinedPlayers) < maxPlayers {
			time.Sleep(time.Millisecond * 100)
		}
		<-maxPlayersJoinedChan
	}()

	<-maxPlayersJoinedChan

	// wait for all the players to join
	for len(joinedPlayers) < maxPlayers {
		time.Sleep(time.Millisecond * 100)
	}
	log.Println("all the players joined the game")

	// give cards
	cards := giveCards(len(joinedPlayers))
	for i := range joinedPlayers {
		joinedPlayers[i].cards = cards[i]
		log.Println("cards have been assigned to " + joinedPlayers[i].name)
	}

	wg.Wait()
}
