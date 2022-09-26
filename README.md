# dubito

[Dubito](https://it.wikipedia.org/wiki/Dubito) (translatable in English as "I doubt") is an Italian card game based on bluff.

## Rules

This card game has many variants, this repository references the Wikipedia version and adds rules whenever something is not specified.

The game requires at least 3 players and makes use of a [standard 52-card deck](https://en.wikipedia.org/wiki/Standard_52-card_deck) (or a 40-card deck, but this is not the case).

### Preparation

The cards are divided among the players in equal parts. The remaining cards, if any, are discarded and cannot be used.

Since this is a computer adaptation of the game, the turns are chosen randomly.

### How to play (official way)

There are 5 rules:

- Every player can place from 1 to 4 cards on the table.
- The card(s) placed by one player should be one rank above the card(s) placed by the previous player, unless the player bluffs.
- Every card placed on the table is placed with its front facing down, so that nobody else actually knows what card it is.
- In any moment, any player can doubt of the card(s) placed by the last player.
- The winner is the first player to end their pack of cards.

The first player begins by placing card(s) at their choosing on the table from their pack. The same is done by the following players. If any player doubts of the card(s) placed by the last player, those card(s) get uncovered and either one of these two things happens:

- The placed cards are correct. In this case the player who doubted gets to take all the cards placed on the table and the last player plays one more time.
- The player bluffed. In this case the player who bluffed gets to take all the cards placed on the table and the turn jumps to the player who doubted.

### Game example

There are 4 players: A, B, C and D. The turns are: A -> B -> C -> D.

1. A places 3 covered cards and says "three 5".
2. B places 2 covered cards and says "two 6".
3. D decides to doubt of the cards placed by B (not A because it's not the last player).
4. In this case, the cards are actually two 6 so D takes all five cards in the table and the next turn is B again.
5. B places 4 covered cards and says "four 7".
6. C places 1 card and says "one 8".
7. A decides to doubt of the cards placed by C.
8. In this case, the player bluffed, which means that at least one of the cards they placed (one 8) is wrong.
9. The game continues until one player finishes all the cards.

### How to play (fun way)

[todo]

## Building

![Go build status](https://github.com/EdoardoLaGreca/dubito/actions/workflows/go.yml/badge.svg)

Due to the GUI toolkit in use, **before compiling the client** you need to install some dependencies, see [this page](https://developer.fyne.io/started/#prerequisites). These dependencies are not managed by Go and are required to statically link the Go object code with the native GUI libraries, which are platform-dependent.

You can build the whole thing using the commands below (you need to install [Go](https://go.dev/dl/) for this). This will download all the Go module dependencies and build the repository.

```
go build ./cmd/client -o dubito
go build ./cmd/server -o dubito-server
```

### Running without compiling

Thanks to the Go design, it is possible to run the program without compiling it. However, you still need to download the dependencies mentioned above. Also, this may affect performances, which is not really relevant anyway.

```
go run ./cmd/client
go run ./cmd/server
```

## Credits

Thanks to [**mehrasaur**](https://opengameart.org/users/mehrasaur) from [opengameart.org](https://opengameart.org) for card and deck assets. [Download page](https://opengameart.org/content/playing-card-assets-52-cards-deck-chips).

## License

![CC0 logo](https://mirrors.creativecommons.org/presskit/buttons/88x31/svg/cc-zero.svg)

dubito is licensed under the [Creative Commons Zero](https://en.wikipedia.org/wiki/Creative_Commons_license#Zero_/_public_domain) (CC0) license.
