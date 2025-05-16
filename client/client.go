package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// Steps for client
// Step1: Log that the client starter
// Step2: Dial to the port that the server is listening on
// Step3: Error handling if dialing failed
// Step4: Keep sending messages to the server and printing the responses

func main() {
	fmt.Println("Client Started")

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)
	res, _ := reader.ReadString('\n')
	fmt.Println(res)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		msg := scanner.Text() + "\n"
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		for {
			response, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Disconnected from server:", err)
				os.Exit(0)
			}

			fmt.Println(response)
		}
	}()

	fmt.Println("Enter your message: \n")
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > 0 && text[0] == '/' {
			switch text {
			case "/quit":
				fmt.Println("Goodbye!")
				conn.Close()
				os.Exit(0)
			case "/help":
				fmt.Println("Available commands: /quit, /help, /users")
			default:
				fmt.Println("Invalid command")
			}
			continue
		}

		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			log.Println("Error writing to server", err)
			break
		}

	}
}
