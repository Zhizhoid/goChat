package main

import (
	"log"
	"net"
)

const adress string = ":8080"

func main() {
	db, err := NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfuly opened the db")

	listener, err := net.Listen("tcp", adress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Started listening on %v\n", adress)

	for {
		conn, err := listener.Accept()
		log.Printf("Accepted connction from %v\n", conn.RemoteAddr())
		if err != nil {
			log.Println("An error occured: ", err)
		}

		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *Database) {
	buf := make([]byte, 2000)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("An error occured: ", err)
			break
		}

		responseBytes := db.HandleRequest(buf[:n])
		n, err = conn.Write(responseBytes)

		if err != nil {
			log.Println("An error occured: ", err)
			break
		} else {
			log.Printf("Wrote %v bytes to %v", n, conn.LocalAddr().String())
		}

		// db.PrintUsers()
		// db.PrintUsernamesToIDs()
		// db.PrintRooms()
	}
}
