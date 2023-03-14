package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

type Response struct {
	Action     string `json:"action"`
	ObjectName string `json:"object"`
	Success    bool   `json:"success"`
	Status     string `json:"status"`
	ID         uint64 `json:"id"`
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 2000)
	for {
		bytes, err := handleCommand()

		if err != nil {
			log.Println(err)
			continue
		}

		_, err = conn.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}

		n, err := conn.Read(buf)

		var response Response

		err = json.Unmarshal(buf[:n], &response)
		if err != nil {
			fmt.Println("Something is wrong with the response")
			continue
		}
		fmt.Println(response)
	}
}

func handleCommand() ([]byte, error) {
	var action, object string
	fmt.Scanf("%v %v\n", &action, &object)

	switch action {
	case "create":
		switch object {
		case "user":
			return handleUserCreate(), nil
		case "room":
		case "message":
		}
	case "update":
		switch object {
		case "user":
			return handleUserUpdate(), nil
		case "room":
		case "message":
		}
	case "delete":
		switch object {
		case "user":
			return handleUserDelete(), nil
		}
	case "login":
		switch object {
		case "user":
			return handleUserLogin(), nil
		}
	}

	return nil, errors.New("Unknown command")
}

// User

func handleUserCreate() []byte {
	var username, password, name string
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)
	fmt.Print("Name: ")
	fmt.Scan(&name)

	return GetUserCreateBytes(username, password, name)
}

type UserCreate struct {
	Action     string `json:"action"`
	ObjectName string `json:"object"`
	Data       struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"data"`
}

func GetUserCreateBytes(username string, password string, name string) []byte {
	userCreate := UserCreate{
		Action:     "create",
		ObjectName: "user",
	}
	userCreate.Data.Username = username
	userCreate.Data.Password = password
	userCreate.Data.Name = name

	bytes, _ := json.Marshal(userCreate)

	return bytes
}

func handleUserUpdate() []byte {
	var id uint64
	var username, password, name string
	fmt.Print("ID: ")
	fmt.Scanf("%v\n", &id)
	fmt.Print("Username: ")
	fmt.Scanf("%v\n", &username)
	fmt.Print("Password: ")
	fmt.Scanf("%v\n", &password)
	fmt.Print("Name: ")
	fmt.Scanf("%v\n", &name)

	log.Println("id: ", id, "uName: ", username, "pass: ", password, "name: ", name)

	return GetUserUpdateBytes(id, username, password, name)
}

type UserUpdate struct {
	Action     string `json:"action"`
	ObjectName string `json:"object"`
	Data       struct {
		ID       uint64 `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"data"`
}

func GetUserUpdateBytes(id uint64, username string, password string, name string) []byte {
	userUpdate := UserUpdate{
		Action:     "update",
		ObjectName: "user",
	}
	userUpdate.Data.ID = id
	userUpdate.Data.Username = username
	userUpdate.Data.Password = password
	userUpdate.Data.Name = name

	log.Println(userUpdate)

	bytes, _ := json.Marshal(userUpdate)

	return bytes
}

func handleUserDelete() []byte {
	var id uint64
	fmt.Print("ID: ")
	fmt.Scan(&id)

	return GetUserDeleteBytes(id)
}

type UserDelete struct {
	Action     string `json:"action"`
	ObjectName string `json:"object"`
	Data       struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

func GetUserDeleteBytes(id uint64) []byte {
	userDelete := UserDelete{
		Action:     "delete",
		ObjectName: "user",
	}

	userDelete.Data.ID = id

	bytes, _ := json.Marshal(userDelete)

	return bytes
}

func handleUserLogin() []byte {
	var username, password string
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)

	return GetUserLoginBytes(username, password)
}

type UserLogin struct {
	Action     string `json:"action"`
	ObjectName string `json:"object"`
	Data       struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"data"`
}

func GetUserLoginBytes(username string, password string) []byte {
	userLogin := UserLogin{
		Action:     "login",
		ObjectName: "user",
	}
	userLogin.Data.Username = username
	userLogin.Data.Password = password

	bytes, _ := json.Marshal(userLogin)

	return bytes
}
