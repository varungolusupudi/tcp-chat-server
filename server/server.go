package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var clients = make(map[net.Conn]string)
var mu sync.Mutex

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, "Enter your username:\n")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	username = username[:len(username)-1]

	mu.Lock()
	clients[conn] = username
	mu.Unlock()

	for {
		msg, error := reader.ReadString('\n')
		if error != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return
		}

		fmt.Printf("Received from %s: %s", conn.RemoteAddr(), msg)

		response := "You said: " + msg

		broadcaster(conn, msg)

		fmt.Fprintf(conn, response)
	}
}

func broadcaster(sender net.Conn, msg string) {
	mu.Lock()
	defer mu.Unlock() // Makes sure the lock is always released even if there's a return or an error. or else could deadlock
	for client := range clients {
		if client != sender {
			time := time.Now().Format("2006-01-02 15:04:05")
			fmt.Fprintf(client, "[%s] %s said: %s", time, clients[sender], msg)
		}
	}
}

func main() {
	fmt.Println("Server Starting")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err) // Prints the error message and stops the program
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err) // Something wrong with the connection
		}
		go handleConnection(conn)
	}
}
