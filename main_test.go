package main

import (
	"encoding/json"
	"log"
	"net"
	"testing"
)

func TestHandleConnection(t *testing.T) {
	input := []byte("{\"method\": \"GIBBERISH\"}")
	expected := AuthResponse{HTTPCode: "200", Message: "Invalid protocol method."}

	server, client := net.Pipe()

	go func() {
		handleConnection(server)
	}()

	_, err := client.Write(input)
	if err != nil {
		log.Fatalf("Programming error: Unable to write to mock client. Is the server listening? Input: %+v\n", input)
	}

	var authResponse AuthResponse
	d := json.NewDecoder(client)
	err = d.Decode(&authResponse)
	if err != nil {
		log.Fatal("Unable to read response from server.")
	}

	if authResponse != expected {
		log.Fatalf("Incorrect result: authResponse. Got: %+v, Expected %+v", authResponse, expected)
	}

	client.Close()
}

func TestHandleConnectionDecodeFailure(t *testing.T) {
	input := []byte("{\"metd:GIBBERISH\n")

	server, client := net.Pipe()

	go func() {
		handleConnection(server)
	}()

	_, err := client.Write(input)
	if err != nil {
		log.Fatalf("Programming error: Unable to write to mock client. Is the server listening? Input: %+v\n", input)
	}
}

func TestHandleConnectionClientClosed(t *testing.T) {
	input := []byte("{\"method\": \"GIBBERISH\"}")

	server, client := net.Pipe()

	go func() {
		handleConnection(server)
	}()

	_, err := client.Write(input)
	if err != nil {
		log.Fatalf("Programming error: Unable to write to mock client. Is the server listening? Input: %+v\n", input)
	}
	client.Close()

}
