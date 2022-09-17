package main

import (
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

var connectedPlayers int
var joinedPlayers []*player = make([]*player, 0)
var currentTurn int = 0
var currentRank cardutils.Rank

var lastPlacedCards []cardutils.Card

func getPlayerByConn(conn net.Conn) (int, *player) {
	repeat := false

	if conn == nil {
		return -1, nil
	}

	for {
		for i, p := range joinedPlayers {
			if p == nil {
				// a race condition happened
				repeat = true
				break
			}

			if p.conn.RemoteAddr() == conn.RemoteAddr() {
				return i, p
			}
		}

		if !repeat {
			break
		} else {
			repeat = false
		}
	}

	return -1, nil
}

func fmtPlayerName(p *player) string {
	return p.name + " (" + p.conn.RemoteAddr().String() + ")"
}

// check if the current turn is the player's turn
func checkPlayerTurn(p *player) bool {
	return joinedPlayers[currentTurn].conn.RemoteAddr() == p.conn.RemoteAddr()
}

// returns the winner or nil if the game is not over yet
func checkWin() *player {
	for _, p := range joinedPlayers {
		if len(p.cards) == 0 {
			return p
		}
	}

	return nil
}

// check if player has cards
// the cards should not be duplicated
func checkPlayerHasCards(p *player, cards []cardutils.Card) bool {
	cardsFound := make(map[cardutils.Card]bool)

	// add cards
	for _, c := range cards {
		cardsFound[c] = false
	}

	for _, c := range cards {
		for _, pc := range p.cards {
			if c == pc {
				cardsFound[c] = true
			}
		}
	}

	for _, found := range cardsFound {
		if !found {
			return false
		}
	}

	return true
}

// return true if all the cards match the rank
func checkLastPlacedCards(rank cardutils.Rank) bool {
	for _, c := range lastPlacedCards {
		if c.Rank != rank {
			return false
		}
	}

	return true
}

func handler(conn net.Conn, maxPlayers int, addPlayer chan<- *player, removePlayer chan<- *player) {
	log.Println("a player connected (IP: " + conn.RemoteAddr().String() + ")")
	hasJoined := false
	var p *player     // read this only if the player has joined
	var indexInJP int // index in joinedPlayers

	// remove player when handler ends
	defer func() { connectedPlayers--; removePlayer <- p }()

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
					netutils.SendMsg(conn, "the game is full")
					break msgLoop
				} else {
					tmpPlayer := new(player)
					tmpPlayer.conn = conn
					tmpPlayer.name = fields[1]
					addPlayer <- tmpPlayer

					// wait for the player to be added to joinedPlayers
					for {
						index, foundPlayer := getPlayerByConn(conn)
						if foundPlayer == nil {
							continue
						}

						p = foundPlayer
						indexInJP = index
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
					for len(p.cards) == 0 {
						// wait to get cards
						time.Sleep(time.Millisecond * 100)
					}

					var cardsStr string
					for _, c := range p.cards {
						cardsStr += cardutils.CardToString(c) + ","
					}
					cardsStr = strings.TrimSuffix(cardsStr, ",")
					netutils.SendMsg(conn, cardsStr)

				case "my-turn":
					winner := checkWin()
					if winner == nil {
						if checkPlayerTurn(p) {
							netutils.SendMsg(conn, "yes")
						} else {
							netutils.SendMsg(conn, "no")
						}
					} else {
						if p == winner {
							netutils.SendMsg(conn, "winner")
						} else {
							netutils.SendMsg(conn, "loser")
						}
					}
				default:
					log.Println("invalid request from " + fmtPlayerName(p) + ": \"" + msg + "\"")
				}
			}
		case "place":
			if hasJoined {
				if !checkPlayerTurn(p) {
					netutils.SendMsg(conn, "wrong turn")
				} else {
					cardsStr := strings.Split(strings.Join(fields[1:], " "), ",")
					cards := make([]cardutils.Card, len(cardsStr))

					// convert cardsStr into cards
					for i, cs := range cardsStr {
						card, err := cardutils.CardByName(cs)
						if err != nil {
							log.Println("invalid card placed: \"" + cs + "\"")
							break
						}
						cards[i] = card
					}

					// check number of cards
					if len(cards) >= 1 && len(cards) <= 4 {
						// check if the player have those cards
						if checkPlayerHasCards(p, cards) {
							// place the cards
							lastPlacedCards = cards
							netutils.SendMsg(conn, "ok")
							currentTurn++

							if currentTurn >= len(joinedPlayers) {
								currentTurn = 0
							}
						} else {
							netutils.SendMsg(conn, "you don't have that card")
						}
					} else {
						netutils.SendMsg(conn, "too many cards")
					}
				}
			}
		case "dubito":
			if hasJoined {
				// send "right" if last player lied, "wrong" otherwise
				if checkLastPlacedCards(currentRank) {
					// last player didn't lie
					netutils.SendMsg(conn, "wrong")

					// repeat the turn for the last player
					currentTurn--
				} else {
					// last player lied
					netutils.SendMsg(conn, "right")
					currentTurn = indexInJP
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

	addPlayerChan := make(chan *player)
	removePlayerChan := make(chan *player)

	// goroutine to add and remove players from joinedPlayers
	go func(addPlayers <-chan *player, removePlayers <-chan *player) {
		for {
			select {
			case p := <-addPlayers:
				joinedPlayers = append(joinedPlayers, p)

			case p := <-removePlayers:
				// find and remove player
				for i, jp := range joinedPlayers {
					if jp == p {
						joinedPlayers = append(joinedPlayers[:i], joinedPlayers[i+1:]...)
						break
					}
				}
			}
		}
	}(addPlayerChan, removePlayerChan)

	for {
		log.Println("waiting for players to join...")

		var wg sync.WaitGroup

		// let players connect
		for len(joinedPlayers) < maxPlayers {
			if connectedPlayers < maxPlayers {
				conn, err := lis.Accept()
				if err != nil {
					log.Println(err.Error())
				}

				connectedPlayers++

				wg.Add(1)
				go func() {
					defer wg.Done()
					handler(conn, maxPlayers, addPlayerChan, removePlayerChan)
				}()
			} else {
				time.Sleep(200 * time.Millisecond)
			}
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
}
