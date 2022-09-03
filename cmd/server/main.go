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
)

type player struct {
	ip    string
	name  string
	cards []cardutils.Card
}

var joinedPlayers []player = make([]player, 0)

func getPlayerByIP(ip net.Addr) (*player, error) {
	for _, p := range joinedPlayers {
		if p.ip == ip.String() {
			return &p, nil
		}
	}

	return &player{}, fmt.Errorf("player not found")
}

func handler(conn net.Conn, maxPlayers int, players chan player) {
	b := make([]byte, 0)
	hasJoined := false
	ip := conn.RemoteAddr()
	var p *player

	for {
		_, err := conn.Read(b)
		if err != nil {
			log.Println(err.Error())
		}

		msg := string(b)
		fields := strings.Fields(msg)

		switch fields[0] {
		case "join":
			if !hasJoined {
				players <- player{ip: ip.String(), name: fields[1]}
				for player, err := getPlayerByIP(ip); err != nil; {
					p = player
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
					conn.Write([]byte(jpStr))

				case "max-players":
					conn.Write([]byte(strconv.Itoa(maxPlayers)))

				case "cards":
					var cardsStr string
					for _, c := range p.cards {
						cardsStr += cardutils.CardToString(c) + ","
					}
					cardsStr = strings.TrimSuffix(cardsStr, ",")
					conn.Write([]byte(cardsStr))

				default:
					log.Println("unknown request: \"" + msg + "\"")
				}
			}
		default:
			log.Println("unknown request: \"" + msg + "\"")
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
	lis, err := net.Listen("tcp", "0.0.0.0")
	if err != nil {
		panic(err.Error())
	}

	maxPlayers, err := getArgMaxPlayers()
	if err != nil {
		panic(err.Error())
	}

	playersChan := make(chan player)

	var wg sync.WaitGroup

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

	// add joined players
	for i := 0; i < maxPlayers; i++ {
		joinedPlayers = append(joinedPlayers, <-playersChan)
	}

	// give cards
	cards := giveCards(len(joinedPlayers))
	for i := range joinedPlayers {
		joinedPlayers[i].cards = cards[i]
	}

	wg.Wait()
}
