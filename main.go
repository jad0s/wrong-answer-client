package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jad0s/wrong-answer-client/internal/config"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type RevealImpostorPayload struct {
	Impostor         string `json:"impostor"`
	MostVoted        string `json:"most_voted"`
	ImpostorQuestion string `json:"impostor_question"`
}

func main() {
	fmt.Println("config file path: ", config.ConfigPath)
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	if err := sendJSON(conn, "set_username", username); err != nil {
		log.Fatal("Set username error:", err)
	}
	fmt.Println("Connected! Waiting for game...")

	inputChan := make(chan bool, 1)
	voteChan := make(chan bool, 1)

	go func() {
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Read error:", err)
				os.Exit(0)
			}
			//DEBUG
			//fmt.Printf("DEBUG: Received message type: '%s', payload: '%s'\n", msg.Type, string(msg.Payload))

			switch msg.Type {
			case "question":
				var q string
				_ = json.Unmarshal(msg.Payload, &q)
				fmt.Println("\nQuestion:", q)
				fmt.Print("Your answer: ")
				inputChan <- true

			case "reveal_answers":
				var txt string
				_ = json.Unmarshal(msg.Payload, &txt)
				fmt.Println("\n--- All Answers ---")
				fmt.Println(txt)
				fmt.Println("-------------------")

			case "vote":
				var open string
				_ = json.Unmarshal(msg.Payload, &open)
				fmt.Println("\nVoting started! Enter the username of the player you think is the impostor:")
				voteChan <- true

			case "reveal_impostor":
				var reveal RevealImpostorPayload
				if err := json.Unmarshal(msg.Payload, &reveal); err != nil {
					fmt.Println("Error decoding impostor reveal payload:", err)
					break
				}
				fmt.Println("\n--- Voting Result ---")
				if reveal.Impostor == reveal.MostVoted {
					fmt.Printf("You guessed correctly! The impostor was: %s\n", reveal.Impostor)
				} else {
					fmt.Printf("Wrong guess! %s was not the impostor. The impostor was: %s\n", reveal.MostVoted, reveal.Impostor)
				}
				fmt.Printf("Impostor's question was: %s\n", reveal.ImpostorQuestion)
				fmt.Println("----------------------")

			default:
				fmt.Println("Unknown message:", msg)
			}
		}
	}()

	for {
		select {
		case <-inputChan:
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == "" {
				fmt.Print("Please enter a non-empty answer: ")
				inputChan <- true
				continue
			}
			if err := sendJSON(conn, "submit_answer", answer); err != nil {
				log.Println("Write error:", err)
				return
			}
			fmt.Println("Answer submitted, waiting for others...")

		case <-voteChan:
			vote, _ := reader.ReadString('\n')
			vote = strings.TrimSpace(vote)
			if vote == "" {
				fmt.Print("Please enter a username to vote for: ")
				voteChan <- true
				continue
			}
			if err := sendJSON(conn, "vote", vote); err != nil {
				log.Println("Write error:", err)
				return
			}
			fmt.Println("Vote submitted, waiting for others...")
		}
	}
}

func sendJSON(conn *websocket.Conn, typ string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return conn.WriteJSON(Message{Type: typ, Payload: data})
}
