package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

const adress string = ":8080"

var db *Database

func main() {
	var err error
	db, err = NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfuly opened the db")

	handleHTTP()

	// go handleTCP()
}

func handleHTTP() {
	http.HandleFunc("/", HTTPHandler)

	err := http.ListenAndServe(adress, nil)
	log.Fatal(err)
}

func HTTPHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	switch req.Method {
	case "POST":
		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println("An error occured: ", err)
		}

		responseBytes := db.HandleRequest(bytes)
		_, err = w.Write(responseBytes)
	case "OPTIONS":

		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func handleTCP() {
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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
	}
}
