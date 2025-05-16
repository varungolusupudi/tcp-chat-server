package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
}

const reset = "\033[0m"

var clients = make(map[net.Conn]string)
var clientColors = make(map[net.Conn]string)
var mu sync.Mutex

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	fmt.Fprintf(conn, "Enter your username: \n")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	username = username[:len(username)-1]
	color := colors[rand.Intn(len(colors))]

	mu.Lock()
	clients[conn] = username
	clientColors[conn] = color
	mu.Unlock()

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Client Disconnected", conn.RemoteAddr())
			broadCaster(conn, clients[conn]+" disconnected")
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return
		}

		fmt.Println("Client %s said %s", conn.RemoteAddr(), msg)
		response := "You said: " + msg

		broadCaster(conn, msg)
		fmt.Fprintf(conn, response)
	}
}

func broadCaster(sender net.Conn, msg string) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		if client != sender {
			fmt.Fprintf(client, "%s%s%s", clientColors[sender], clients[sender], reset)
			fmt.Fprintf(client, ": %s", msg)
		}
	}
}

func handleInterrupt() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received terminate signal")
		mu.Lock()
		for conn := range clients {
			conn.Write([]byte("Server is shutting down...\n"))
			conn.Close()
		}
		mu.Unlock()
		os.Exit(0)
	}()
}

// Steps for the server
// Step 1: Print a statement indicating that the server starting
// Step 2: Start the server using the net Listen function to listen to connections on a port(eg: 8080) and tcp protocol
// Step 3: Error handling for listening to connections
// Step 4: Accept connections on the port
// Step 5: Have a go routine that handles each connection and send a response back for ack
func main() {
	fmt.Println("Server starting")
	handleInterrupt()
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}

}
