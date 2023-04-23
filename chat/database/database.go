package database

import (
	"chat/safety"

	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/argon2"
)

const SALT_LENGTH int = 30
const PASSWORDHASH_LENGTH uint32 = 60
const ARGON_TIME uint32 = 1
const ARGON_MEMORY uint32 = 47104
const ARGON_THREADS uint8 = 1

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
	case "read":
		defAction, err = object.GetReadAction()
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
	q := "INSERT INTO users(Username, PasswordHash, Salt, `Name`) VALUES(?, ?, ?, ?);"

	salt, err := safety.RandomString(SALT_LENGTH)
	passwordHash := argon2.IDKey([]byte(password), salt, ARGON_TIME, ARGON_MEMORY, ARGON_THREADS, PASSWORDHASH_LENGTH)

	_, err = db.sqlDB.Exec(q, username, base64.StdEncoding.EncodeToString(passwordHash), base64.StdEncoding.EncodeToString(salt), name)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateUser(sessionId uint64, newUsername string, newPassword string, newName string) error { // TOTEST
	id, err := db.sm.GetUserIDFromSession(sessionId)
	if err != nil {
		return err
	}

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
		q += " PasswordHash = ?, Salt = ?"

		salt, err := safety.RandomString(SALT_LENGTH)
		if err != nil {
			return err
		}

		newPasswordHash := argon2.IDKey([]byte(newPassword), salt, ARGON_TIME, ARGON_MEMORY, ARGON_THREADS, PASSWORDHASH_LENGTH)

		queryArgs = append(queryArgs, base64.StdEncoding.EncodeToString(newPasswordHash), base64.StdEncoding.EncodeToString(salt))
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

	_, err = db.sqlDB.Query(q, queryArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteUser(sessionId uint64) error {
	id, err := db.sm.GetUserIDFromSession(sessionId)
	if err != nil {
		return err
	}

	q := "DELETE FROM users WHERE id = ?;"
	_, err = db.sqlDB.Query(q, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) LoginUser(username string, password string) (uint64, error) {
	q := "SELECT id, Salt, PasswordHash FROM users WHERE Username = ?;"
	row := db.sqlDB.QueryRow(q, username)

	var (
		id                       uint64
		salt64, dbPasswordHash64 string
	)

	err := row.Scan(&id, &salt64, &dbPasswordHash64)
	if err != nil {
		// return 0, errors.New("Username may be invalid")
		return 0, err
	}

	salt, err := base64.StdEncoding.DecodeString(salt64)
	if err != nil {
		return 0, err
	}

	inputPasswordHash := argon2.IDKey([]byte(password), salt, ARGON_TIME, ARGON_MEMORY, ARGON_THREADS, PASSWORDHASH_LENGTH)

	dbPasswordHash, err := base64.StdEncoding.DecodeString(dbPasswordHash64)
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(PASSWORDHASH_LENGTH); i++ {
		if inputPasswordHash[i] != dbPasswordHash[i] {
			return 0, errors.New("Invalid password")
		}
	}

	return db.sm.NewSession(id)
}

func (db *Database) ReadUser(username string) (name string, err error) {
	return "", errors.New("Unimplemented")
}

func (db *Database) ReadUserRooms(username string) (rooms []string, err error) {
	return nil, errors.New("Unimplemented")
}

// Room
func (db *Database) AddRoom(name string, ownerSessionId uint64) error {
	ownerId, err := db.sm.GetUserIDFromSession(ownerSessionId)
	if err != nil {
		return err
	}

	q := "INSERT INTO rooms(`Name`, OwnerID) VALUES(?, ?);"
	_, err = db.sqlDB.Exec(q, name, ownerId)
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

func (db *Database) LoginRoom(roomName string, sessionID uint64) error {
	userId, err := db.sm.GetUserIDFromSession(sessionID)
	if err != nil {
		return err
	}

	//Getting room ID
	q := "SELECT id FROM rooms WHERE `Name` = ?;"
	row := db.sqlDB.QueryRow(q, roomName)

	var roomId uint64
	err = row.Scan(roomId)
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
