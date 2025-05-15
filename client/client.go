package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Client Starting")
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// Enter username

	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	scanner1 := bufio.NewScanner(os.Stdin)
	if scanner1.Scan() {
		msg := scanner1.Text() + "\n"
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Listen to server messages
	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Disconnected from Server")
				return
			}
			fmt.Println(msg)

		}
	}()

	// Send messages to server
	fmt.Println("Enter your message: \n")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			msg := scanner.Text() + "\n"
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Println("Error sending the message", err)
				break
			}
		}

	}
}
