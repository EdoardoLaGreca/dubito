# Source code structure

Many parts could have been done better. In my defense, I did not put all my effort into it. I may make the code cleaner in the future.

This repository follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

This game, speaking of architecture, is made of a client and a server. Clients are used by the players while the server acts as a coordinator and must be hosted somewhere. The message exchange between the client and the server is "request-response"-like, in which only the client can start a request, while the server only listens for client requests and responds to them. This helps keeping the code clear, organized and bug-free.

## Client

All the client code is located in `cmd/client`. It is split into these source files:

 - `main.go`, which contains the minimum code needed to start the program
 - `ui.go`, which handles the user interface
 - `net.go`, which has network-related stuff

In `ui.go`, many functions have `*fyne.Container` as return type, which is where widgets are placed, and `fyne.Window` as one of the parameter types. Those functions can obviously call each other, which is how a window gets its future content. This is usually done while reacting to a user input such as a button click. Notice how these functions keep the code well-divided depending on the context and enable to switch from container to container in an easy and flexible way.

In `net.go`, many functions have `request` at the beginning of their name. The reason for that, as explained above, is that they perform a request to the server and wait for a response. The responses from the server are collected by a coroutine which sends the received messages to the requesting functions through a channel and checks for lost connections without polling, avoiding unnecessary compute overhead.

## Server

All the client code is located in `cmd/server`. It is split into these source files:

 - `main.go`, which contains almost all the code
 - `cli.go`, which handles the command line arguments

In `main.go`, the `handler` function does most of the work, providing a single place to manage all the possible requests.

The possible command line arguments are:

 - `-a [addr]`, which specifies the address to listen to
 - `-p [port]`, which specifies the port to listen to
 - `-m [number]`, which specifies the maximum number of players which can join a game (which is also the minimum number to start the game)

## Internal

The code placed in the `internal` directory is meant to be shared between the client and the server. It usually consists of utility functions made to ease some task.

The internal code is divided into two packages: `cardutils` and `netutils`.

The functions in `cardutils` are related to cards. Those functions are related, although not directly, to network functions since cards are sent as their string representation.

In `netutils`, there are two files: `queue.go` and `utils.go`. The first one manages the message queue while the seconds provides functions for reading and writing strings from and to the connection stream.

The message queue is a buffer for the incoming messages: the `RecvMsg` function fills it with all the incoming messages present in the connection stream and pops the first element of the queue to return it. Then, until the queue will be empty again, it will continue to pop messages from the queue. In this way, it feels like every call to `RecvMsg` reads exactly one string from the connection and returns it, which may be harder and way messier due to corner cases.

Choosing strings as universal encoding for messages passed through the network is an optimal choice: there is no need for specialized fields telling the length of the message, articulated ways of organizing data in the payload or many data structures for representing each kind of message. Instead, at the expense of a little larger payload and the presence of conversion functions, strings enable to know where they end (null bytes), what they contain and they are clearer when debugging.

## Assets

The `assets` directory contains all the assets and a source file (`assets.go`) which embeds them. The main reason for this choice is that it makes it possible to provide a single executable file instead of a huge directory with sub-directories.
