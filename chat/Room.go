package main

import (
	"fmt"
)

type Room struct{}

// func NewRoom(id uint64, name string) *Room {
// 	var room Room

// 	room.ID = id
// 	room.Name = name

// 	room.messages = make([]*Message, 1)

// 	return &room
// }

func (room *Room) SendMessage(content Content, senderId, replyToId uint64) {
	// room.messages = append(room.messages, &Message{
	// 	Content:   content,
	// 	SentAt:    time.Now(),
	// 	AuthorID:  senderId,
	// 	ReplyToID: replyToId,
	// 	Edited:    false,
	// })
}

func (room *Room) UpdateMessage(messageId uint64, newContent Content, newReplyToId uint64) error {
	// if messageId >= uint64(len(room.messages)) {
	// 	return errors.New("Invalid message ID")
	// }

	// if newReplyToId >= messageId {
	// 	newReplyToId = 0
	// }

	// room.messages[messageId].Content = newContent
	// room.messages[messageId].ReplyToID = newReplyToId
	// room.messages[messageId].Edited = true

	return nil
}

func (room *Room) DeleteMessage(messageId uint64) error {
	// if messageId >= uint64(len(room.messages)) {
	// 	return errors.New("Invalid message ID")
	// }

	// room.messages[messageId] = nil
	return nil
}

// Create action TOCHECK
type RoomCreate struct {
	Data struct {
		SessionID uint64 `json:"sessionId"`
		Name      string `json:"name"`
	} `json:"data"`
}

func (room *Room) GetCreateAction() (DefinedAction, error) {
	return &RoomCreate{}, nil
}

func (action *RoomCreate) Process(db *Database) Response {
	userId, err := db.sm.GetUserIDFromSession(action.Data.SessionID)
	if err != nil {
		return roomResponse("create", false, err.Error())
	}

	if action.Data.Name == "" {
		return roomResponse("create", false, "Creating room failed, room name cannon be empty")
	}

	err = db.AddRoom(action.Data.Name, userId)
	if err != nil {
		return roomResponse("create", false, err.Error())
	}

	return roomResponse("create", true, fmt.Sprintf("Room %v successfully created", action.Data.Name))
}

// Update action TODO
type RoomUpdate struct {
	Data struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

func (room *Room) GetUpdateAction() (DefinedAction, error) {
	return &RoomUpdate{}, nil
}

func (action *RoomUpdate) Process(db *Database) Response {
	fmt.Printf("Updating room: %v %v\n", action.Data.ID, action.Data.Name)

	if action.Data.Name == "" {
		return roomResponse("update", false, "Updating room failed, room name cannon be empty")
	}

	err := db.UpdateRoom(action.Data.ID, action.Data.Name)
	if err != nil {
		return roomResponse("update", false, err.Error())
	}

	return userResponse("update", true, fmt.Sprintf("Room %v successfully created", action.Data.Name))
}

// Delete action TODO
type RoomDelete struct {
	Data struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

func (room *Room) GetDeleteAction() (DefinedAction, error) {
	return &RoomDelete{}, nil
}

func (action *RoomDelete) Process(db *Database) Response {
	fmt.Printf("Deleting room %v\n", action.Data.ID)

	err := db.DeleteRoom(action.Data.ID)
	if err != nil {
		return roomResponse("delete", false, err.Error())
	}
	return roomResponse("delete", true, fmt.Sprintf("Room successfully deleted"))
}

// Login action
type RoomLogin struct {
	Data struct {
		SessionID uint64 `json:"sessionId"`
		RoomName  string `json:"roomName"`
	} `json:"data"`
}

func (m *Room) GetLoginAction() (DefinedAction, error) {
	return &RoomLogin{}, nil
}

func (action *RoomLogin) Process(db *Database) Response {
	userId, err := db.sm.GetUserIDFromSession(action.Data.SessionID)
	if err != nil {
		return roomResponse("login", false, err.Error())
	}

	err = db.LoginRoom(action.Data.RoomName, userId)
	if err != nil {
		return roomResponse("login", false, err.Error())
	}

	return roomResponse("login", true, "Successfully joined the room")
}

// OTHER
func roomResponse(action string, success bool, status string) Response {
	return Response{
		Action:     action,
		ObjectName: "room",
		Success:    success,
		Status:     status,
	}
}
