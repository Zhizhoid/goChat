package main

import (
	"chat/database"
	"io"
	"log"
	"net/http"
)

const adress string = ":8080"

var db *database.Database

func main() {
	var err error
	db, err = database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfuly opened the db")

	handleHTTP()
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
