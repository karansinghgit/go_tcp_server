package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

var db *sql.DB

func getDBConnection() *sql.DB {
	//The format for DSN is <username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/entry_task_db")
	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println("The connection to the database could not be established!")
	}

	return db
}

func getResponse(authRequest AuthRequest) AuthResponse {
	var authResponse AuthResponse
	if authRequest.Method == "LOGIN" {
		authResponse = loginUser(authRequest.UserInfo)
	} else {
		authResponse = registerUser(authRequest.UserInfo)
	}

	return authResponse
}

func handleConnection(conn net.Conn) {
	var authRequest AuthRequest
	var authResponse AuthResponse
	var connReader = bufio.NewReader(conn)
	defer conn.Close()

	for {
		rawRequest, err := connReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Hung connection. Error: %s", err)
				return
			}

			log.Printf("Unable to read request. Error: %s\n", err)
			_, _ = conn.Write([]byte("Server error: Unable to read request message\n"))
			continue
		}

		err = json.Unmarshal([]byte(rawRequest), &authRequest)
		if err != nil {
			log.Printf("Failed to unmarshal payload. Error: %s\n", err)
			_, _ = conn.Write([]byte("Server error: Malformed request\n"))
			continue
		}

		authResponse = getResponse(authRequest)

		response, err := json.Marshal(authResponse)
		if err != nil {
			log.Printf("Failed to marshal JSON object. Error: %s, Object: %+v\n", err, authResponse)
			response, _ = json.Marshal(
				AuthResponse{})
		}

		_, err = conn.Write(response)
		if err != nil {
			log.Printf("Unable to respond to client. Error: %s, Request: %+v\n", err, authRequest)
		}
	}
}
