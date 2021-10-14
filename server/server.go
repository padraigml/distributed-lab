package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error, conn net.Conn) {
	if err.Error() == "EOF" {
		fmt.Println("okay mate")
		conn.Close()
	}

}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, e := ln.Accept()
		if e != nil {
			handleError(e, conn)
		}
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	reader := bufio.NewReader(client)
	for {
		msg, e := reader.ReadString('\n')
		if e != nil {
			handleError(e, client)
		}
		m := Message{sender: clientid, message: msg}

		msgs <- m
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	ln, _ := net.Listen("tcp", *portPtr)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			id := len(clients)
			clients[id] = conn
			go handleClient(conn, id, msgs)
		case msg := <-msgs:
			for i := 0; i < len(clients); i++ {
				if clients[i] != clients[msg.sender] {
					fmt.Fprintf(clients[i], msg.message)
				}
			}
		}
	}
}
