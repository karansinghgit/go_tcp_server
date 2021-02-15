package main

import (
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Printf("TCP SERVER BEGUN")
	db := getDBConnection()
	defer db.Close()

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Unable to listen on port %s. Error: %s", ":9000", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Unable to accept a connection on TCP Port")
		}
		go handleConnection(conn)
	}
}
