package main

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	sqlDB *sql.DB
	sm    *SessionManager
}

func NewDatabase() (*Database, error) {
	var db Database

	cfg := mysql.NewConfig()
	cfg.Addr = "localhost"
	cfg.User = "root"
	cfg.Passwd = "masterkey"
	cfg.DBName = "gochatdb"
	cfg.ParseTime = true

	var err error

	db.sqlDB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	db.sm = NewSessionManager()

	return &db, nil
}

func (db *Database) HandleRequest(requestBytes []byte) (responseBytes []byte) {
	var action Action
	err := json.Unmarshal(requestBytes, &action)
	if err != nil {
		return unknownFailedResponseBytes(err.Error())
	}

	var object JsonObject
	switch action.ObjectName {
	case "user":
		object = &User{}
	case "room":
		object = &Room{}
	case "message":
		object = &Message{}
	default:
		return unknownFailedResponseBytes("Undefined object name")
	}

	var defAction DefinedAction
	switch action.Action {
	case "create":
		defAction, err = object.GetCreateAction()
	case "update":
		defAction, err = object.GetUpdateAction()
	case "delete":
		defAction, err = object.GetDeleteAction()
	case "login":
		defAction, err = object.GetLoginAction()
	default:
		return unknownFailedResponseBytes("Undefined action")
	}

	if err != nil {
		return unknownFailedResponseBytes(err.Error())
	}

	err = json.Unmarshal(requestBytes, defAction)
	if err != nil {
		return unknownFailedResponseBytes(err.Error())
	}

	response := defAction.Process(db)
	responseBytes, err = json.Marshal(response)
	if err != nil {
		return unknownFailedResponseBytes(err.Error())
	}

	return
}

func unknownFailedResponseBytes(status string) (responseBytes []byte) {
	responseBytes, err := json.Marshal(Response{
		Action:     "unknown",
		ObjectName: "unknown",
		Success:    false,
		Status:     status,
	})
	if err != nil {
		return nil
	}
	return
}

// User
func (db *Database) AddUser(username string, password string, name string) error {
	q := "INSERT INTO users(Username, `Password`, `Name`) VALUES(?, ?, ?);"

	_, err := db.sqlDB.Query(q, username, password, name)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUser(id uint64, newUsername string, newPassword string, newName string) error {
	q := "UPDATE users SET"

	everyFieldIsEmpty := true
	queryArgs := make([]interface{}, 0)

	if newUsername != "" {
		q += " Username = ?"
		queryArgs = append(queryArgs, newUsername)
		everyFieldIsEmpty = false
	}

	if newPassword != "" {
		if !everyFieldIsEmpty {
			q += ","
		}
		q += " `Password` = ?"
		queryArgs = append(queryArgs, newPassword)
		everyFieldIsEmpty = false
	}

	if newName != "" {
		if !everyFieldIsEmpty {
			q += ","
		}
		q += " `Name` = ?"
		queryArgs = append(queryArgs, newName)
		everyFieldIsEmpty = false
	}

	if everyFieldIsEmpty {
		return errors.New("At least one field should be changed")
	}

	q += " WHERE id = ?;"
	queryArgs = append(queryArgs, id)

	_, err := db.sqlDB.Query(q, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteUser(id uint64) error {
	q := "DELETE FROM users WHERE id = ?;"

	_, err := db.sqlDB.Query(q, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) LoginUser(username string, password string) (uint64, error) {
	q := "SELECT id FROM users WHERE Username = ? AND `Password` = ?;"

	rows, err := db.sqlDB.Query(q, username, password)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, errors.New("Invalid username and/or password")
	}

	var id uint64
	rows.Scan(&id)

	return db.sm.NewSession(id), nil
}

// Room
func (db *Database) AddRoom(name string) {
	// db.rooms = append(db.rooms, NewRoom(uint64(len(db.rooms)), name))
}

func (db *Database) UpdateRoom(id uint64, newName string) error {
	// if id >= uint64(len(db.rooms)) || id == 0 {
	// 	return errors.New("Invalid room ID")
	// }

	// db.rooms[id].Name = newName

	return nil
}

func (db *Database) DeleteRoom(id uint64) error {
	// if id >= uint64(len(db.rooms)) || id == 0 {
	// 	return errors.New("Invalid room ID")
	// }

	// db.rooms = append(db.rooms[:id], db.rooms[id+1:]...)

	return nil
}

// Message
func (db *Database) SendMessage(content Content, senderId, roomId, replyToId uint64) error {
	// if roomId >= uint64(len(db.rooms)) || roomId == 0 {
	// 	return errors.New("Invalid room ID")
	// }
	// if senderId >= uint64(len(db.rooms)) || senderId == 0 {
	// 	return errors.New("Invalid sender ID")
	// }
	// if senderId >= uint64(len(db.rooms)) {
	// 	return errors.New("Invalid replyTo ID")
	// }

	// db.rooms[roomId].SendMessage(content, senderId, replyToId)

	return nil
}

func (db *Database) UpdateMessage(roomId, messageId uint64, newContent Content, newReplyToId uint64) error {
	// if roomId >= uint64(len(db.rooms)) {
	// 	return errors.New("Invalid room ID")
	// }

	// return db.rooms[roomId].UpdateMessage(messageId, newContent, newReplyToId)

	return nil
}

func (db *Database) DeleteMessage(roomId, messageId uint64) error {
	// if roomId >= uint64(len(db.rooms)) || roomId == 0 {
	// 	return errors.New("Invalid room ID")
	// }

	// return db.rooms[roomId].DeleteMessage(messageId)
	return nil
}

// package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// )

// type Database struct {
// 	usernamesToIDs map[string]uint64
// 	users          []User
// 	rooms          []*Room
// }

// func NewDatabase() *Database {
// 	var db Database
// 	db.users = make([]User, 1)
// 	db.rooms = make([]*Room, 1)
// 	db.usernamesToIDs = make(map[string]uint64)

// 	return &db
// }

// func (db *Database) HandleRequest(requestBytes []byte) (responseBytes []byte) {
// 	var action Action
// 	err := json.Unmarshal(requestBytes, &action)
// 	if err != nil {
// 		return unknownFailedResponseBytes(err.Error())
// 	}

// 	var object JsonObject
// 	switch action.ObjectName {
// 	case "user":
// 		object = &User{}
// 	case "room":
// 		object = &Room{}
// 	case "message":
// 		object = &Message{}
// 	default:
// 		return unknownFailedResponseBytes("Undefined object name")
// 	}

// 	var defAction DefinedAction
// 	switch action.Action {
// 	case "create":
// 		defAction, err = object.GetCreateAction()
// 	case "update":
// 		defAction, err = object.GetUpdateAction()
// 	case "delete":
// 		defAction, err = object.GetDeleteAction()
// 	case "login":
// 		defAction, err = object.GetLoginAction()
// 	default:
// 		return unknownFailedResponseBytes("Undefined action")
// 	}

// 	if err != nil {
// 		return unknownFailedResponseBytes(err.Error())
// 	}

// 	err = json.Unmarshal(requestBytes, defAction)
// 	if err != nil {
// 		return unknownFailedResponseBytes(err.Error())
// 	}

// 	response := defAction.Process(db)
// 	responseBytes, err = json.Marshal(response)
// 	if err != nil {
// 		return unknownFailedResponseBytes(err.Error())
// 	}

// 	return
// }

// func unknownFailedResponseBytes(status string) (responseBytes []byte) {
// 	responseBytes, err := json.Marshal(Response{
// 		Action:     "unknown",
// 		ObjectName: "unknown",
// 		Success:    false,
// 		Status:     status,
// 	})
// 	if err != nil {
// 		return nil
// 	}
// 	return
// }

// func (db *Database) PrintUsers() {
// 	fmt.Println(db.users)
// }

// func (db *Database) PrintRooms() {
// 	fmt.Print("Rooms: ")
// 	for _, room := range db.rooms {
// 		if room == nil {
// 			continue
// 		}
// 		fmt.Print("{")
// 		room.Print()
// 		fmt.Print("}, ")
// 	}
// 	fmt.Println()
// }

// func (db *Database) PrintUsernamesToIDs() {
// 	fmt.Println(db.usernamesToIDs)
// }

// // User
// func (db *Database) AddUser(username string, password string, name string) error {
// 	_, exists := db.usernamesToIDs[username]
// 	if exists {
// 		return errors.New("Username already used")
// 	}

// 	db.usernamesToIDs[username] = uint64(len(db.users))

// 	db.users = append(db.users, User{
// 		// ID:       uint64(len(db.users)),
// 		Username: username,
// 		Name:     name,
// 		Password: password,
// 	})

// 	return nil
// }

// func (db *Database) UpdateUser(id uint64, newName string, newPassword string) error {
// 	if id >= uint64(len(db.users)) || id == 0 {
// 		return errors.New("Invalid user ID")
// 	}

// 	db.users[id].Name = newName
// 	db.users[id].Password = newPassword

// 	return nil
// }

// func (db *Database) DeleteUser(id uint64) error {
// 	if id >= uint64(len(db.users)) || id == 0 {
// 		return errors.New("Invalid user ID")
// 	}

// 	username := db.users[id].Username
// 	db.users = append(db.users[:id], db.users[id+1:]...)
// 	delete(db.usernamesToIDs, username)

// 	return nil
// }

// func (db *Database) LoginUser(username string, password string) (uint64, error) {
// 	id, exists := db.usernamesToIDs[username]

// 	if !exists {
// 		return 0, errors.New("Invalid username")
// 	}

// 	if db.users[id].Password != password {
// 		return 0, errors.New("Invalid password")
// 	}

// 	return id, nil
// }

// // Room
// func (db *Database) AddRoom(name string) {
// 	db.rooms = append(db.rooms, NewRoom(uint64(len(db.rooms)), name))
// }

// func (db *Database) UpdateRoom(id uint64, newName string) error {
// 	if id >= uint64(len(db.rooms)) || id == 0 {
// 		return errors.New("Invalid room ID")
// 	}

// 	db.rooms[id].Name = newName

// 	return nil
// }

// func (db *Database) DeleteRoom(id uint64) error {
// 	if id >= uint64(len(db.rooms)) || id == 0 {
// 		return errors.New("Invalid room ID")
// 	}

// 	db.rooms = append(db.rooms[:id], db.rooms[id+1:]...)

// 	return nil
// }

// // Message
// func (db *Database) SendMessage(content Content, senderId, roomId, replyToId uint64) error {
// 	if roomId >= uint64(len(db.rooms)) || roomId == 0 {
// 		return errors.New("Invalid room ID")
// 	}
// 	if senderId >= uint64(len(db.rooms)) || senderId == 0 {
// 		return errors.New("Invalid sender ID")
// 	}
// 	if senderId >= uint64(len(db.rooms)) {
// 		return errors.New("Invalid replyTo ID")
// 	}

// 	db.rooms[roomId].SendMessage(content, senderId, replyToId)

// 	return nil
// }

// func (db *Database) UpdateMessage(roomId, messageId uint64, newContent Content, newReplyToId uint64) error {
// 	if roomId >= uint64(len(db.rooms)) {
// 		return errors.New("Invalid room ID")
// 	}

// 	return db.rooms[roomId].UpdateMessage(messageId, newContent, newReplyToId)
// }

// func (db *Database) DeleteMessage(roomId, messageId uint64) error {
// 	if roomId >= uint64(len(db.rooms)) || roomId == 0 {
// 		return errors.New("Invalid room ID")
// 	}

// 	return db.rooms[roomId].DeleteMessage(messageId)
// }
