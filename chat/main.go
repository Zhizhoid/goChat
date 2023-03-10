package main

import (
	"log"
	"net"
)

func main() {
	db := NewDatabase()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

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

		db.PrintUsers()
		db.PrintUsernamesToIDs()
		db.PrintRooms()
	}
}
