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
	_, err := db.sqlDB.Exec(q, username, password, name)
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
	row := db.sqlDB.QueryRow(q, username, password)

	var id uint64
	err := row.Scan(&id)
	if err != nil {
		// return 0, errors.New("Invalid username and/or password")
		return 0, err
	}

	return db.sm.NewSession(id), nil
}

// Room
func (db *Database) AddRoom(name string, ownerId uint64) error {
	q := "INSERT INTO rooms(`Name`, OwnerID) VALUES(?, ?);"
	_, err := db.sqlDB.Exec(q, name, ownerId)
	if err != nil {
		return err
	}

	return nil
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

func (db *Database) LoginRoom(roomName string, userId uint64) error {
	//Getting room ID
	q := "SELECT id FROM rooms WHERE `Name` = ?;"
	row := db.sqlDB.QueryRow(q, roomName)

	var roomId uint64
	err := row.Scan(roomId)
	if err != nil {
		return err
	}

	//Creating new row in users_rooms
	q = "INSERT INTO users_rooms(UserID, RoomID) VALUES(?, ?);"
	_, err = db.sqlDB.Exec(q, userId, roomId)
	if err != nil {
		return err
	}

	return nil
}

// Message TODO
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
